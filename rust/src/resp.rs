use std::hash::{Hash, Hasher};
use std::io::Write;

use bytes::{BufMut, Bytes, BytesMut};

use crate::misc::SetMode;

type RespLengthResult<'a> = std::result::Result<usize, RespParseError>;
type RespResult = std::result::Result<Resp, RespParseError>;

const CR: u8 = b'\r';
const LF: u8 = b'\n';

#[derive(Debug, Clone, Eq, PartialEq)]
pub enum Resp {
    String(BytesMut),
    Error(BytesMut),
    Integer(i64),
    BulkString(BytesMut),
    NilBulk,
    Array(RespArray),
    NilArray,
}

impl Resp {
    pub fn res_from_str(input: &str) -> Self {
        return Resp::Error(BytesMut::from(input.as_bytes()));
    }

    pub fn res_from_static_str(input: &'static str) -> Self {
        return Resp::Error(BytesMut::from(input.as_bytes()));
    }

    pub fn res_from_bytes(input: &[u8]) -> Self {
        Resp::Error(BytesMut::from(input))
    }

    pub fn rbs_from_str(input: &str) -> Self {
        return Resp::BulkString(BytesMut::from(input.as_bytes()));
    }

    pub fn rbs_from_static_str(input: &'static str) -> Self {
        return Resp::BulkString(BytesMut::from(input.as_bytes()));
    }

    pub fn rbs_from_bytes(input: &[u8]) -> Self {
        Resp::BulkString(BytesMut::from(input))
    }

    pub fn rbs_from_static_bytes(input: &'static [u8]) -> Self {
        Resp::BulkString(BytesMut::from(input))
    }

    pub fn rss_from_static_bytes(input: &'static [u8]) -> Self {
        Resp::String(BytesMut::from(input))
    }

    pub fn ri_from_u64(input: i64) -> Self {
        Resp::Integer(input)
    }

    pub fn write_to_writer<W>(&self, writer: &mut W) -> Result<(), RespParseError>
    where
        W: Write,
    {
        match self {
            Resp::String(s) => {
                writer.write_all(b"+")?;
                writer.write_all(s)?;
                writer.write_all(b"\r\n")?;
            }
            Resp::Error(s) => {
                writer.write_all(b"-")?;
                writer.write_all(s)?;
                writer.write_all(b"\r\n")?;
            }
            Resp::Integer(s) => {
                writer.write_all(b":")?;
                writer.write_all(format!("{}", s).as_bytes())?;
                writer.write_all(b"\r\n")?;
            }
            Resp::BulkString(s) => {
                writer.write_all(b"$")?;
                writer.write_all(format!("{}", s.len()).as_bytes())?;
                writer.write_all(b"\r\n")?;
                writer.write_all(s)?;
                writer.write_all(b"\r\n")?;
            }
            Resp::NilBulk => writer.write_all(b"$-1\r\n")?,
            Resp::Array(a) => {
                writer.write_all(b"*")?;
                writer.write_all(format!("{}", a.0).as_bytes())?;
                writer.write_all(b"\r\n")?;
                for s in a.iter() {
                    s.write_to_writer(writer)?
                }
            }
            Resp::NilArray => writer.write_all(b"*-1\r\n")?,
        };
        Ok(())
    }

    pub fn as_bytes(&self) -> &[u8] {
        match self {
            Resp::String(b) => b.as_ref(),
            Resp::Error(b) => b.as_ref(),
            Resp::Integer(b) => &format!("{}", b).as_bytes(),
            Resp::BulkString(b) => b.as_ref(),
            Resp::NilBulk => b"-1", // TODO(yuyuri) - What should be returned here?
            Resp::Array(a) => a.1.as_ref(),
            Resp::NilArray => b"-1", // TODO(yuyuri) - What should be returned here?
        }
    }

    // The number of bytes required to fill this resp.
    pub fn resp_len(&self) -> usize {
        match self {
            Resp::String(b) => 1 + b.as_ref().len() + 2,
            Resp::Error(b) => 1 + b.as_ref().len() + 2,
            Resp::Integer(b) => 1 + format!("{}", *b).len() + 2,
            Resp::BulkString(b) => 1 + format!("{}", b.len()).len() + b.as_ref().len() + 2,
            Resp::NilBulk => 5,
            Resp::Array(a) => 1 + a.iter().map(|s| s.resp_len()).sum::<usize>(),
            Resp::NilArray => 5,
        }
    }

