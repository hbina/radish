use nom::{bytes::complete::tag_no_case, character::complete::digit1};
use std::{
    collections::{hash_map::Entry, HashMap},
    sync::{Arc, Mutex},
    time::SystemTime,
};

use crate::{config::get_default_config, resp::Resp};

#[derive(Clone)]
pub struct Server {
    pub config: Arc<Mutex<HashMap<Resp, Resp>>>,
    pub dbs: Arc<Mutex<HashMap<u64, LockDb>>>,
}

pub type LockDb = Arc<Mutex<Db>>;

pub struct Db {
    pub kv: HashMap<Resp, Resp>,
    pub expiry: HashMap<Resp, u128>,
}

impl Db {
    pub fn set_value(&mut self, key: Resp, value: Resp) -> Option<Resp> {
        self.kv.insert(key, value)
    }

    pub fn remove_value(&mut self, key: &Resp) -> Option<Resp> {
        self.kv.remove(key)
    }

    pub fn get_value(&mut self, key: &Resp) -> Option<&Resp> {
        let current_ms = get_current_ms();
        if let Some(expiry_ms) = self.expiry.get(&key) {
            if current_ms > *expiry_ms {
                self.kv.remove(&key);
                self.expiry.remove(&key);
            }
        }
        self.kv.get(&key)
    }

    pub fn entry_value<'a>(&'a mut self, key: Resp) -> Entry<'a, Resp, Resp> {
        self.kv.entry(key)
    }

    // Set expiry time measured in milliseconds since UNIX epoch
    pub fn set_expiry(&mut self, key: Resp, value: u128) -> Option<u128> {
        self.expiry.insert(key, value)
    }

    // Get expiry time measured in milliseconds since UNIX epoch
    pub fn get_expiry(&self, key: &Resp) -> Option<&u128> {
        self.expiry.get(&key)
    }

    pub fn remove_expiry(&mut self, key: &Resp) -> Option<u128> {
        self.expiry.remove(key)
    }

    pub fn entry_expiry<'a>(&'a mut self, key: Resp) -> Entry<'a, Resp, u128> {
        self.expiry.entry(key)
    }

    pub fn contains_value(&mut self, key: &Resp) -> bool {
        let current_ms = get_current_ms();
        if let Some(expiry_ms) = self.expiry.get(&key) {
            if current_ms > *expiry_ms {
                self.kv.remove(&key);
                self.expiry.remove(&key);
            }
        }
        self.kv.contains_key(key)
    }

    pub fn contains_expiry(&mut self, key: &Resp) -> bool {
        self.expiry.contains_key(key)
    }

    pub fn clear(&mut self) {
        self.kv.clear();
        self.expiry.clear();
    }
}

impl Default for Db {
    fn default() -> Self {
        Self {
            kv: Default::default(),
            expiry: Default::default(),
        }
    }
}

impl Server {
    /// Returns the database pointed by `id`. Creates the db if it doesn't exist.
    pub fn get_db(&mut self, id: u64) -> LockDb {
        let mut dbs_lock = self.dbs.lock().unwrap();
        let result = dbs_lock.entry(id).or_default();
        result.clone()
    }
}

impl Default for Server {
    fn default() -> Self {
        let config = get_default_config();
        Self {
            config: Arc::new(Mutex::new(config)),
            dbs: Arc::new(Mutex::new(HashMap::default())),
        }
    }
}

#[derive(Debug, Clone)]
pub enum CommandError {
    NotEnoughArgs(&'static str),
    NotImplemented(String),
    SyntaxError,
    NotInteger,
    InvalidArgument,
}

pub type CommandResult = ::std::result::Result<Resp, CommandError>;

#[derive(Debug, Copy, Clone)]
pub enum SetMode {
    Nx,
    Xx,
    Gt,
    Lt,
}

pub fn get_current_ms() -> u128 {
    SystemTime::now()
        .duration_since(SystemTime::UNIX_EPOCH)
        .unwrap()
        .as_millis()
}

// [EX seconds | PX milliseconds | EXAT unix-time-seconds | PXAT unix-time-milliseconds | PERSIST]
#[derive(Debug, Copy, Clone)]
pub enum SetExpireMode {
    /// seconds
    Ex(u128),
    /// milliseconds
    Px(u128),
    /// seconds
    Exat(u128),
    /// milliseconds
    Pxat(u128),
    Persist,
}

impl SetExpireMode {
    pub fn parse_set_expire_mode(left: &[u8], right: &[u8]) -> Result<SetExpireMode, CommandError> {
        let tag_matcher = tag_no_case::<&[u8], &[u8], ()>;
        let value_matcher = |x| -> Option<u128> {
            // TODO(need to check that leftover is empty?)
            let (_, x_bytes) = digit1::<&[u8], ()>(x).ok()?;
            let x_str = std::str::from_utf8(x_bytes).ok()?;
            let x_u128 = x_str.parse::<u128>().ok()?;
            Some(x_u128)
        };
        if tag_matcher(b"EX")(left).is_ok() {
            let value = value_matcher(right).ok_or(CommandError::SyntaxError)?;
            return Ok(SetExpireMode::Ex(value));
        } else if tag_matcher(b"PX")(left).is_ok() {
            let value = value_matcher(right).ok_or(CommandError::SyntaxError)?;
            return Ok(SetExpireMode::Px(value));
        } else if tag_matcher(b"EXAT")(left).is_ok() {
            let value = value_matcher(right).ok_or(CommandError::SyntaxError)?;
            return Ok(SetExpireMode::Exat(value));
        } else if tag_matcher(b"PXAT")(left).is_ok() {
            let value = value_matcher(right).ok_or(CommandError::SyntaxError)?;
            return Ok(SetExpireMode::Pxat(value));
        } else if tag_matcher(b"PERSIST")(left).is_ok() {
            return Ok(SetExpireMode::Persist);
        } else {
            return Err(CommandError::SyntaxError);
        }
    }
}
