use std::fs;
use std::io::prelude::*;
use std::net::TcpListener;
use std::net::TcpStream;
use std::str::from_utf8;

fn main(){
    let species: &'static str = "human";
    let position: [i32;2] = [100,200]; 
    let bundle: (&'static str, [i32;2]) = (species, position); 
    let listener = TcpListener::bind("127.0.0.1:8000").unwrap();
    
    for stream in listener.incoming(){
        let stream = stream.unwrap();
        handle_connection(stream, bundle);
    }
}

fn handle_connection(mut stream: TcpStream, bundle: (&'static str, [i32;2])){

    let correct_msg = b"myinformation";
    
    let (species, position) = bundle;
    
    let mut buffer = [0;13];
    
    match stream.read(&mut buffer) {
        Ok(_) => {
            let text_read = std::str::from_utf8(&buffer).unwrap();
            println!{"text i read is : {}",text_read};
            if &buffer == correct_msg {
                println!("The key is correct!\n");
                let text = format!{"client basic information:\n client's species:{}\n client's position:{:?}", species, position};
                stream.write(text.as_bytes()).unwrap();
                stream.flush().unwrap();
            } else {
                let response = format!{"The keyword is wrong, please input the right one!"};
                println!{"{}",response};
                stream.write(response.as_bytes()).unwrap();
                stream.flush().unwrap();
            }
        }
        Err(e) => {
            println!("Error happens while reading the stream: {}",e);
        }
    }
    /*
    let text = from_utf8(&buffer).unwrap();

    println!("what we read is : {}",text);
    
    let response = format!("are you ok");
   
    stream.write(response.as_bytes()).unwrap();
    stream.flush().unwrap();
   
   let get = b"GET / HTTP/1.1\r\n";
    
    let (status_line, filename) = if buffer.starts_with(get) {
        ("HTTP/1.1 200 OK", "hello.html")
    } else {
        ("HTTP/1.1 404 NOT FOUND","404.html")
    };
    
    
    if buffer.starts_with(get){
        println!("start with correct headers!");
        let contents = fs::read_to_string(filename).unwrap();
        let response = format!("{}\r\nContent-Length: {}\r\n\r\n{}",
                                status_line,               
                                contents.len(),
                                contents);
        stream.write(response.as_bytes()).unwrap();
        stream.flush().unwrap();
    } else {
        println!("Error again123!");
    }
    */
}