    pub fn to_u64(&self) -> Option<i64> {
        match self {
            Resp::String(s) => std::str::from_utf8(s.as_ref())
                .ok()
                .and_then(|s| s.parse::<i64>().ok()),
            Resp::Error(s) => std::str::from_utf8(s.as_ref())
                .ok()
                .and_then(|s| s.parse::<i64>().ok()),
            Resp::Integer(s) => Some(*s),
            Resp::BulkString(s) => std::str::from_utf8(s.as_ref())
                .ok()
                .and_then(|s| s.parse::<i64>().ok()),
            _ => None,
        }
    }

    pub fn to_str(&self) -> String {
        String::new()
    }

    pub fn to_u128(&self) -> Option<u128> {
        match self {
            Resp::String(s) => std::str::from_utf8(s.as_ref())
                .ok()
                .and_then(|s| s.parse::<u128>().ok()),
            Resp::Error(s) => std::str::from_utf8(s.as_ref())
                .ok()
                .and_then(|s| s.parse::<u128>().ok()),
            Resp::Integer(s) => Some(*s as u128),
            Resp::BulkString(s) => std::str::from_utf8(s.as_ref())
                .ok()
                .and_then(|s| s.parse::<u128>().ok()),
            _ => None,
        }
    }

    pub fn to_set_mode(&self) -> Option<SetMode> {
        match self {
            Resp::String(b) => {
                if b.as_ref().eq_ignore_ascii_case(b"NX") {
                    Some(SetMode::Nx)
                } else if b.as_ref().eq_ignore_ascii_case(b"XX") {
                    Some(SetMode::Xx)
                } else if b.as_ref().eq_ignore_ascii_case(b"GT") {
                    Some(SetMode::Gt)
                } else if b.as_ref().eq_ignore_ascii_case(b"LT") {
                    Some(SetMode::Lt)
                } else {
                    None
                }
            }
            Resp::Error(_) => None,
            Resp::Integer(_) => None,
            Resp::BulkString(b) => {
                if b.as_ref().eq_ignore_ascii_case(b"NX") {
                    Some(SetMode::Nx)
                } else if b.as_ref().eq_ignore_ascii_case(b"XX") {
                    Some(SetMode::Xx)
                } else if b.as_ref().eq_ignore_ascii_case(b"GT") {
                    Some(SetMode::Gt)
                } else if b.as_ref().eq_ignore_ascii_case(b"LT") {
                    Some(SetMode::Lt)
                } else {
                    None
                }
            }
            Resp::NilBulk => None,
            Resp::Array(_) => None,
            Resp::NilArray => None,
        }
    }
}

impl Hash for Resp {
    fn hash<H: Hasher>(&self, state: &mut H) {
        match self {
            Resp::String(s) => s.hash(state),
            Resp::Error(s) => s.hash(state),
            Resp::Integer(s) => s.hash(state),
            Resp::BulkString(s) => s.hash(state),
            Resp::NilBulk => 0.hash(state),
            Resp::Array(s) => s.hash(state),
            Resp::NilArray => 1.hash(state),
        }
    }
}

#[derive(Debug, Clone, Eq, PartialEq, Default)]
pub struct RespArray(pub usize, pub BytesMut);

impl RespArray {
    pub fn as_bytes(&self) -> &[u8] {
        self.1.as_ref()
    }

    pub fn len(&self) -> usize {
        self.0
    }

    pub fn iter(&self) -> RespArrayIterator {
        RespArrayIterator(0, self.1.into())
    }

    pub fn get(&self, index: usize) -> Option<Resp> {
        self.iter().nth(index)
    }
}

#[derive(Debug, Clone, Eq, PartialEq, Default)]
pub struct RespArrayBuilder(pub usize, pub BytesMut);

impl RespArrayBuilder {
    pub fn push_resp(&mut self, resp: Resp) {
        self.0 += 1;
        self.1.reserve(resp.resp_len());
        match resp {
            Resp::String(bytes) => {
                let content_bytes = bytes.as_ref();
                self.1.put_slice(b"+");
                self.1.put_slice(content_bytes);
                self.1.put_slice(b"\r\n");
            }
            Resp::Error(bytes) => {
                let content_bytes = bytes.as_ref();
                self.1.put_slice(b"-");
                self.1.put_slice(content_bytes);
                self.1.put_slice(b"\r\n");
            }
            Resp::Integer(b) => {
                let content_bytes = format!("{}", b);
                self.1.put_slice(b":");
                self.1.put_slice(content_bytes.as_ref());
                self.1.put_slice(b"\r\n");
            }
            Resp::BulkString(bytes) => {
                let len_str = format!("{}", bytes.len());
                let len_bytes = len_str.as_bytes();
                let content_bytes = bytes.as_ref();
                self.1.put_slice(b"$");
                self.1.put_slice(len_bytes);
                self.1.put_slice(b"\r\n");
                self.1.put_slice(content_bytes);
                self.1.put_slice(b"\r\n");
            }
            Resp::NilBulk => self.1.put_slice(b"$-1\r\n"),
            Resp::Array(arr) => {
                let len_str = format!("{}", arr.0);
                let len_bytes = len_str.as_bytes();
                self.1.put_slice(b"*");
                self.1.put_slice(len_bytes);
                self.1.put_slice(b"\r\n");
                for b in arr {
                    self.push_resp(b);
                }
            }
            Resp::NilArray => self.1.put_slice(b"*-1\r\n"),
        }
    }

