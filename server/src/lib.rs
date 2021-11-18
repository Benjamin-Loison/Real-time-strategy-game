use std::collections::HashMap;
use std::cmp;

#[derive(Debug, Clone, Copy)]
pub enum Arch {
    Floor,
    ElfBuilding,
    OrcBuilding,
    HumanBuilding,
}
impl Arch {
    pub fn create(arch: &str) -> Self {
        match arch {
            "Floor" => {Arch::Floor}
            "ElfBuilding" => {Arch::ElfBuilding}
            "OrcBuilding" => {Arch::OrcBuilding}
            _ => {Arch::HumanBuilding}
        }
    }
    pub fn provide(&self) -> String {
        match *self {
            Arch::HumanBuilding => {format!("HumanBuilding")}
            Arch::ElfBuilding => {format!("ElfBuilding")}
            Arch::OrcBuilding => {format!("Orcbuilding")}
            _ => {format!("Floor")}
        }
    }
}

#[derive(Debug, Clone, Copy)]
pub enum Type {
    Elves,
    Orcs,
    Human,
}
impl Type {
    pub fn create(thetype: &str) -> Self {
        match thetype {
            "elves" => {Type::Elves}
            "orcs" => {Type::Orcs}
            _ => {Type::Human}
        }
    }

    pub fn provide(&self) -> String {
        match *self {
            Type::Elves => {format!("type:elves")}
            Type::Orcs  => {format!("type:orcs")}
            _ => {format!("type:human")}
        }
    }
}

#[derive(Debug, Clone, Copy, Eq, Hash, PartialEq)]
pub struct Position {
    pub x: i64,
    pub y: i64,
}

impl Position {
    pub fn new() -> Self {
        Self { x: 0, y: 0 }
    }
    pub fn displacement(x: i64, y: i64) -> Self {
        Self { x, y }
    }
    pub fn show(&self) {
        println!("position({},{})", self.x, self.y);
    }
    pub fn provide(&self) -> String {
        //format!("position:({},{})", self.x, self.y)
        format!("{},{}", self.x, self.y)
    }
}

#[derive(Debug, Clone)]
pub struct Entity {
    pub thetype: Type,
    pub id: i64,
    pub position: Position,
    pub owner: String,
    pub view: Position,
}
impl Entity {
    pub fn create(thetype1: &str, id: i64, x: i64, y: i64, owner: String, width: i64, height: i64) -> Self {
        Self {
            thetype: Type::create(thetype1),
            id: id,
            position: Position::displacement(x, y),
            owner: owner,
            view: Position::displacement(width, height),
        }
    }

    pub fn provide(&self) -> String {
        let type_info = self.thetype.provide();
        let id_info = format!("id:{}", self.id);
        let position_info = self.position.provide();
        let owner_info = format!("owner:{}", self.owner);
        let view_info = format!("view: ({},{})", self.view.x, self.view.y);
        format!("{} {} {} {} {}", type_info, id_info, position_info, owner_info, view_info)
    } 
}

#[derive(Debug, Clone)]
pub struct Thing {
    thetype: Arch,
    id: i64,
    thestring: String,
}
impl Thing {
    pub fn create(thetype1: &str, id: i64, thestring: &str) -> Self {
        Self {
            thetype: Arch::create(thetype1),
            id: id,
            thestring: thestring.to_string(),
        }
    }
    pub fn provide(&self) -> String {
        let type_info = self.thetype.provide();
        let id_info = format!("{}", self.id);
        format!{"{},{},{}", type_info, id_info, self.thestring}
    }
}

#[derive(Debug, Clone)]
pub struct EntityRegistry {
    pub entitylist: HashMap<i64, Entity>,
    pub world: HashMap<Position, Thing>,
    pub size: i64,
    pub ownershiplist: HashMap<String, Vec<i64>>,
}

