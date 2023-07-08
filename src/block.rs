use crate::tx::Tx;
use crypto::digest::Digest;

#[derive(Debug, Clone, serde::Serialize, serde::Deserialize)]
pub struct Block {
    pub index: u32,
    pub timestamp: u64,
    pub hash: String,
    pub prev_hash: String,
    pub data: Vec<Tx>,
    pub merkle_root: String,
    pub difficulty: u32,
    pub nonce: u32,
}

impl Block {
    pub fn calculate_hash(&self) -> String {
        let mut hasher = crypto::sha2::Sha256::new();
        let data = format!(
            "{}{}{}{}{}",
            self.index, self.timestamp, self.prev_hash, self.difficulty, self.nonce
        );
        hasher.input_str(&data);
        hasher.result_str()
    }

    pub fn calculate_merkle_root(&self) -> String {
        let mut hasher = crypto::sha2::Sha256::new();
        let mut data = "".to_string();
        for tx in &self.data {
            data.push_str(&tx.hash);
        }
        hasher.input_str(&data);
        hasher.result_str()
    }

    pub fn mine(&mut self) {
        loop {
            self.hash = self.calculate_hash();
            if self.hash.starts_with(&"0".repeat(self.difficulty as usize)) {
                break;
            }
            self.nonce += 1;
        }
    }
}
