use std::{
    collections::HashMap,
    io::{Read, Write},
    net::{TcpListener, TcpStream},
    sync::{Arc, Mutex},
};

use anyhow::{anyhow, Result};
use redis_protocol_parser::{parse_resp, RespError, RespOwned, RespRef};

#[derive(Debug, Default)]
struct State {
    map: HashMap<Vec<u8>, RespOwned>,
}

impl State {
    pub fn select(&self, key: &[u8]) -> Option<&RespOwned> {
        self.map.get(key)
    }

    pub fn insert(&mut self, key: &[u8], resp: RespRef) -> Option<RespOwned> {
        self.map.insert(key.to_vec(), resp.as_owned())
    }
}

fn handle_client(state: Arc<Mutex<State>>, mut stream: TcpStream) -> Result<()> {
    stream.set_read_timeout(Some(std::time::Duration::from_secs(1)))?;

    let mut tmp = [0; 1024];
    let mut buffer = Vec::with_capacity(3);

    loop {
        let b = match stream.read(&mut tmp) {
            Ok(b) => b,
            Err(err) => {
                if buffer.len() == 0 {
                    return Ok(());
                } else {
                    return Err(anyhow!("{}", err));
                }
            }
        };

        if b == 0 {
            return Ok(());
        } else {
            buffer.extend_from_slice(&tmp[..b]);

            // Parse buffer
            match parse_resp(buffer.as_ref()) {
                Ok((resp, leftover)) => {
                    let state = state.lock().unwrap();
                    match &resp {
                        RespRef::Array(s) => {
                            if let Some(first) = s.first() {
                                match String::from_utf8(first.as_bytes())?.as_str() {
                                    "PING" => stream
                                        .write(&[b'+', b'P', b'I', b'N', b'G', b'\r', b'\n'])?,
                                    o => return Err(anyhow!("Unknown command '{}'", o)),
                                }
                            } else {
                                return Err(anyhow!("Resp must not be empty when used in request"));
                            }
                        }
                        o => return Err(anyhow!("Excepted Resp array but got '{:?}'", o)),
                    };
                    resp.write_to_writer(&mut stream)
                        .map_err(|s| anyhow!("{}", s))?;
                    stream.flush()?;
                    buffer.drain(0..buffer.len() - leftover.len());
                }
                Err(err) => match err {
                    RespError::NotEnoughBytes => continue,
                    err => {
                        return Err(anyhow!("err:{err}"));
                    }
                },
            }
        }
    }
}

fn main() -> std::io::Result<()> {
    let listener = TcpListener::bind("127.0.0.1:6380")?;

    let mut state = Arc::new(Mutex::new(State::default()));

    // accept connections and process them serially
    for stream in listener.incoming() {
        let state = state.clone();
        if let Err(err) = handle_client(state, stream?) {
            eprintln!("{err}");
        }
    }
    Ok(())
}
