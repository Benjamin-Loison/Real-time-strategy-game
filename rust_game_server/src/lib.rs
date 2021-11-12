use std::collections::HashMap;

#[derive(Debug, Clone, Copy)]
pub struct Query {
    ID: u32,
    ttype: &'static str,
    option: &'static str,
    
}

#[derive(Debug, Clone, Copy)]
pub struct Answer {
    ID: u32,
    feedback: &'static str,
}

#[derive(Debug, Clone, Copy)]
pub struct Status {
    ID: u32,
    feedback: bool,
}

#[derive(Debug, Clone, Copy)]
pub enum Message {
    Query{ID:u32, ttype:&'static str, option:&'static str},
    Answer{ID:u32, feedback:&'static str},
    Status{ID:u32, feedback1:bool},
    None,
}
impl Message {
    pub fn call(&self) {
        match *self{
        Message::Query{ID, ttype, option} => println!("query:\nid:{}\ntype:{}\noption:{}",ID,ttype,option),
        Message::Answer{ID, feedback} => println!("answer:\nid:{}\nfeedback:{}",ID,feedback),
        Message::Status{ID, feedback1} => println!("status:\nid:{}\nfeedback:{}",ID,feedback1),
        _ => println!("none"),
        }
    }
}

#[derive(Debug, Clone, Copy)]
pub struct Position {
    x: u32,
    y: u32,
}

impl Position {
    pub fn new() -> Self {
        Self { x: 0, y: 0 }
    }
    pub fn displacement(x: u32, y: u32) -> Self {
        Self { x, y }
    }
    pub fn show(&self) {
        println!("position({},{})",self.x, self.y);
    }
}


#[derive(Debug, Clone)]
pub struct PositionRegistry {
    components: HashMap<u32, Position>,
}

impl PositionRegistry {
    pub fn create() -> Self {
        Self {
            components: HashMap::new(),
        }
    }
    pub fn show(&mut self) {
        println!("current players' postions:");
        for (address, pos) in &self.components {
            println!("player{} at ({},{})",address, pos.x, pos.y);
        }
    }
    
    pub fn update(&mut self, address:u32, pos: Position) {
        *self.components.entry(address).or_insert(pos) = pos;
        self.show();
    }
    pub fn provide(&mut self) -> String {
        let mut output = String::new();
        for (address, pos) in &self.components {
            let x = format!("player{} at ({},{})",address, pos.x, pos.y);
            output.push_str(&x);
            output.push_str("\n");
        }
        return output;
    }
}


