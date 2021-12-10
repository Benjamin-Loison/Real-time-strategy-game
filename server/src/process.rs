extern crate bufstream;
extern crate regex;
use crossbeam_channel::{unbounded, Sender, Receiver};
use std::net::{TcpStream};
use std::sync::Arc;
use std::sync::Mutex;
pub use super::lib::{EntityRegistry, Position, Type, Entity};
pub use super::handle::{logging, query_client, answer_client,status_client};




// update all the data of playerlist
pub fn position_update(stream: TcpStream, id: &str, client_x: &str, client_y: &str, width: &str, height: &str, playerlist: Arc<Mutex<EntityRegistry>>, sender: Sender<Arc<Mutex<EntityRegistry>>>) {
	let client_latest_position = Position::displacement(client_x.parse::<i64>().unwrap(), client_y.parse::<i64>().unwrap());
	let playerlist_to_send = Arc::clone(&playerlist);
	let mut latest_playerlist = playerlist.lock().unwrap();
	let owner: String = stream.peer_addr().unwrap().to_string();
	logging("position_update", "");
	latest_playerlist.update("show", "human", id, client_latest_position, owner, width.parse::<i64>().unwrap(), height.parse::<i64>().unwrap());
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
	let mut playerlist_unwrap = playerlist.lock().unwrap();
	//let mut ownershiplist = playerlist_unwrap.ownershiplist.get(&owner).unwrap();
	match playerlist_unwrap.ownershiplist.get(&owner) {
		Some(mut ownershiplist) => {
			for id in ownershiplist.iter() {
				let stream_clone1 = stream.try_clone().unwrap();
				match thestring {
					"set" => {
						let list = playerlist_unwrap.players_nearby(id);
						for (ids, entitys) in list {
							println!();
							let stream_clone1_clone = stream_clone1.try_clone().unwrap();
							let mut output = String::new();
							output.push_str(&format!("set,{},{}",ids, entitys.position.provide()));
							query_client(stream_clone1_clone, &id.to_string(), &output);
						}
					}
					_ => {
						query_client(stream_clone1, &id.to_string(), "get");
					}
				}
			}
		}
		_ => {
			logging("prepare_inside_players", 
			"There is no such player in our world, please send us your player info following the protocol!");
		}
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
	let map_data: String = playerlist_clone2_unwrap_clone.bigframe(id);
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
	logging("query_after_update", "");
	prepare_inside_players(stream_clone1, owner.to_string(), playerlist, thestring);
}

