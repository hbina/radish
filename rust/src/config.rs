use std::collections::HashMap;

use crate::resp::Resp;

pub fn get_default_config() -> HashMap<Resp, Resp> {
    let mut result = HashMap::with_capacity(100);

    result.insert(
        Resp::rbs_from_static_bytes(b"save"),
        Resp::rbs_from_static_bytes(b"900 1 300 10 60 10000"),
    );
    // "appendonly":                      "no",
    result.insert(
        Resp::rbs_from_static_bytes(b"appendonly"),
        Resp::rbs_from_static_bytes(b"no"),
    );


    result
}
