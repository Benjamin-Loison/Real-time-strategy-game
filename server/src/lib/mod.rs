use std::collections::HashMap;
use std::cmp;
pub mod lib;
//pub mod new;
pub use self::lib::{Arch, Type, Position, Entity, Thing};
//pub use self::new;

#[derive(Debug, Clone)]
pub struct EntityRegistry {
    pub entitylist: HashMap<String, Entity>,
    pub world: HashMap<Position, Thing>,
    pub size: i64,
    pub ownershiplist: HashMap<String, Vec<String>>,
}

impl EntityRegistry {
    pub fn create(size: i64) -> Self {
        //let size: i64 = 10;
        let mut new_world = HashMap::new();
        for i in 0..size {
            for j in 0..size {
                let pos_ij = Position::displacement(i, j);
                new_world.insert(pos_ij, Thing::create("Floor", "0", "string"));
            }
        }
        Self {
            entitylist: HashMap::new(),
            world: new_world,
            size: size,
            ownershiplist: HashMap::new(),
        }
    }
    /*UNFINISHED*/
    pub fn show(&mut self) {
        let output = self.provide();
    }
    
    pub fn update(&mut self, show_or_not: &str, thetype: &str, id: &str, pos: Position, owner: String, width: i64, height: i64) {
        let owner_clone1 = owner.clone();
        let owner_clone2 = owner.clone();
        let entity = Entity::create(thetype, id, pos.x, pos.y, owner, width, height);
        let mut new_ownership = Vec::new();
        match self.ownershiplist.get(&owner_clone1) {
            Some(old_ownership) => {
                let mut old_ownership1 = old_ownership.clone();
                if old_ownership1.contains(&id.to_string()) == false {
                    old_ownership1.push(id.to_string());
                }
                new_ownership = old_ownership1.to_vec();
            }
            None => {
                let mut temp_ownership = Vec::<String>::new();
                temp_ownership.push(id.to_string());
                new_ownership = temp_ownership;
            }
        }
        let entity_clone = entity.clone();
        let new_ownership_clone = new_ownership.clone();
        *self.entitylist.entry(id.to_string()).or_insert(entity_clone) = entity;
        *self.ownershiplist.entry(owner_clone2).or_insert(new_ownership_clone) = new_ownership;
        if show_or_not == "show" {
            println!("current game view after update:");
            self.show();
        }
    }


    pub fn provide_world(&mut self) -> String {
        let mut output = String::new();
        for i in 0..self.size {
            for j in 0..self.size {
                let new_pos = Position::displacement(i,j);
                let new_thing = self.world.get(&new_pos).unwrap();
                output.push_str(",,");
                output.push_str(&new_thing.provide());
            }
            output.push_str("\n");
        }
        println!("{}", &output);
        return output;
    }

    pub fn provide_players(&mut self) -> String {
        let mut output = String::new();
        for (id, entity) in &self.entitylist {
            let x = format!("The information of players{}:{}", id, entity.provide());
            output.push_str(x.as_str());
            output.push_str("\n");
        }
        println!("{}", &output);
        return output;
    }

    pub fn provide_ownership(&mut self) -> String {
        let mut output = String::new();
        for (owners, id_list) in &self.ownershiplist {
            let x = format!("The information of owner{}:{:?}", owners, id_list);
            output.push_str(x.as_str());
            output.push_str("\n");
        }
        println!("ownershiplist:\n{}", &output);
        return output;
    }
    
    pub fn provide(&mut self) -> String {
        format!("{}\n{}\n{}\n", self.provide_world(), self.provide_players(), self.provide_ownership())
    }

    pub fn bigframe(&self, id: &str) -> String {
        let id_entity = self.entitylist.get(id).unwrap();
        let xlowerbound = cmp::max(0, id_entity.position.x - id_entity.view.x);
        let xupperbound = cmp::min(self.size, id_entity.position.x + id_entity.view.x + 1);
        let ylowerbound = cmp::max(0, id_entity.position.y - id_entity.view.y);
        let yupperbound = cmp::min(self.size, id_entity.position.y + id_entity.view.y + 1);
        let mut output = String::new();
        for i in xlowerbound..xupperbound {
            for j in ylowerbound..yupperbound {
                let new_pos = Position::displacement(i,j);
                let new_thing = self.world.get(&new_pos).unwrap();
                output.push_str(&format!(",,"));
                output.push_str(&new_thing.provide());
            }
           // output.push_str("\n");
        }
        //println!("the bigframe of player {}: \n{}", id, &output);
        return output;
    }

    // check if the new position is inide the scope of fixed position with width and height
    pub fn inside_or_not(&self, fixed_entity: Entity, new_entity: Entity) -> bool {
        let w = fixed_entity.view.x;
        let h = fixed_entity.view.y;
        let x_ok: bool = ((fixed_entity.position.x - new_entity.position.x).abs() <= w);
        let y_ok: bool = ((fixed_entity.position.y - new_entity.position.y).abs() <= h);
        return (x_ok & y_ok);
    }
    
    pub fn players_nearby(&self, id: &str) -> HashMap<String, Entity>{
        let mut internal_playerlist = EntityRegistry::create(10);
		let fixed_entity = self.entitylist.get(id).unwrap();
		let playerlist_unwrap_clone = self.clone();
		for (new_id, new_entity) in playerlist_unwrap_clone.entitylist{
			let fixed_entity_clone = fixed_entity.clone();
			let mut owner_clone = new_entity.owner.clone();
			let mut view_clone = new_entity.view.clone();
			let mut position_clone = new_entity.position.clone();
			if self.inside_or_not(fixed_entity_clone, new_entity) & (new_id != id) {
				internal_playerlist.update("notshow", "human", &new_id, position_clone, owner_clone, view_clone.x, view_clone.y);
			}
		}
        return internal_playerlist.entitylist;
    }
}



