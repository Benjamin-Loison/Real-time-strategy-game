extern crate bufstream;
extern crate regex;

use std::str::FromStr;
use std::io::Write;
use std::net::{TcpListener, TcpStream};
use std::net::SocketAddr;
use std::thread::spawn;
use bufstream::BufStream;
use std::io::BufRead;
use regex::Regex;

fn logging(from: &str, msg: &str) {
	println!("[{}]: {}", from, msg);
}

fn answer_client(stream: &mut BufStream<TcpStream>, id: &str, option: &str) {
	let ans = format!("A{}.{}\n", id, option);
	stream.write(ans.as_bytes());
	stream.flush();
}

fn parse_query(stream: &mut BufStream<TcpStream>, id: &str, command: &str, options: &str) {
	match command {
		"info" => {
			answer_client(stream, id, "MyId"); // TODO
		}
		_ => {
			logging("parse_query", "Unknown command.");
		}
	}
}
fn parse_answer(stream: &mut BufStream<TcpStream>, id: &str, command: &str, options: &str) {
	// TODO
}

fn parse_status(stream: &mut BufStream<TcpStream>, id: &str, options: &str) {
	match options {
		"ok" => {
			logging(id, "Status Ok");
		}
		_ => {
			logging(id, "Status not Ok");
		}
	}
}

fn parse_incoming(stream: &mut BufStream<TcpStream>, msg: &str) {
    let re_query = Regex::new(r"Q([a-zA-Z0-9]+).([a-zA-Z0-9]+):([a-zA-Z0-9,]*)").unwrap();
    let re_answer = Regex::new(r"A([a-zA-Z0-9]+).([a-zA-Z0-9]+):([a-zA-Z0-9,]*)").unwrap();
    let re_status = Regex::new(r"S([a-zA-Z0-9]+).(ok|nok)").unwrap();
    if re_query.is_match(msg) {
        let cap = re_query.captures(msg).unwrap();
        parse_query(stream, &cap[1], &cap[2], &cap[3]);
    } else if re_answer.is_match(msg) {
        let cap = re_answer.captures(msg).unwrap();
        parse_answer(stream, &cap[1], &cap[2], &cap[3]);
    } else if re_status.is_match(msg) {
        let cap = re_status.captures(msg).unwrap();
        parse_status(stream, &cap[1], &cap[2]);
    } else {
        println!("Invalid incoming message");
    }
}

fn handle_connection(stream: &mut BufStream<TcpStream>) {
	loop {
		let mut reads = String::new();
		stream.read_line(&mut reads).unwrap(); //TODO: non-blocking read
		if reads.trim().len() != 0 {
			parse_incoming(stream, reads.trim());
		}
	}
}

fn main() {
	let addr: SocketAddr = SocketAddr::from_str("127.0.0.1:8888").unwrap();
	let listener = TcpListener::bind(addr).unwrap();
    
	for stream in listener.incoming() {
		match stream {
			Err(_) => println!("listen error"),
			Ok(stream) => {
				println!("connection from {} to {}",
						 stream.peer_addr().unwrap(),
						 stream.local_addr().unwrap());
				spawn(move|| {
					let mut stream = BufStream::new(stream);
					handle_connection(&mut stream);
				});
			}
		}
	}
}
