use crossbeam_channel::{Sender};
use std::net::{TcpStream};
use regex::Regex;
use std::sync::Arc;
use std::sync::Mutex;
pub use super::lib::{EntityRegistry, Position, Type, Entity};
pub use super::process::{position_update, generate_vertex, prepare_inside_players, send_map_data};
pub use super::handle::{logging, query_client, answer_client,status_client};

pub fn parse_query(stream: TcpStream, id: &str, command: &str, options: &str, playerlist: Arc<Mutex<EntityRegistry>>, sender: Sender<Arc<Mutex<EntityRegistry>>>) {
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

pub fn parse_answer(stream: TcpStream, id: &str, options: &str, playerlist: Arc<Mutex<EntityRegistry>>, sender: Sender<Arc<Mutex<EntityRegistry>>>) {
	let original_bigframe = "10";
	let mut new_options = format!("{},{},{}", options, original_bigframe, original_bigframe);
	println!("new_options:{}", new_options);
	send_map_data(stream, id, &new_options, playerlist,"A", sender);
}

pub fn parse_status(stream: TcpStream, id: &str, options: &str) {
	match options {
		"ok" => {
			logging(id, "Status Ok");
		}
		_ => {
			logging(id, "Status not Ok");
		}
	}
}

pub fn parse_incoming(stream: TcpStream, msg: &str, playerlist: Arc<Mutex<EntityRegistry>>, sender: Sender<Arc<Mutex<EntityRegistry>>>) {
    let re_query = Regex::new(r"Q([a-zA-Z0-9]+).([a-zA-Z0-9]+):([\d,]*)").unwrap();
    let re_answer = Regex::new(r"A([a-zA-Z0-9]+).([\d,]*)").unwrap();
    let re_status = Regex::new(r"S([a-zA-Z0-9]+).(ok|nok)").unwrap();
    if re_query.is_match(msg) {
		println!("are you kiding me?");
        let cap = re_query.captures(msg).unwrap();
		println!("error here:{}",&cap[3]);
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