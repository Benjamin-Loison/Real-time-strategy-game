extern crate bufstream;
extern crate regex;
extern crate crossbeam_channel;
mod lib;
//mod prepare;

//use std::sync::mpsc::channel;
use crossbeam_channel::{unbounded, Sender, Receiver};
use std::time::Duration;
use std::str::FromStr;
use std::io::{Write, BufWriter, BufReader};
use std::net::{TcpListener, TcpStream};
use std::net::SocketAddr;
use std::thread::spawn;
use std::thread;
use bufstream::BufStream;
use std::io::BufRead;
use regex::Regex;
use std::sync::Arc;
use std::sync::Mutex;
use std::sync::mpsc;
//use std::sync::mpsc::{Sender, Receiver};
pub use self::lib::{EntityRegistry, Position, Type, Entity};
//pub use self::prepare::{position_update, inside_or_not, prepare_inside_players, send_map_data, query_after_update};


pub fn logging(from: &str, msg: &str) {
	println!("[{}]: {}", from, msg);
}

pub fn answer_client(stream: TcpStream, id: &str, options: &str) {
	let ans = format!("A{}.{}\n", id, options);
	let mut streamwriter = BufWriter::with_capacity(100, stream);
	streamwriter.write(ans.as_bytes());
	streamwriter.flush();
}

pub fn query_client(stream: TcpStream, id: &str, options: &str) {
	let ans = format!("Q{}.location:{}\n", id, options);
	let mut streamwriter = BufWriter::with_capacity(100, stream);
	streamwriter.write(ans.as_bytes());
	streamwriter.flush();
}

pub fn status_client(stream: TcpStream, id: &str, options: &str) {
	let ans = format!("S{}.{}\n", id, options);
	let mut streamwriter = BufWriter::with_capacity(100, stream);
	streamwriter.write(ans.as_bytes());
	streamwriter.flush();
}

// update all the data of playerlist
pub fn position_update(stream: TcpStream, id: &str, client_x: &str, client_y: &str, width: &str, height: &str, playerlist: Arc<Mutex<EntityRegistry>>, sender: Sender<Arc<Mutex<EntityRegistry>>>) {
	let client_latest_position = Position::displacement(client_x.parse::<i64>().unwrap(), client_y.parse::<i64>().unwrap());
	let playerlist_to_send = Arc::clone(&playerlist);
	let mut latest_playerlist = playerlist.lock().unwrap();
	let owner: String = stream.peer_addr().unwrap().to_string();
	latest_playerlist.update("human", id.parse::<i64>().unwrap(), client_latest_position, owner, width.parse::<i64>().unwrap(), height.parse::<i64>().unwrap());
	let sender_clone = sender.clone();
	sender.send(playerlist_to_send).unwrap();
}

//generate the left upper vertex of the sub-rectangle
pub fn generate_vertex(client_x: &str, client_y: &str, width: &str, height: &str) -> String{
	let vertex_x = client_x.parse::<i64>().unwrap() + width.parse::<i64>().unwrap();
	let vertex_y = client_y.parse::<i64>().unwrap() + height.parse::<i64>().unwrap();
	format!("{},{}", vertex_x.to_string(), vertex_y.to_string())
}

//prepare the infomation of other players who get into the big frame of current player
pub fn prepare_inside_players(stream: TcpStream, owner: String, playerlist: Arc<Mutex<EntityRegistry>>, thestring: &str) {
	//println!("second");
	let mut playerlist_unwrap = playerlist.lock().unwrap();
	//println!("second");
	let mut ownershiplist = playerlist_unwrap.ownershiplist.get(&owner).unwrap();
	for id in ownershiplist.iter() {
		let stream_clone1 = stream.try_clone().unwrap();
		match thestring {
			"set" => {
				let list = playerlist_unwrap.players_nearby(*id);
				//let stream_clone1_clone = stream_clone1.try_clone().unwrap();
				for (ids, entitys) in list {
					let stream_clone1_clone = stream_clone1.try_clone().unwrap();
					let mut output = String::new();
					output.push_str(&format!("set ids:{} {}",ids, entitys.position.provide()));
					query_client(stream_clone1_clone, &id.to_string(), &output);
				}
			}
			_ => {
				query_client(stream_clone1, &id.to_string(), "get");
			}
		}
		/*
		let mut output = String::from("set ");
		match thestring {
			"set" => {
				output.push_str(&playerlist_unwrap.players_nearby(*id));
				output.push_str("\n");
				query_client(stream_clone1, &id.to_string(), &output);
			}
			_ => {
				query_client(stream_clone1, &id.to_string(), "get");
			}
		}
		*/
	}
}

