mod block;
mod solo;
mod tx;

#[tokio::main]
async fn main() {
    let args: Vec<String> = std::env::args().collect();
    if args.len() != 3 {
        println!("Usage: pixelpay-miner <miner_address> <node_uri>");
        std::process::exit(1);
    }

    let miner_address = &args[1];
    let node_uri = &args[2];

    solo::solo_mine(node_uri, miner_address).await;
}