impl EntityRegistry {
    pub fn create(size: i64) -> Self {
        //let size: i64 = 10;
        let mut new_world = HashMap::new();
        for i in 0..size {
            for j in 0..size {
                let pos_ij = Position::displacement(i, j);
                new_world.insert(pos_ij, Thing::create("Floor", 0, "string"));
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
    
    pub fn update(&mut self, thetype: &str, id: i64, pos: Position, owner: String, width: i64, height: i64) {
        let owner_clone1 = owner.clone();
        let owner_clone2 = owner.clone();
        let entity = Entity::create(thetype, id, pos.x, pos.y, owner, width, height);
        let mut new_ownership = Vec::new();
        match self.ownershiplist.get(&owner_clone1) {
            Some(old_ownership) => {
                //println!("printold:{:?}", old_ownership);
                let mut old_ownership1 = old_ownership.clone();
                //println!("printid:{}", id);
                if old_ownership1.contains(&id) == false {
                    old_ownership1.push(id);
                }
                //old_ownership1.push(id);
                //println!("printold1:{:?}", old_ownership1);
                new_ownership = old_ownership1.to_vec();
                //println!("printnew:{:?}", new_ownership);
            }
            None => {
                let mut temp_ownership = Vec::new();
                temp_ownership.push(id);
                new_ownership = temp_ownership;
            }
        }
        let entity_clone = entity.clone();
        let new_ownership_clone = new_ownership.clone();
        *self.entitylist.entry(id).or_insert(entity_clone) = entity;
        *self.ownershiplist.entry(owner_clone2).or_insert(new_ownership_clone) = new_ownership;
        println!("current game view after update:");
        self.show();
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

    pub fn bigframe(&self, id: i64) -> String {
        let id_entity = self.entitylist.get(&id).unwrap();
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
        println!("the bigframe of player {}: \n{}", id, &output);
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
    /*
    pub fn players_nearby(&self, id: i64) -> String {
        let mut internal_playerlist = EntityRegistry::create(10);
		let fixed_entity = self.entitylist.get(&id).unwrap();
		let playerlist_unwrap_clone = self.clone();
		for (new_id, new_entity) in playerlist_unwrap_clone.entitylist{
			let fixed_entity_clone = fixed_entity.clone();
			let mut owner_clone = new_entity.owner.clone();
			let mut view_clone = new_entity.view.clone();
			let mut position_clone = new_entity.position.clone();
			if self.inside_or_not(fixed_entity_clone, new_entity) == true {
				internal_playerlist.update("human", new_id, position_clone, owner_clone, view_clone.x, view_clone.y);
			}
		}
        let mut output = String::from(format!("All the players near player {}:\n", id));
		for (ids, entity) in internal_playerlist.entitylist {
			let x = format!("player{}: {}", ids, entity.provide());
			output.push_str(&x);
			output.push_str("\n");
		}
        output
    }
    */
    pub fn players_nearby(&self, id: i64) -> HashMap<i64, Entity>{
        let mut internal_playerlist = EntityRegistry::create(10);
		let fixed_entity = self.entitylist.get(&id).unwrap();
		let playerlist_unwrap_clone = self.clone();
		for (new_id, new_entity) in playerlist_unwrap_clone.entitylist{
			let fixed_entity_clone = fixed_entity.clone();
			let mut owner_clone = new_entity.owner.clone();
			let mut view_clone = new_entity.view.clone();
			let mut position_clone = new_entity.position.clone();
			if self.inside_or_not(fixed_entity_clone, new_entity) & (new_id != id) {
				internal_playerlist.update("human", new_id, position_clone, owner_clone, view_clone.x, view_clone.y);
			}
		}
        return internal_playerlist.entitylist;
        /*
        let mut output = String::from(format!("All the players near player {}:\n", id));
		for (ids, entity) in internal_playerlist.entitylist {
			let x = format!("player{}: {}", ids, entity.provide());
			output.push_str(&x);
			output.push_str("\n");
		}
        output
        */
    }
}



