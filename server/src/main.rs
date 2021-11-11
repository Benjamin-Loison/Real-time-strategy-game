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

fn parse_query(id: &str, command: &str, options: &str) {
}
fn parse_answer(id: &str, command: &str, options: &str) {
}
fn parse_status(id: &str, options: &str) {
}


fn parse_incoming(msg: &str) {
    let re_query = Regex::new(r"Q([a-zA-Z0-9]+).([a-zA-Z0-9]+):([a-zA-Z0-9,]*)").unwrap();
    let re_answer = Regex::new(r"A([a-zA-Z0-9]+).([a-zA-Z0-9]+):([a-zA-Z0-9,]*)").unwrap();
    let re_status = Regex::new(r"S([a-zA-Z0-9]+).(ok|nok)").unwrap();
    if re_query.is_match(msg) {
        let cap = re_query.captures(msg).unwrap();
        parse_query(&cap[1], &cap[2], &cap[3]);
    } else if re_answer.is_match(msg) {
        let cap = re_answer.captures(msg).unwrap();
        parse_answer(&cap[1], &cap[2], &cap[3]);
    } else if re_status.is_match(msg) {
        let cap = re_query.captures(msg).unwrap();
        parse_status(&cap[1], &cap[2]);
    } else {
        println!("Invalid incoming message");
    }
}

fn handle_connection(stream: &mut BufStream<TcpStream>) {
	loop {
		stream.write(b" > ").unwrap();
		stream.flush().unwrap();

		let mut reads = String::new();
		stream.read_line(&mut reads).unwrap(); //TODO: non-blocking read
		if reads.trim().len() != 0 {
            parse_incoming(reads.trim());
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
