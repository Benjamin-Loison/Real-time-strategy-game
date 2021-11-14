extern crate bufstream;
extern crate regex;
mod lib;
//mod prepare;

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
use std::sync::mpsc::{Sender, Receiver};
pub use self::lib::{EntityRegistry, Position, Type, Entity};
//pub use self::prepare::{position_update, inside_or_not, prepare_map, send_map_data, query_after_update};


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

// update all the data of playerlist
pub fn position_update(stream: TcpStream, id: &str, client_x: &str, client_y: &str, playerlist: Arc<Mutex<EntityRegistry>>) {
	let client_latest_position = Position::displacement(client_x.parse::<i64>().unwrap(), client_y.parse::<i64>().unwrap());
	let mut latest_playerlist = playerlist.lock().unwrap();
	let owner: String = stream.peer_addr().unwrap().to_string();
	latest_playerlist.update("human", id.parse::<i64>().unwrap(), client_latest_position, owner);
}

// check if the new position is inide the scope of fixed position with width and height
pub fn inside_or_not(fixed_entity: Entity, new_entity: Entity, width: &str, height: &str) -> bool {
	let w = width.parse::<i64>().unwrap();
	let h = height.parse::<i64>().unwrap();
	let x_ok: bool = ((fixed_entity.position.x - new_entity.position.x).abs() <= w);
	let y_ok: bool = ((fixed_entity.position.y - new_entity.position.y).abs() <= h);
	return (x_ok & y_ok);
}

//generate the left upper vertex of the sub-rectangle
pub fn generate_vertex(client_x: &str, client_y: &str, width: &str, height: &str) -> String{
	let vertex_x = client_x.parse::<i64>().unwrap() + width.parse::<i64>().unwrap();
	let vertex_y = client_y.parse::<i64>().unwrap() + height.parse::<i64>().unwrap();
	format!("{},{},\n", vertex_x.to_string(), vertex_y.to_string())
}


// prepare the string map_data from playerlist to send back to client
pub fn prepare_map(stream: TcpStream, id: &str, width: &str, height: &str, playerlist: Arc<Mutex<EntityRegistry>>) -> String {
	let mut internal_playerlist = EntityRegistry::create();
	let playerlist_unwrap = playerlist.lock().unwrap();
	let fixed_entity = playerlist_unwrap.entitylist.get(&id.parse::<i64>().unwrap()).unwrap();
	let playerlist_unwrap_clone = playerlist_unwrap.clone();
	//let owner: String = stream.peer_addr().unwrap().to_string();
	for (new_id, new_entity) in playerlist_unwrap_clone.entitylist{
		let fixed_entity_clone = fixed_entity.clone();
		let mut owner_clone = new_entity.owner.clone();
		let mut position_clone = new_entity.position.clone();
		if inside_or_not(fixed_entity_clone, new_entity, width, height) == true {
			//let owner = playerlist_unwrap.ownerlist.get(new_id).unwrap();
			internal_playerlist.update("human", new_id, position_clone, owner_clone);
		}
	}
	return internal_playerlist.provide();
}

//send map data to client after receiving query with key word "map"
pub fn send_map_data(stream: TcpStream, id: &str, options: &str, playerlist: Arc<Mutex<EntityRegistry>>) {
	let xywh: Vec<&str> = options.split(',').collect();
	let mut playerlist_clone1 = Arc::clone(&playerlist);
	let mut playerlist_clone2 = Arc::clone(&playerlist);
	let mut stream_clone1 = stream.try_clone().unwrap();
	let mut stream_clone2 = stream.try_clone().unwrap();
	position_update(stream, id, xywh[0], xywh[1], playerlist_clone1);
	let map_data: String = prepare_map(stream_clone1, id, xywh[2], xywh[3], playerlist_clone2);
	let vertex_string : String = generate_vertex(xywh[0], xywh[1], xywh[2], xywh[3]);
	let new_options = format!("{}{}",vertex_string, map_data);
	answer_client(stream_clone2, id, new_options.as_str());
}

//UNFINISHED: send new query to clients after detecting conflicts from latest update
pub fn query_after_update(stream: TcpStream, id: &str, playerlist: Arc<Mutex<EntityRegistry>>) {
	let owner: &str = &stream.peer_addr().unwrap().to_string();	
	query_client(stream, id, "get");
}


