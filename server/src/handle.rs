extern crate bufstream;
extern crate crossbeam_channel;
use std::time::SystemTime;
use std::io::{Write, BufWriter};
use std::net::{TcpStream};
pub use super::lib::{EntityRegistry, Position, Type, Entity};




pub fn logging(from: &str, msg: &str) {
	println!("[{:?}]---[{}]: {}", SystemTime::now(), from, msg);
}

pub fn answer_client(stream: TcpStream, id: &str, options: &str) {
	let ans = format!("A{}.{}\n", id, options);
	let mut streamwriter = BufWriter::with_capacity(100, stream);
	streamwriter.write(ans.as_bytes());
	streamwriter.flush();
	logging("answer_client", &ans);
}

pub fn query_client(stream: TcpStream, id: &str, options: &str) {
	let ans = format!("Q{}.location:{}\n", id, options);
	let mut streamwriter = BufWriter::with_capacity(100, stream);
	streamwriter.write(ans.as_bytes());
	streamwriter.flush();
	logging("query_client", &ans);
}

pub fn status_client(stream: TcpStream, id: &str, options: &str) {
	let ans = format!("S{}.{}\n", id, options);
	let mut streamwriter = BufWriter::with_capacity(100, stream);
	streamwriter.write(ans.as_bytes());
	streamwriter.flush();
	logging("status_client", &ans);
}