    pub fn build(self) -> RespArray {
        RespArray(self.0, self.1)
    }
}

impl Hash for RespArray {
    fn hash<H: Hasher>(&self, state: &mut H) {
        self.0.hash(state);
        self.1.hash(state);
    }
}

impl IntoIterator for RespArray {
    type Item = Resp;

    type IntoIter = RespArrayIterator;

    fn into_iter(self) -> Self::IntoIter {
        RespArrayIterator(0, self.1)
    }
}

impl IntoIterator for &RespArray {
    type Item = Resp;

    type IntoIter = RespArrayIterator;

    fn into_iter(self) -> Self::IntoIter {
        RespArrayIterator(0, self.1.clone())
    }
}

impl IntoIterator for &mut RespArray {
    type Item = Resp;

    type IntoIter = RespArrayIterator;

    fn into_iter(self) -> Self::IntoIter {
        RespArrayIterator(0, self.1.clone())
    }
}

pub struct RespArrayIterator(usize, BytesMut);

impl Iterator for RespArrayIterator {
    type Item = Resp;

    fn next(&mut self) -> Option<Self::Item> {
        if self.0 == self.1.len() {
            return None;
        }
        let resp_length = get_resp_length(&self.1)
            .expect("we should have checked that every construction of RESP is valid");
        let resp_bytes = &self.1.as_ref()[self.0..self.0 + resp_length];
        self.0 += resp_length;
        Some(
            cast_bytes_to_resp(resp_bytes)
                .expect("we should have checked that every construction of RESP is valid"),
        )
    }
}

#[derive(Debug)]
pub enum RespParseError {
    // Cannot find CRLF at index
    NotEnoughBytes,
    // Incorrect format detected
    IncorrectFormat,
    Other(Box<dyn std::error::Error + Sync + Send>),
}

impl std::fmt::Display for RespParseError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            RespParseError::NotEnoughBytes => write!(f, "No enough bytes"),
            RespParseError::IncorrectFormat => write!(f, "Incorrect format"),
            RespParseError::Other(err) => write!(f, "{}", err),
        }
    }
}

impl std::error::Error for RespParseError {}

impl From<std::str::Utf8Error> for RespParseError {
    fn from(from: std::str::Utf8Error) -> Self {
        Self::Other(Box::new(from))
    }
}

impl From<std::num::ParseIntError> for RespParseError {
    fn from(from: std::num::ParseIntError) -> Self {
        Self::Other(Box::new(from))
    }
}

impl From<std::io::Error> for RespParseError {
    fn from(from: std::io::Error) -> Self {
        Self::Other(Box::new(from))
    }
}

/// Given a bytes, calculating the length of the first RESP in that byte
/// Includes the first marker symbol and the ending CRLF.
pub fn get_resp_length(input: &[u8]) -> Result<usize, RespParseError> {
    if input.is_empty() {
        Err(RespParseError::NotEnoughBytes)
    } else {
        let marker_byte = input[0];
        let payload_bytes = &input[1..];
        let length = match marker_byte {
            b'+' => get_rss_length(payload_bytes)? + 1,
            b':' => get_ri_length(payload_bytes)? + 1,
            b'$' => get_rbs_length(payload_bytes)? + 1,
            b'*' => get_ra_length(payload_bytes)? + 1,
            b'-' => get_res_length(payload_bytes)? + 1,
            _ => get_rss_length(input)?,
        };
        // Compensate for the first byte
        Ok(length)
    }
}

