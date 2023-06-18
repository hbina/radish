use bytes::{BufMut, Bytes, BytesMut};
use clap::Parser;

use itertools::Itertools;
use log::{debug, LevelFilter};
use misc::{CommandResult, LockDb, Server, SetMode};
use resp::{RespArray, RespArrayBuilder};
use std::{
    io::Write,
    sync::atomic::{AtomicU64, Ordering},
    time::Duration,
};
use tokio::{
    io::AsyncWriteExt,
    net::{TcpListener, TcpStream},
    task,
    time::timeout,
};

use crate::{
    misc::{get_current_ms, CommandError, SetExpireMode},
    resp::{cast_bytes_to_resp, Resp, RespParseError},
};

mod config;
mod misc;
mod resp;

#[tokio::main]
async fn main() -> std::io::Result<()> {
    let args = Args::parse();
    let port = args.port;
    let listener = TcpListener::bind(format!("127.0.0.1:{}", port)).await?;

    // Set up logger
    env_logger::Builder::new()
        .format(|buf, record| {
            writeln!(
                buf,
                "{}:{} {} [{}] - {}",
                record.file().unwrap_or("unknown"),
                record.line().unwrap_or(0),
                // chrono::Local::now().format("%Y-%m-%dT%H:%M:%S"),
                chrono::Local::now().format("%H:%M:%S"),
                record.level(),
                record.args()
            )
        })
        .filter(None, LevelFilter::Debug)
        .init();

    let server = Server::default();
    let client_id = AtomicU64::default();
    loop {
        let (stream, _) = listener.accept().await.unwrap();
        let id = client_id.fetch_add(1, Ordering::Relaxed);
        let server = server.clone();
        let mut client = Client {
            id,
            // The default connection will always be to db 0 according to redis
            db_id: 0,
            stream,
            server,
        };
        task::spawn(async move {
            client.handle_client_connection().await;
            debug!("client_id={} dies", id)
        });
    }
}

#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
struct Args {
    #[arg(short = 'p', long = "port")]
    port: u32,
}

pub struct Client {
    id: u64,
    db_id: u64,
    stream: TcpStream,
    // Handle to the server
    server: Server,
}

impl Client {
    /// Returns the database pointed by `id`. Creates the db if it doesn't exist.
    pub fn get_db(&mut self) -> LockDb {
        self.server.get_db(self.db_id)
    }

    async fn handle_client_connection(&mut self) {
        debug!("client_id={} starting", self.id);

        let mut total_buffer = BytesMut::with_capacity(1024);

        loop {
            // TODO(yuyuri) - Should we even have a timeout here? A redis-cli client that never sends anything is a perfectly reasonable client I think.
            // Currently just set to 1 hour.
            if timeout(Duration::from_secs(60 * 60), self.stream.readable())
                .await
                .unwrap()
                .is_err()
            {
                debug!("client_id={} timed out", self.id);
                return;
            }

            let mut read_buffer = [0; 1024];
            let count = match self.stream.try_read(&mut read_buffer) {
                Ok(count) => count,
                Err(ref e) if e.kind() == std::io::ErrorKind::WouldBlock => {
                    continue;
                }
                Err(err) => {
                    debug!("client_id={} tcp stream returned err={:#?}", self.id, err);
                    return;
                }
            };

            // TODO(yuyuri) - How to distinguish between reading zero because we haven't received anything or zero because the client ends?
            if count == 0 {
                debug!("client_id={} read nothing", self.id);
                return;
            }

            total_buffer.put_slice(&read_buffer[..count]);

            loop {
                match resp::get_resp_length(&total_buffer) {
                    Ok(resp_length) => {
                        let resp_bytes = total_buffer.split_to(resp_length).freeze();
                        debug!(
                            "<= {:#?}",
                            std::str::from_utf8(resp_bytes.as_ref()).unwrap_or("bad bytes")
                        );
                        let input_resp = cast_bytes_to_resp(resp_bytes).unwrap();
                        let r_output_resp = self.handle_client_resp(input_resp);
                        let output_resp = match r_output_resp {
                            Ok(output_resp) => output_resp,
                            Err(err) => match err {
                                CommandError::NotEnoughArgs(e) => Resp::res_from_str(&format!(
                                    "ERR wrong number of arguments for '{}' command",
                                    e
                                )),
                                CommandError::NotImplemented(e) => Resp::res_from_str(&format!(
                                    "ERR command '{}' is not yet implemented",
                                    e
                                )),
                                CommandError::SyntaxError => {
                                    Resp::res_from_static_str("ERR syntax error")
                                }
                                CommandError::NotInteger => Resp::res_from_static_str(
                                    "ERR value is not an integer or out of range",
                                ),
                                CommandError::InvalidArgument => {
                                    Resp::res_from_static_str("ERR invalid argument")
                                }
                            },
                        };
                        debug!(
                            "=> {:#?}",
                            std::str::from_utf8(output_resp.as_bytes().as_ref())
                                .unwrap_or("bad bytes")
                        );
                        let mut output_buffer = Vec::with_capacity(output_resp.as_bytes().len());
                        output_resp.write_to_writer(&mut output_buffer).unwrap();
                        let required_len = output_buffer.len();

                        let mut count = 0;
                        while count != required_len {
                            count += self.stream.write(&output_buffer).await.unwrap();
                        }
                    }
                    // For now just keep trying, but we should handle timeout or bad clients somehow
                    Err(RespParseError::NotEnoughBytes) => break,
                    Err(err) => {
                        debug!("client_id={} tcp stream returned err={:#?}", self.id, err);
                        return;
                    }
                }
            }
        }
    }

