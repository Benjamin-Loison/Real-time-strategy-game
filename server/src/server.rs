extern crate bufstream;

use std::str::FromStr;
use std::io::Write;
use std::net::{TcpListener, TcpStream};
use std::net::SocketAddr;
use std::thread::spawn;
use bufstream::BufStream;
use std::io::BufRead;
use std::sync::{Arc,RwLock};
use std::sync::mpsc;
use std::sync::mpsc::{Sender, Receiver};

fn handle_connection(stream: &mut BufStream<TcpStream>) {
	loop {
		stream.write(b" > ").unwrap();
		stream.flush().unwrap();

		let mut reads = String::new();
		stream.read_line(&mut reads).unwrap(); //TODO: non-blocking read
		if reads.trim().len() != 0 {
			let mut answer = String::new();
			match &(reads.trim()) as &str {
				"Hello" => {
					answer = "Hello there!\n".to_string(); }
				_ => {
					answer = reads; }
			}
			stream.write(answer.as_bytes()).unwrap();
			stream.flush().unwrap();
		}
	}
}

fn main() {
	let addr: SocketAddr = SocketAddr::from_str("127.0.0.1:8888").unwrap();
	let listener = TcpListener::bind(addr).unwrap();

	for stream in listener.incoming() {
		match stream {
			Err(_) => println!("listen error"),
			Ok(mut stream) => {
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
