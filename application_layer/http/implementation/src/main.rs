// todo
// run tcp server which is handling http requests

use std::{
    io::{BufRead, BufReader, Write},
    net::{TcpListener, TcpStream},
};

/* HTTP Request Format

Method Request-URI HTTP-Version CRLF
headers CRLF
message-body

----------------------------------------------

HTTP Response Format

HTTP-Version Status-Code Reason-Phrase CRLF
headers CRLF
message-body

*/

const CRLF: &str = "\r\n";

fn main() {
    let listener = TcpListener::bind("127.0.0.1:5000").unwrap();

    for stream in listener.incoming() {
        let mut conn = stream.unwrap();
        let n = conn.write(&[1]).unwrap();

        println!("connection established; {}", n);

        handle_connection(conn);
    }
}

fn handle_connection(conn: TcpStream) {
    let buf_reader = BufReader::new(&conn);

    let http_request: Vec<String> = buf_reader
        .lines()
        .map(|res| res.unwrap())
        .take_while(|line| !line.is_empty())
        .collect();

    println!("request: {http_request:#?}");
    // handle_http_request(http_request);
    // todo handle incoming request from connection
}

fn handle_http_request(req: Vec<String>) {}