fn parse_query(stream: TcpStream, id: &str, command: &str, options: &str, playerlist: Arc<Mutex<EntityRegistry>>) {
	match command {
		"info" => {
			answer_client(stream, id, "0"); // TODO: the server id is 0
		}
        "map" => {
            send_map_data(stream, id, options, playerlist);
        }
        _ => {
			logging("parse_query", "Unknown command.");
		}
	}
}
fn parse_answer(stream: TcpStream, id: &str, command: &str, options: &str, playerlist: Arc<Mutex<EntityRegistry>>) {
	match command {
		"map" => {
			send_map_data(stream, id, options, playerlist);
		}
		_ => {
			logging("parse_query", "Unknown command.");
		}
	}
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

fn parse_incoming(stream: TcpStream, msg: &str, playerlist: Arc<Mutex<EntityRegistry>>) {
    let re_query = Regex::new(r"Q([a-zA-Z0-9]+).([a-zA-Z0-9]+):([a-zA-Z0-9,]*)").unwrap();
    let re_answer = Regex::new(r"A([a-zA-Z0-9]+).([a-zA-Z0-9]+):([a-zA-Z0-9,]*)").unwrap();
    let re_status = Regex::new(r"S([a-zA-Z0-9]+).(ok|nok)").unwrap();
    if re_query.is_match(msg) {
        let cap = re_query.captures(msg).unwrap();
        parse_query(stream, &cap[1], &cap[2], &cap[3], playerlist);
    } else if re_answer.is_match(msg) {
        let cap = re_answer.captures(msg).unwrap();
        parse_answer(stream, &cap[1], &cap[2], &cap[3], playerlist);
    } else if re_status.is_match(msg) {
        let cap = re_status.captures(msg).unwrap();
        parse_status(stream, &cap[1], &cap[2]);
    } else {
        println!("Invalid incoming message");
    }
}

fn handle_connection(stream: TcpStream, playerlist: Arc<Mutex<EntityRegistry>>) {
	let mut playerlist_clone1 = Arc::clone(&playerlist);
	let mut playerlist_clone2 = Arc::clone(&playerlist);
	let stream_clone1 = stream.try_clone().unwrap();
	let stream_clone2 = stream.try_clone().unwrap();

	//receive client's query and answer and response
	thread::spawn(move || {
		loop {
			let stream_clone1_new = stream_clone1.try_clone().unwrap();
			let mut streamreader = BufReader::new(&stream_clone1);
			let mut reads = String::new();
			streamreader.read_line(&mut reads).unwrap(); //TODO: non-blocking read
			let mut playerlist_cloned1 = Arc::clone(&playerlist_clone1);
			if reads.trim().len() != 0 {
				parse_incoming(stream_clone1_new, reads.trim(), playerlist_cloned1);
			}	
		}
	});

	//update by querying other potential clients about their new infomations
	thread::spawn(move || {
		let sleep_time: u64 = (10000) as u64;
		loop {
			thread::sleep(Duration::from_millis(sleep_time));
			let mut playerlist_cloned2 = Arc::clone(&playerlist_clone2);
			let stream_clone2_new = stream_clone2.try_clone().unwrap();
			query_after_update(stream_clone2_new, "312", playerlist_cloned2);
		}
	});
}

fn main() {
	let addr: SocketAddr = SocketAddr::from_str("192.168.30.1:80").unwrap();
	let listener = TcpListener::bind(addr).unwrap();
	let mut playerlist = EntityRegistry::create();
	let mut playerlist_new = Arc::new(Mutex::new(playerlist));


	for stream in listener.incoming() {
		match stream {
			Err(_) => println!("listen error"),
			Ok(mut stream) => {
				println!("connection from {} to {}",
						 stream.peer_addr().unwrap(),
						 stream.local_addr().unwrap());
						 
		        //let mut adr = stream.peer_addr().unwrap().port();
		        let mut playerlist_clone = Arc::clone(&playerlist_new);
				spawn(move|| {
					//let mut stream = BufStream::new(stream);
					handle_connection(stream, playerlist_clone);
				});
			}
		}
	}
}