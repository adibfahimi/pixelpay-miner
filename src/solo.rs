use crate::{block::Block, tx::Tx};
use colored::*;
use serde::{Deserialize, Serialize};
use std::time::{SystemTime, UNIX_EPOCH};

#[derive(Serialize, Deserialize)]
struct NodeResp {
    pub block: Block,
    pub miner_reward: u32,
}

async fn get_block(node_uri: &str) -> String {
    reqwest::get(format!("{}/mine", node_uri))
        .await
        .unwrap()
        .text()
        .await
        .unwrap()
}

async fn send_block(block: &Block, node_uri: &str) -> Result<reqwest::Response, reqwest::Error> {
    reqwest::Client::new()
        .post(format!("{}/mine", node_uri))
        .header("Content-Type", "application/json")
        .body(serde_json::to_string(&block).unwrap())
        .send()
        .await
}

pub async fn solo_mine(node_uri: &str, miner_address: &str) {
    let mut show_no_tx_to_mine = true;

    loop {
        let resp = get_block(node_uri).await;

        if resp == "No transactions to mine" {
            if show_no_tx_to_mine {
                println!(
                    "{}",
                    "No transactions to mine, waiting for 10 seconds..."
                        .to_string()
                        .bright_black(),
                );
                show_no_tx_to_mine = false;
            }
            tokio::time::sleep(tokio::time::Duration::from_secs(10)).await;
            continue;
        }

        println!("mining...");
        let start_time = SystemTime::now();

        let node_resp: NodeResp = serde_json::from_str(&resp).unwrap();

        let mut block = node_resp.block.clone();
        let timestamp = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs();

        let coin_base_tx = Tx {
            sender: "".to_string(),
            receiver: miner_address.to_string(),
            amount: node_resp.miner_reward,
            signature: "".to_string(),
            hash: "".to_string(),
            timestamp,
        };
        block.data.push(coin_base_tx);
        block.merkle_root = block.calculate_merkle_root();
        block.mine();

        println!();
        println!("Block mined: {}", block.hash.bright_green());
        println!("Nonce: {}", block.nonce.to_string().bright_red());
        println!(
            "Time taken to mine: {} seconds",
            SystemTime::now()
                .duration_since(start_time)
                .unwrap()
                .as_secs()
                .to_string()
                .bright_yellow()
        );

        show_no_tx_to_mine = true;

        for _ in 0..10 {
            let resp = send_block(&block, node_uri).await;
            if resp.is_ok() {
                break;
            }
            tokio::time::sleep(tokio::time::Duration::from_secs(10)).await;
        }

        println!("Block sent to node");
    }
}