//send map data to client after receiving query with key word "map"
pub fn send_map_data(stream: TcpStream, id: &str, options: &str, playerlist: Arc<Mutex<EntityRegistry>>, msg: &str, sender: Sender<Arc<Mutex<EntityRegistry>>>) {
	let xywh: Vec<&str> = options.split(',').collect();
	let mut playerlist_clone1 = Arc::clone(&playerlist);
	let mut playerlist_clone2 = Arc::clone(&playerlist);
	let mut stream_clone2 = stream.try_clone().unwrap();
	position_update(stream, id, xywh[0], xywh[1], xywh[2], xywh[3], playerlist_clone1, sender);
	let playerlist_clone2_unwrap = playerlist_clone2.lock().unwrap();
	let playerlist_clone2_unwrap_clone = playerlist_clone2_unwrap.clone();
	let map_data: String = playerlist_clone2_unwrap_clone.bigframe(id.parse::<i64>().unwrap());
	let vertex_string : String = generate_vertex(xywh[0], xywh[1], xywh[2], xywh[3]);
	let new_options = format!("{}{}",vertex_string, map_data);
	match msg {
		"Q" => {
			answer_client(stream_clone2, id, new_options.as_str());
		}
		_ => {
			status_client(stream_clone2, id, "OK");
		}
	}
}

//UNFINISHED: send new query to clients after detecting conflicts from latest update
pub fn query_after_update(stream: TcpStream, playerlist: Arc<Mutex<EntityRegistry>>, thestring: &str) {
	let owner: &str = &stream.peer_addr().unwrap().to_string();
	let stream_clone2 = stream.try_clone().unwrap();
	let stream_clone1 = stream.try_clone().unwrap();
	//println!("first");
	prepare_inside_players(stream_clone1, owner.to_string(), playerlist, thestring);
	//println!("lets take a look:\n{}", set);
	//println!("first");
	
	//query_client(stream_clone1, id, &set);
	//println!("first");
}


fn parse_query(stream: TcpStream, id: &str, command: &str, options: &str, playerlist: Arc<Mutex<EntityRegistry>>, sender: Sender<Arc<Mutex<EntityRegistry>>>) {
	match command {
		"info" => {
			let stream_clone1 = stream.try_clone().unwrap();
			let stream_clone2 = stream.try_clone().unwrap();
			answer_client(stream_clone1, id, "0"); // TODO: the server id is 0
			query_client(stream_clone2, id, "get");
		}
        "map" => {
            send_map_data(stream, id, options, playerlist,"Q", sender);
        }
        _ => {
			logging("parse_query", "Unknown command.");
		}
	}
}
fn parse_answer(stream: TcpStream, id: &str, options: &str, playerlist: Arc<Mutex<EntityRegistry>>, sender: Sender<Arc<Mutex<EntityRegistry>>>) {
	let original_bigframe = "10";
	let mut new_options = format!("{},{},{}", options, original_bigframe, original_bigframe);
	println!("new_options:{}", new_options);
	send_map_data(stream, id, &new_options, playerlist,"A", sender);
}

fn parse_status(stream: TcpStream, id: &str, options: &str) {
	match options {
		"ok" => {
			logging(id, "Status Ok");
		}
		_ => {
			logging(id, "Status not Ok");
		}
	}
}

fn parse_incoming(stream: TcpStream, msg: &str, playerlist: Arc<Mutex<EntityRegistry>>, sender: Sender<Arc<Mutex<EntityRegistry>>>) {
    let re_query = Regex::new(r"Q([a-zA-Z0-9]+).([a-zA-Z0-9]+):([a-zA-Z0-9,]*)").unwrap();
    let re_answer = Regex::new(r"A([a-zA-Z0-9]+).([a-zA-Z0-9,]*)").unwrap();
    let re_status = Regex::new(r"S([a-zA-Z0-9]+).(ok|nok)").unwrap();
    if re_query.is_match(msg) {
        let cap = re_query.captures(msg).unwrap();
        parse_query(stream, &cap[1], &cap[2], &cap[3], playerlist, sender);
    } else if re_answer.is_match(msg) {
        let cap = re_answer.captures(msg).unwrap();
        parse_answer(stream, &cap[1], &cap[2], playerlist, sender);
    } else if re_status.is_match(msg) {
        let cap = re_status.captures(msg).unwrap();
        parse_status(stream, &cap[1], &cap[2]);
    } else {
        println!("Invalid incoming message");
    }
}

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
			//println!("the message I receive: {}", reads);
			let mut playerlist_cloned1 = Arc::clone(&playerlist_clone1);
			if reads.trim().len() != 0 {
				parse_incoming(stream_clone1_new, reads.trim(), playerlist_cloned1, sender_thread);
			}	
		}
	});

	thread::spawn(move || {
		let sleep_time: u64 = (10000) as u64;
		loop {
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
			//let stream_clone3 = stream.try_clone().unwrap();
			//println!("maybe error here:");
			let playerlist3: Arc<Mutex<EntityRegistry>> = receiver_thread.recv().unwrap();
			//println!("maybe error here:");
            let playerlist_clone3 = Arc::clone(&playerlist3);
			let stream_clone31 = stream_clone3.try_clone().unwrap();
			query_after_update(stream_clone31, playerlist_clone3, "set");
		}
	});
}

fn main() {
	let addr: SocketAddr = SocketAddr::from_str("138.231.144.134:80").unwrap();
	let listener = TcpListener::bind(addr).unwrap();
	let mut playerlist = EntityRegistry::create(10);
	let mut playerlist_new = Arc::new(Mutex::new(playerlist));
    let (sender, receiver) = unbounded();

	for stream in listener.incoming() {
		match stream {
			Err(_) => println!("listen error"),
			Ok(mut stream) => {
				println!("connection from {} to {}",
						 stream.peer_addr().unwrap(),
						 stream.local_addr().unwrap());
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