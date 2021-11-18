extern crate bufstream;
extern crate regex;
mod lib;

use std::str::FromStr;
use std::io::Write;
use std::net::{TcpListener, TcpStream};
use std::net::SocketAddr;
use std::thread;
use bufstream::BufStream;
use std::io::BufRead;
use regex::Regex;
use std::sync::Arc;
use std::sync::Mutex;
use std::sync::mpsc;
use std::sync::mpsc::{Sender, Receiver};
use self::lib::{PositionRegistry,Position};

pub mod prepare {
    pub fn logging(from: &str, msg: &str) {
        println!("[{}]: {}", from, msg);
    }
    
    pub fn answer_client(stream: TcpStream, id: &str, options: &str) {
        let ans = format!("A{}.{}\n", id, options);
        let mut streamwriter = BufWriter::with_capacity(100, stream);
        streamwriter.write(ans.as_bytes());
        streamwriter.flush();
    }
    
    // update all the data of playerlist
    pub fn position_update(stream: TcpStream, id: &str, client_x: &str, client_y: &str, playerlist: Arc<Mutex<PositionRegistry>>) {
        let client_latest_position = Position::displacement(client_x.parse::<i64>().unwrap(), client_y.parse::<i64>().unwrap());
        let mut latest_playerlist = playerlist.lock().unwrap();
        let owner: &str = stream.peer_addr().unwrap().to_string();
        latest_playerlist.update(id.parse::<i64>().unwrap(), client_latest_position, owner.parse::<i64>().unwrap());
    }
    
    // check if the new position is inide the scope of fixed position with width and height
    pub fn inside_or_not(fixed_position: Position, new_position: Position, width: &str, height: &str) -> bool {
        let w = width.parse::<i64>().unwrap();
        let h = height.parse::<i64>().unwrap();
        let x_ok: bool = ((fixed_position.x - new_position.x).abs() <= w);
        let y_ok: bool = ((fixed_position.y - new_position.y).abs() <= h);
        return (x_ok & y_ok);
    }
    
    
    // prepare the string map_data from playerlist to send back to client
    pub fn prepare_map(stream: TcpStream, id: &str, width: &str, height: &str, playerlist: Arc<Mutex<PositionRegistry>>) -> String {
        let mut internal_playerlist = PositionRegistry::create();
        let playerlist_unwrap = playerlist.lock().unwrap();
        let fixed_position = playerlist_unwrap.positionlist.get(&id.parse::<i64>().unwrap()).unwrap();
        let owner: &str = stream.peer_addr().unwrap().to_string();
        for (new_id, new_position) in &playerlist_unwrap.positionlist{
            if inside_or_not(*fixed_position, *new_position, width, height) == true {
                internal_playerlist.update(*new_id, *new_position, owner.parse::<i64>().unwrap());
            }
        }
        return internal_playerlist.provide();
    }
    
    //send map data to client after receiving query with key word "map"
    pub fn send_map_data(stream: TcpStream, id: &str, options: &str, playerlist: Arc<Mutex<PositionRegistry>>) {
        let xywh: Vec<&str> = options.split(',').collect();
        let mut playerlist_clone1 = Arc::clone(&playerlist);
        let mut playerlist_clone2 = Arc::clone(&playerlist);
        position_update(stream, id, xywh[0], xywh[1], playerlist_clone1);
        let map_data: String = prepare_map(stream, id, xywh[2], xywh[3], playerlist_clone2);
        answer_client(stream, id, map_data.as_str());
    }
    
    //send new query to clients after detecting conflicts from latest update
    pub fn query_after_update(stream: TcpStream, id: &str, playerlist: Arc<Mutex<PositionRegistry>>) {
        let owner: &str = stream.peer_addr().unwrap();
        answer_client(stream, id, "you need to tell me your info");
    }
}
