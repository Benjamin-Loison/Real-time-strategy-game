extern crate bufstream;
extern crate regex;
extern crate crossbeam_channel;
mod lib;
mod process;
mod handle;
pub mod parse;

use crossbeam_channel::{unbounded, Sender, Receiver};
use std::time::Duration;
use std::str::FromStr;
use std::io::{BufReader};
use std::net::{TcpListener, TcpStream};
use std::net::SocketAddr;
use std::thread::spawn;
use std::thread;
use bufstream::BufStream;
use std::io::BufRead;
use std::sync::Arc;
use std::sync::Mutex;
pub use self::lib::{EntityRegistry, Position, Type, Entity};
pub use self::process::{position_update, generate_vertex, prepare_inside_players, send_map_data, query_after_update, query_client, logging, answer_client};
pub use self::parse::{parse_query, parse_answer, parse_incoming, parse_status};



fn handle_connection(stream: TcpStream, playerlist: Arc<Mutex<EntityRegistry>>, receiver: Receiver<Arc<Mutex<EntityRegistry>>>, 
sender: Sender<Arc<Mutex<EntityRegistry>>>) {
	let mut playerlist_clone1 = Arc::clone(&playerlist);
	let mut playerlist_clone2 = Arc::clone(&playerlist);
	let stream_clone1 = stream.try_clone().unwrap();
	let stream_clone2 = stream.try_clone().unwrap();
	let stream_clone3 = stream.try_clone().unwrap();
	let receiver_thread = receiver.clone();

	
	//receive client's query and answer and response
	thread::spawn(move || {
		//let sleep_time: u64 = (1000) as u64;
		loop {
			let sender_thread = sender.clone();
			//thread::sleep(Duration::from_millis(sleep_time));
			let stream_clone1_new = stream_clone1.try_clone().unwrap();
			let mut streamreader = BufReader::new(&stream_clone1);
			let mut reads = String::new();
			streamreader.read_line(&mut reads).unwrap(); //TODO: non-blocking read
			println!("the message I receive: {}", reads);
			let mut playerlist_cloned1 = Arc::clone(&playerlist_clone1);
			if reads.trim().len() != 0 {
				parse_incoming(stream_clone1_new, reads.trim(), playerlist_cloned1, sender_thread);
			}	
		}
	});

	//query to all the clients the latest info every 5 seconds
	thread::spawn(move || {
		let sleep_time: u64 = (10000) as u64;
		loop {
			println!("we've get in");
			thread::sleep(Duration::from_millis(sleep_time));
			let stream_clone21 = stream_clone2.try_clone().unwrap();
			let mut playerlist_clone21 = playerlist_clone2.clone();
			query_after_update(stream_clone21, playerlist_clone21, "get");
		}
	});

	//update by querying other potential clients about their new infomations
	thread::spawn(move || {
		let sleep_time: u64 = (10000) as u64;
		loop {
			thread::sleep(Duration::from_millis(sleep_time));
			let playerlist3: Arc<Mutex<EntityRegistry>> = receiver_thread.recv().unwrap();
            let playerlist_clone3 = Arc::clone(&playerlist3);
			let stream_clone31 = stream_clone3.try_clone().unwrap();
			//init_msg_sending(stream_clone31);
			query_after_update(stream_clone31, playerlist_clone3, "set");
		}
	});
}

fn main() {
	let addr: SocketAddr = SocketAddr::from_str("127.0.0.1:80").unwrap();
	let listener = TcpListener::bind(addr).unwrap();
	let mut playerlist = EntityRegistry::create(10);
	let mut playerlist_new = Arc::new(Mutex::new(playerlist));
    let (sender, receiver) = unbounded();

	for stream in listener.incoming() {
		match stream {
			Err(_) => println!("listen error"),
			Ok(mut stream) => {
				let mut connection = format!("connection from {} to {}",
				         stream.peer_addr().unwrap(),
				         stream.local_addr().unwrap());
				logging("rust server", &connection);
		        let mut playerlist_clone1 = Arc::clone(&playerlist_new);
				let mut playerlist_clone2 = Arc::clone(&playerlist_new);
				let receiver_stream = receiver.clone();
				let sender_stream = sender.clone();
				sender_stream.send(playerlist_clone1);
				spawn(move|| {
					handle_connection(stream, playerlist_clone2, receiver_stream, sender_stream);
				});
			}
		}
	}
}