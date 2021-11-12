extern crate bufstream;
mod lib;

use std::str::FromStr;
use std::io::Write;
use std::net::{TcpListener, TcpStream};
use std::net::SocketAddr;
use std::thread::spawn;
use bufstream::BufStream;
use std::io::BufRead;
use std::sync::{Arc,RwLock};
use std::sync::Mutex;
use std::sync::mpsc;
use std::sync::mpsc::{Sender, Receiver};
pub use self::lib::{Message,PositionRegistry,Position};

fn string_to_static_str(s: String) -> &'static str {
    Box::leak(s.into_boxed_str())
}

fn processor(read: &str) -> Message{
    let k = read.chars().next().unwrap();
    println!("The message we get:");
    match k {
        'Q' => {
            let v1: Vec<&str> = read.split('.').collect();
            let value: &str = v1[0];
            let mut r = value.to_string();
            r.remove(0);
            let id: &str = &r[..];
            let v2: Vec<&str> = v1[1].split(':').collect();
            let thetype:&'static str = string_to_static_str(v2[0].to_string());
            let option:&'static str = string_to_static_str(v2[1].trim().to_string());
            let msg1 = Message::Query{ID:id.parse::<u32>().unwrap(), ttype:thetype, option:option};
            msg1.call();
            return msg1;
        }
        'A' => {
            let v1: Vec<&str> = read.split(':').collect();
            let value: &str = v1[0];
            let mut r = value.to_string();
            r.remove(0);
            let id: &str = &r[..];
            let feedback:&'static str = string_to_static_str(v1[1].trim().to_string());
            let msg2 = Message::Answer{ID:id.parse::<u32>().unwrap(),feedback:feedback};
            msg2.call();
            return msg2;
        }
        'S' => {
            let v1: Vec<&str> = read.split(':').collect();
            let value: &str = v1[0];
            let mut r = value.to_string();
            r.remove(0);
            let id: &str = &r[..];
            let bool1: &str = v1[1].trim();
            let feedback:bool = bool1.parse::<bool>().unwrap();
            let msg3 = Message::Status{ID:id.parse::<u32>().unwrap(),feedback1:feedback};
            msg3.call();
            return msg3;
        }
        _ => {
            let msg4 = Message::None;
            return msg4;
        }
    }
}
    
fn client(msg: Message, port: u32, mut back1: &mut String, playerlist: Arc<Mutex<PositionRegistry>>) {
    match msg {
        Message::Query{ID, ttype, option} => {
            let mut list = playerlist.lock().unwrap();
            let prov = list.provide();
            *back1 = format!{"A{}:{}", port, prov};
        }
        Message::Answer{ID, feedback} => {
            let back2:&str = "OK";
            *back1 = format!{"S{}:{}", port, back2};
        }
        Message::Status{ID, feedback1} => {
        }
        Message::None => {
            let back2:&str = "what are you asking for?";
            *back1 = format!{"Q{}.{}:{}", port, "ask", back2};
        }    
    }
}

fn pos_get(pos: &str) -> Position{
    let v1: Vec<&str> = pos.split(',').collect();
    let v2: Vec<&str> = v1[0].split('(').collect();
    let v3: Vec<&str> = v1[1].split(')').collect();
    let x2 = v2[1];
    let y2 = v3[0];
    let x1:u32 = x2.parse::<u32>().unwrap();
    let y1:u32 = y2.parse::<u32>().unwrap();
    let pos1:Position = Position::displacement(x1, y1);   
    pos1.show();
    pos1
}

fn get(msg: Message) -> Result<Position, &'static str> {
    match msg {
        Message::Query{ID, ttype, option} => {
            let new_pos = pos_get(&option);
            return Ok(new_pos);
        }
        Message::Answer{ID, feedback} => {
            let new_pos = pos_get(&feedback);
            return Ok(new_pos);
        }
        Message::Status{ID, feedback1} => {
            return Err("none");
        }
        Message::None => {
            return Err("none");
        }    
    }
}


fn handle_connection(stream: &mut BufStream<TcpStream>, playerlist: Arc<Mutex<PositionRegistry>>, adr: u32) {
	loop {
	    
		stream.write(b" > ").unwrap();
		stream.flush().unwrap();

		let mut reads = String::new();
		stream.read_line(&mut reads).unwrap(); //TODO: non-blocking read
		
		if reads.trim().len() != 0 {
		
			let mut ans2 = processor(&reads);
			let mut pos5 = get(ans2);
			let mut playerlist2 = Arc::clone(&playerlist);
			match pos5 {
			    Ok(pos) => {
			        let mut list = playerlist2.lock().unwrap();
			        list.update(adr, pos);
			        list.show();
			    }
			    Err(e) => {
			        let mut list = playerlist2.lock().unwrap();
			        list.show();
			    }
			}
			
			let mut ans3 = format!("sdfb").to_string();
			client(ans2,123, &mut ans3, playerlist2);
			//println!("back2:{}",k);
			/*
			match &(reads.trim()) as &str {
				"Hello" => {
					answer = "Hello there!\n".to_string(); }
				_ => {
					answer = reads; }
			}
			*/
			stream.write(ans3.as_bytes()).unwrap();
			stream.flush().unwrap();
		}
		
		
	}
}

fn main() {
	let addr: SocketAddr = SocketAddr::from_str("127.0.0.1:8888").unwrap();
	let listener = TcpListener::bind(addr).unwrap();
	let mut playerlist = PositionRegistry::create();
	let mut playerlist1 = Arc::new(Mutex::new(playerlist));

	for stream in listener.incoming() {
		match stream {
			Err(_) => println!("listen error"),
			Ok(mut stream) => {
				println!("connection from {} to {}",
						 stream.peer_addr().unwrap(),
						 stream.local_addr().unwrap());
						 
		        let mut adr = stream.peer_addr().unwrap().port();
		        let mut playerlist2 = Arc::clone(&playerlist1);
				spawn(move|| {
					let mut stream = BufStream::new(stream);
					handle_connection(&mut stream, playerlist2, adr.into());
				});
			}
		}
	}
}

