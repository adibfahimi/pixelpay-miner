use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Tx {
    pub sender: String,
    pub receiver: String,
    pub amount: u32,
    pub signature: String,
    pub hash: String,
    pub timestamp: u64,
}
