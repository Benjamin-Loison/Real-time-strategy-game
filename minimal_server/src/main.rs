// Require to load the JSON conf
use std::path::Path;
use std::fs::File;
use serde::Deserialize;

// Require to run the server
use std::net::SocketAddr;
use std::net::Ipv4Addr;
use std::net::IpAddr;
use std::str::FromStr;
use std::net::TcpListener;
use std::thread::spawn;

// exit the program if required
use std::process::exit;


// Structure that stores the config server config
#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct ServerConf {
	name: String,
	address: String,
	port: u16
}

// Load a config file
fn load_conf(file_name: &str) -> ServerConf {
	let json_file_path = Path::new(file_name);
	let file = File::open(json_file_path).expect("file not found");
	return serde_json::from_reader(file).expect("error while reading");
}



// Main function
fn main() {
	// Load the configuration and verbose
	let conf : ServerConf = load_conf("conf/conf.json");
	println!("[Conf] {} listens on {}:{}", conf.name, conf.address, conf.port);

	// Start the server
	let sock_addr: SocketAddr = SocketAddr::new(
		IpAddr::V4(Ipv4Addr::from_str(&conf.address.to_owned()).unwrap()),
		conf.port);
	let listener = TcpListener::bind(sock_addr).unwrap();

	// Listen to new connection and fork for each one of them
	for stream in listener.incoming() {
		match stream {
			Err(_) => println!("listen error"),
			Ok(mut stream) => {
				println!("connection from {} to {}",
						 stream.peer_addr().unwrap(),
						 stream.local_addr().unwrap());
				spawn(move|| {
					//
					exit(0);
				});
			}
		}
	}
}
