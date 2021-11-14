use std::collections::HashMap;

#[derive(Debug, Clone, Copy)]
pub enum Type {
    Elves,
    Orcs,
    Human,
    Tree,
    Building,
}
impl Type {
    pub fn new() -> Self {
        Type::Building
    }

    pub fn create(thetype: &str) -> Self {
        match thetype {
            "elves" => {Type::Elves}
            "orcs" => {Type::Orcs}
            "human" => {Type::Human}
            "tree" => {Type::Tree}
            _ => {Type::Building}
        }
    }

    pub fn provide(&self) -> String {
        match *self {
            Type::Elves => {format!("type:elves")}
            Type::Orcs  => {format!("type:orcs")}
            Type::Human => {format!("type:human")}
            Type::Tree  => {format!("type:tree")}
            Type::Building => {format!("type:building")}
        }
    }
}

#[derive(Debug, Clone, Copy)]
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
        format!("position:({},{})", self.x, self.y)
    }
}

#[derive(Debug, Clone)]
pub struct Entity {
    pub thetype: Type,
    pub id: i64,
    pub position: Position,
    pub owner: String,
}
impl Entity {
    pub fn create(thetype1: &str, id: i64, x: i64, y: i64, owner: String) -> Self {
        Self {
            thetype: Type::create(thetype1),
            id  : id,
            position: Position::displacement(x, y),
            owner: owner,
        }
    }

    pub fn provide(&self) -> String {
        //let mut info = String::new();
        let type_info = self.thetype.provide();
        let id_info = format!("id:{}", self.id);
        let position_info = self.position.provide();
        let owner_info = format!("owner:{}", self.owner);
        format!("{} {} {} {}", type_info, id_info, position_info, owner_info)
    } 
}

#[derive(Debug, Clone)]
pub struct EntityRegistry {
    //pub positionlist: HashMap<i64, Position>,
    //pub ownerlist: HashMap<i64, String>,
    pub entitylist: HashMap<i64, Entity>,
}

impl EntityRegistry {
    pub fn create() -> Self {
        Self {
            //positionlist: HashMap::new(),
            //ownerlist: HashMap::new(),
            entitylist: HashMap::new(),
        }
    }
    pub fn show(&mut self) {
        println!("current players' informations:");
        for (id, entity) in &self.entitylist {
            println!("The infomation of player{}:{}", id, entity.provide());
        }
    }
    
    pub fn update(&mut self, thetype: &str, id: i64, pos: Position, owner: String) {
        //let mut owner_clone = owner.clone();
        let entity = Entity::create(thetype, id, pos.x, pos.y, owner);
        let entity_clone = entity.clone();
        *self.entitylist.entry(id).or_insert(entity_clone) = entity;
        //*self.positionlist.entry(id).or_insert(pos) = pos;
        //*self.ownerlist.entry(id).or_insert(owner) = owner_clone;
        self.show();
    }

    pub fn provide(&mut self) -> String {
        let mut output = String::new();
        for (id, entity) in &self.entitylist {
            let x = format!("The infomation of player{}:{}", id, entity.provide());
            output.push_str(x.as_str());
            output.push_str("\n");
        }
        return output;
    }
}