    fn handle_client_resp(&mut self, resp: Resp) -> CommandResult {
        match resp {
            Resp::String(r) => self.handle_resp_string_command(r),
            Resp::Error(_) => todo!("handle error"),
            Resp::Integer(_) => todo!("handle integer"),
            Resp::BulkString(_) => todo!("handle bulk string"),
            Resp::NilBulk => todo!("handle nil bulk"),
            Resp::Array(arr) => self.handle_resp_array_command(arr),
            Resp::NilArray => todo!("handle nil array"),
        }
    }

    // Figure out what command is being executed and let the specific command handler handles it
    fn handle_resp_string_command(&self, first_bytes: Bytes) -> CommandResult {
        if first_bytes.eq_ignore_ascii_case(b"ping") {
            return Ok(Resp::rbs_from_static_bytes(b"PONG"));
        }

        Err(CommandError::NotImplemented("string command".to_string()))
    }

    // Figure out what command is being executed and let the specific command handler handles it
    fn handle_resp_array_command(&mut self, arr: RespArray) -> CommandResult {
        let cmd_name = arr
            .get(0)
            .ok_or(CommandError::NotEnoughArgs("array command"))?;

        if cmd_name.as_bytes().eq_ignore_ascii_case(b"ping") {
            return self.handle_command_ping(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"config") {
            return self.handle_command_config(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"command") {
            return self.handle_cmd_command(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"set") {
            return self.handle_cmd_set(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"get") {
            return self.handle_cmd_get(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"incr") {
            return self.handle_cmd_incr(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"select") {
            return self.handle_cmd_select(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"flushall") {
            return self.handle_cmd_flushall(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"function") {
            return self.handle_cmd_function(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"del") {
            return self.handle_cmd_del(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"setnx") {
            return self.handle_cmd_setnx(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"expire") {
            return self.handle_cmd_expire(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"setex") {
            return self.handle_cmd_setex(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"getex") {
            return self.handle_cmd_getex(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"ttl") {
            return self.handle_cmd_ttl(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"pttl") {
            return self.handle_cmd_pttl(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"debug") {
            return self.handle_cmd_debug(arr);
        } else if cmd_name.as_bytes().eq_ignore_ascii_case(b"flushdb") {
            return self.handle_cmd_flushdb(arr);
        }

        return Err(CommandError::NotImplemented(
            std::str::from_utf8(cmd_name.as_bytes().as_ref())
                .unwrap_or("command name")
                .to_string(),
        ));
    }

    fn handle_command_ping(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() == 1 {
            Ok(Resp::rbs_from_static_bytes(b"PONG"))
        } else if arr.len() == 2 {
            return arr.get(1).ok_or(CommandError::NotEnoughArgs("ping"));
        } else {
            return Err(CommandError::NotEnoughArgs("ping"));
        }
    }

    fn handle_command_config(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() < 2 {
            return Err(CommandError::NotEnoughArgs("config"));
        }

        let second_arg = arr.get(1).ok_or(CommandError::NotEnoughArgs("config"))?;

        if second_arg.as_bytes().eq_ignore_ascii_case(b"get") {
            self.handle_cmd_config_get(arr)
        } else if second_arg.as_bytes().eq_ignore_ascii_case(b"set") {
            return self.handle_cmd_config_set(arr);
        } else {
            return Err(CommandError::NotEnoughArgs("config"));
        }
    }

    fn handle_cmd_config_get(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() < 3 {
            return Err(CommandError::NotEnoughArgs("config"));
        }

        let config = self.server.config.lock().unwrap();

        let mut result_arr = RespArrayBuilder::default();
        for key in arr.iter().skip(2) {
            if let Some(value) = config.get(&key) {
                result_arr.push_resp(key);
                result_arr.push_resp(value.clone());
            }
        }
        Ok(Resp::Array(result_arr.build()))
    }

    fn handle_cmd_config_set(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() < 3 {
            return Err(CommandError::NotEnoughArgs("config"));
        }

        let mut config = self.server.config.lock().unwrap();

        let pairs = arr.iter().skip(2).collect::<Vec<_>>();

        if pairs.len() % 2 != 0 {
            return Err(CommandError::NotEnoughArgs("config"));
        }

        // TODO(yuyuri) - Should only allow valid keys. Also need to check the exact behavior in redis.
        // For example, what does it do if some keys are valid and some are not? Does it check everything is good beforehand?
        // What kind of values are valid? Is it all just strings?
        for (key, value) in pairs.iter().tuples() {
            config.insert(key.clone(), value.clone());
        }
        Ok(Resp::rbs_from_static_bytes(b"OK"))
    }

    fn handle_cmd_command(&mut self, _: RespArray) -> CommandResult {
        Err(CommandError::NotImplemented("command".to_string()))
    }

    fn handle_cmd_set(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() < 3 {
            return Err(CommandError::SyntaxError);
        }

        let key = arr.get(1).ok_or(CommandError::NotEnoughArgs("set"))?;
        let value = arr.get(2).ok_or(CommandError::NotEnoughArgs("set"))?;

        let db_lock = self.get_db();
        let mut db = db_lock.lock().unwrap();

        // TODO(yuyuri) - Reuse this value somehow?
        let _ = db.kv.insert(key, value);

        Ok(Resp::rss_from_static_bytes(b"OK"))
    }

    fn handle_cmd_get(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() != 2 {
            return Err(CommandError::SyntaxError);
        }

        let key = arr.get(1).ok_or(CommandError::NotEnoughArgs("get"))?;

        let db_lock = self.get_db();
        let mut db = db_lock.lock().unwrap();

        // TODO(yuyuri) - Reuse this value somehow?
        let value = db.get_value(&key).cloned().unwrap_or(Resp::NilArray);

        Ok(value)
    }

    fn handle_cmd_del(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() != 2 {
            return Err(CommandError::SyntaxError);
        }

        let key = arr.get(1).ok_or(CommandError::NotEnoughArgs("get"))?;

        let db_lock = self.get_db();
        let mut db = db_lock.lock().unwrap();

        // TODO(yuyuri) - Reuse this value somehow?
        let _ = db.remove_value(&key);

        Ok(Resp::rbs_from_static_str("OK"))
    }

    fn handle_cmd_incr(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() != 2 {
            return Err(CommandError::SyntaxError);
        }

        let key = arr.get(1).ok_or(CommandError::NotEnoughArgs("incr"))?;

        let db_lock = self.get_db();
        let mut db = db_lock.lock().unwrap();

        let value = db.entry_value(key).or_insert(Resp::Integer(0));

        match value {
            Resp::Integer(v) => *v += 1,
            _ => return Err(CommandError::NotInteger),
        };

        Ok(value.clone())
    }

    fn handle_cmd_select(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() != 2 {
            return Err(CommandError::SyntaxError);
        }

        let key = arr
            .get(1)
            .ok_or(CommandError::NotEnoughArgs("select"))?
            .to_u64()
            .ok_or(CommandError::InvalidArgument)?;

        if key < 0 {
            return Err(CommandError::SyntaxError);
        }

        // TODO(yuyuri) - Check the conversion here
        self.db_id = key as u64;

        Ok(Resp::rbs_from_static_str("OK"))
    }

    fn handle_cmd_flushall(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() > 2 {
            return Err(CommandError::SyntaxError);
        }

        // TODO(yuyuri) - handle SYNC/ASYNC flushall.
        // For now we just ignore
        let _ = arr.get(1);

        let dbs_lock = self.server.dbs.lock().unwrap();

        // TODO(yuyuri) - this is very dangerous because we might have a deadlock!!
        for db in dbs_lock.values() {
            let mut db_lock = db.lock().unwrap();
            db_lock.clear();
        }

        Ok(Resp::rbs_from_static_str("OK"))
    }

    fn handle_cmd_function(&mut self, _: RespArray) -> CommandResult {
        // TODO(yuyuri) - Just a stub implementation
        Ok(Resp::rbs_from_static_str("OK"))
    }

    // SETNX key value
    fn handle_cmd_setnx(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() != 3 {
            return Err(CommandError::SyntaxError);
        }

        let key = arr.get(1).ok_or(CommandError::NotEnoughArgs("setnx"))?;
        let value = arr.get(2).ok_or(CommandError::NotEnoughArgs("setnx"))?;

        let db_lock = self.get_db();
        let mut db = db_lock.lock().unwrap();

        if !db.contains_value(&key) {
            db.set_value(key, value);
            return Ok(Resp::ri_from_u64(1));
        }
        return Ok(Resp::ri_from_u64(0));
    }

    // EXPIRE key seconds [NX | XX | GT | LT]
    fn handle_cmd_expire(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() != 3 && arr.len() != 4 {
            return Err(CommandError::SyntaxError);
        }

        let current_time = get_current_ms();
        let key = arr.get(1).ok_or(CommandError::NotEnoughArgs("expire"))?;
        let seconds = arr
            .get(2)
            .ok_or(CommandError::NotEnoughArgs("expire"))?
            .to_u128()
            .ok_or(CommandError::SyntaxError)?;
        let new_when = current_time + (seconds * 1000);

        let db_lock = self.get_db();
        let mut db = db_lock.lock().unwrap();

        let should_set = {
            match arr.get(3) {
                Some(m_expire_type) => match m_expire_type.to_set_mode() {
                    Some(expire_type) => match expire_type {
                        SetMode::Nx => !db.contains_expiry(&key),
                        SetMode::Xx => db.contains_expiry(&key),
                        SetMode::Gt => db.get_expiry(&key).is_none(),
                        SetMode::Lt => db
                            .get_expiry(&key)
                            .map(|when| new_when < *when)
                            .unwrap_or(false),
                    },
                    None => return Err(CommandError::SyntaxError),
                },
                None => true,
            }
        };

        if should_set {
            db.set_expiry(key, new_when);
            return Ok(Resp::ri_from_u64(1));
        } else {
            return Ok(Resp::ri_from_u64(0));
        }
    }

    // SETEX key seconds value
    pub fn handle_cmd_setex(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() != 4 {
            return Err(CommandError::SyntaxError);
        }

        let key = arr.get(1).ok_or(CommandError::NotEnoughArgs("setex"))?;
        let seconds = arr
            .get(2)
            .ok_or(CommandError::NotEnoughArgs("setex"))?
            .to_u128()
            .ok_or(CommandError::SyntaxError)?;
        let value = arr.get(3).ok_or(CommandError::NotEnoughArgs("setex"))?;
        let current_time = get_current_ms();
        let new_when = current_time + (1000 * seconds);

        let db_lock = self.get_db();
        let mut db = db_lock.lock().unwrap();

        db.set_expiry(key.clone(), new_when);
        db.set_value(key, value);

        Ok(Resp::rbs_from_static_str("OK"))
    }

    // GETEX key [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds | PERSIST]
    pub fn handle_cmd_getex(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() != 2 && arr.len() != 4 {
            return Err(CommandError::SyntaxError);
        }

        let key = arr.get(1).ok_or(CommandError::NotEnoughArgs("getex"))?;

        let m_set_expiry_mode = match (arr.get(2), arr.get(3)) {
            (Some(l), Some(r)) => Some(SetExpireMode::parse_set_expire_mode(
                &l.as_bytes(),
                &r.as_bytes(),
            )?),
            (None, None) => None,
            _ => return Err(CommandError::SyntaxError),
        };

        let db_lock = self.get_db();
        let mut db = db_lock.lock().unwrap();

        let value = db.get_value(&key).unwrap_or(&Resp::NilArray).clone();

        if value == Resp::NilArray {
            return Ok(value);
        }

        if let Some(set_expiry_mode) = m_set_expiry_mode {
            match set_expiry_mode {
                SetExpireMode::Ex(seconds) => {
                    let current_time = get_current_ms();
                    let new_when = current_time + (1000 * seconds);
                    db.set_expiry(key, new_when);
                }
                SetExpireMode::Px(mseconds) => {
                    let current_time = get_current_ms();
                    let new_when = current_time + mseconds;
                    db.set_expiry(key, new_when);
                }
                SetExpireMode::Exat(new_when_seconds) => {
                    db.set_expiry(key, new_when_seconds * 1000);
                }
                SetExpireMode::Pxat(new_when_mseconds) => {
                    db.set_expiry(key, new_when_mseconds);
                }
                SetExpireMode::Persist => todo!(),
            }
        };

        Ok(value)
    }

    // TTL key
    fn handle_cmd_ttl(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() != 2 {
            return Err(CommandError::SyntaxError);
        }

        let key = arr.get(1).ok_or(CommandError::NotEnoughArgs("ttl"))?;

        let db_lock = self.get_db();
        let mut db = db_lock.lock().unwrap();
        let current_time = get_current_ms();

        if !db.contains_value(&key) {
            return Ok(Resp::Integer(-2));
        } else {
            match db.get_expiry(&key) {
                Some(ttl) => {
                    let ts_diff = *ttl - current_time;
                    let ts_s = ts_diff / 1_000;
                    // TODO(yuyuri) - Handle conversion properly here
                    return Ok(Resp::Integer(ts_s as i64));
                }
                None => {
                    return Ok(Resp::Integer(-1));
                }
            }
        }
    }

    // PTTL key
    fn handle_cmd_pttl(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() != 2 {
            return Err(CommandError::SyntaxError);
        }

        let key = arr.get(1).ok_or(CommandError::NotEnoughArgs("pttl"))?;

        let db_lock = self.get_db();
        let mut db = db_lock.lock().unwrap();
        let current_time = get_current_ms();

        if !db.contains_value(&key) {
            return Ok(Resp::Integer(-2));
        } else {
            match db.get_expiry(&key) {
                Some(ttl) => {
                    let ts_diff = *ttl - current_time;
                    // TODO(yuyuri) - Handle conversion properly here
                    return Ok(Resp::Integer(ts_diff as i64));
                }
                None => {
                    return Ok(Resp::Integer(-1));
                }
            }
        }
    }

    // DEBUG
    fn handle_cmd_debug(&mut self, arr: RespArray) -> CommandResult {
        debug!("{:?}", arr);
        Ok(Resp::rbs_from_static_str("OK"))
    }

    // FLUSHDB [ASYNC | SYNC]
    fn handle_cmd_flushdb(&mut self, arr: RespArray) -> CommandResult {
        if arr.len() > 2 {
            return Err(CommandError::SyntaxError);
        }

        let mode = match arr.get(1) {
            Some(b) => match b.as_bytes().make_ascii_lowercase().as_bytes() {
                b"sync" => true,
                b"async" => false,
                _ => return Err(CommandError::SyntaxError),
            },
            None => false,
        };

        debug!("{:?}", arr);
        Ok(Resp::rbs_from_static_str("OK"))
    }
}
