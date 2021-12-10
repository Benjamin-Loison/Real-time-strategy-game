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
    pub id: String,
    pub position: Position,
    pub owner: String,
    pub view: Position,
}
impl Entity {
    pub fn create(thetype1: &str, id: &str, x: i64, y: i64, owner: String, width: i64, height: i64) -> Self {
        Self {
            thetype: Type::create(thetype1),
            id: id.to_string(),
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
    id: String,
    thestring: String,
}
impl Thing {
    pub fn create(thetype1: &str, id: &str, thestring: &str) -> Self {
        Self {
            thetype: Arch::create(thetype1),
            id: id.to_string(),
            thestring: thestring.to_string(),
        }
    }
    pub fn provide(&self) -> String {
        let type_info = self.thetype.provide();
        let id_info = format!("{}", self.id);
        format!{"{},{},{}", type_info, id_info, self.thestring}
    }
}