pub fn cast_bytes_to_resp(input: &u) -> RespResult {
    let marker_byte = input[0];
    let resp = match marker_byte {
        b'+' => Resp::String(input.slice(1..input.len() - 2)),
        b':' => {
            let content_bytes = input.slice(1..input.len() - 2);
            let content_value = std::str::from_utf8(content_bytes.as_ref())?.parse::<i64>()?;
            Resp::Integer(content_value)
        }
        b'$' => cast_bytes_to_rbs(input.slice(1..))?,
        b'*' => cast_bytes_to_ra(input.slice(1..))?,
        b'-' => Resp::Error(input.slice(1..input.len() - 2)),
        _ => Resp::String(input.slice(..input.len() - 2)),
    };
    Ok(resp)
}

fn parse_everything_until_crlf(
    input: &[u8],
) -> std::result::Result<(&[u8], &[u8]), RespParseError> {
    for (index, (first, second)) in input.iter().zip(input.iter().skip(1)).enumerate() {
        if first == &CR && second == &LF {
            let before_crlf = &input[..index];
            let after_crlf = &input[index + 2..];
            return Ok((before_crlf, after_crlf));
        }
    }
    Err(RespParseError::NotEnoughBytes)
}

fn parse_everything_until_crlf_from_bytes(
    input: Bytes,
) -> std::result::Result<(Bytes, Bytes), RespParseError> {
    for (index, (first, second)) in input.iter().zip(input.iter().skip(1)).enumerate() {
        if first == &CR && second == &LF {
            let before_crlf = input.slice(..index);
            let after_crlf = input.slice(index + 2..);
            return Ok((before_crlf, after_crlf));
        }
    }
    Err(RespParseError::NotEnoughBytes)
}

// Similar to `parse_everything_until_crlf` except that we already tell it where to find the CRLF.
fn parse_everything_until_index(
    input: &[u8],
    index: usize,
) -> Result<(&[u8], &[u8]), RespParseError> {
    if input.len() <= index {
        Err(RespParseError::NotEnoughBytes)
    } else if input[index] == b'\r' && input[index + 1] == b'\n' {
        let before_crlf = &input[..index];
        let after_crlf = &input[index + 2..];
        return Ok((before_crlf, after_crlf));
    } else {
        return Err(RespParseError::IncorrectFormat);
    }
}

pub fn get_rss_length(input: &[u8]) -> RespLengthResult {
    parse_everything_until_crlf(input).map(|(x, _)| x.len() + 2)
}

pub fn get_res_length(input: &[u8]) -> RespLengthResult {
    parse_everything_until_crlf(input).map(|(x, _)| x.len() + 2)
}

pub fn get_ri_length(input: &[u8]) -> RespLengthResult {
    parse_everything_until_crlf(input).map(|(x, _)| x.len() + 2)
}

pub fn get_rbs_length(input: &[u8]) -> RespLengthResult {
    let mut count = 0;
    let (size_bytes, leftover) = parse_everything_until_crlf(input)?;
    count += size_bytes.len() + 2;

    if size_bytes == b"-1" {
        Ok(5)
    } else {
        let size = std::str::from_utf8(size_bytes)?.parse::<usize>()?;
        let (result, _) = parse_everything_until_index(leftover, size)?;
        count += result.len();
        Ok(count + 2)
    }
}

pub fn cast_bytes_to_rbs(input: Bytes) -> RespResult {
    let (size_bytes, leftover) = parse_everything_until_crlf_from_bytes(input)?;

    if size_bytes.as_ref() == b"-1" {
        Ok(Resp::NilBulk)
    } else {
        let resp = Resp::BulkString(leftover.slice(..leftover.len() - 2));
        Ok(resp)
    }
}

pub fn cast_bytes_to_ra(input: Bytes) -> RespResult {
    // Skip the first byte
    let (size_bytes, leftover) = parse_everything_until_crlf_from_bytes(input)?;

    if size_bytes.as_ref() == b"-1" {
        Ok(Resp::NilArray)
    } else {
        let size = std::str::from_utf8(size_bytes.as_ref())?.parse::<usize>()?;
        let resp = Resp::Array(RespArray(size, leftover));
        Ok(resp)
    }
}

pub fn get_ra_length(input: &[u8]) -> RespLengthResult {
    let mut count = 0;
    let (size_bytes, leftover) = parse_everything_until_crlf(input)?;

    // Add the CRLF
    count += size_bytes.len() + 2;

    if size_bytes == b"-1" {
        Ok(5)
    } else {
        let size = std::str::from_utf8(size_bytes)?.parse::<usize>()?;
        let mut left = leftover;
        for _ in 0..size {
            let length = get_resp_length(left)?;
            count += length;
            left = &left[length..];
        }
        Ok(count)
    }
}
