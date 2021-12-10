extern crate bufstream;
extern crate crossbeam_channel;
use std::time::SystemTime;
use std::io::{Write, BufWriter};
use std::net::{TcpStream};
pub use super::lib::{EntityRegistry, Position, Type, Entity, ServerMessage};

pub fn init_msg_sending(stream: TcpStream) {
    let init_msg1 = ServerMessage{MessageType: , GameMap: ,Players: , Id: };
}