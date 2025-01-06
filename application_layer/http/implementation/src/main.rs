// todo
// run tcp server which is handling http requests

use std::{
    fs,
    io::{self, BufRead, BufReader, Write},
    net::{TcpListener, TcpStream},
};

use implementation::ThreadPool;

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
const PROTOCOL_VERSION: &str = "HTTP/1.1";

fn main() {
    let listener = TcpListener::bind("127.0.0.1:5000").unwrap();
    let thread_pool = ThreadPool::new(4);

    for stream in listener.incoming() {
        let mut conn = stream.unwrap();
        let n = conn.write(&[1]).unwrap();

        println!("connection established; {}", n);

        thread_pool.exec(|| match handle_connection(conn) {
            Ok(n) => println!("send response; length={}", n),
            Err(e) => println!("got error during sending response; error={:?}", e),
        });
    }
}

fn handle_connection(mut conn: TcpStream) -> io::Result<usize> {
    let buf_reader = BufReader::new(&conn);

    let http_request: Vec<String> = buf_reader
        .lines()
        .map(|res| res.unwrap())
        .take_while(|line| !line.is_empty())
        .collect();

    println!("request: {http_request:#?}");

    let request_line = &http_request[0];
    let elements: Vec<&str> = request_line.split_whitespace().collect();

    if elements.len() != 3 {
        return conn.write(&bad_request());
    }

    let method = elements[0];
    let path = elements[1];
    let protocol_version = elements[2];

    match validate_method(method) {
        false => return conn.write(&bad_request()),
        _ => (),
    };

    match validate_protocol_version(protocol_version) {
        false => return conn.write(&bad_request()),
        _ => (),
    };

    match validate_path(path) {
        Some(content) => conn.write(&status_ok(content)),
        None => conn.write(&not_found(path)),
    }
}

// validating for correct formating
fn validate_method(method: &str) -> bool {
    match method {
        "GET" => true,
        _ => false,
    }
}

fn validate_protocol_version(version: &str) -> bool {
    match version {
        PROTOCOL_VERSION => true,
        _ => false,
    }
}

fn validate_path(path: &str) -> Option<String> {
    match path {
        "/" => return Some(String::from("ROOT")),
        _ => (),
    };

    let path = path.trim_matches('/');

    match fs::read_to_string(path) {
        Ok(content) => Some(content),
        Err(e) => {
            println!("got error: {:?}", e);
            None
        }
    }
}

fn status_ok(content: String) -> Vec<u8> {
    let status_line = "200 STATUS_OK";
    let length = content.len();

    let response = format!(
        "{PROTOCOL_VERSION} {status_line}{CRLF}Content-Length: {length}{CRLF}{CRLF}{content}"
    );

    response.into_bytes()
}

fn bad_request() -> Vec<u8> {
    let status_line = "400 BAD_REQUEST";
    let response = format!("{PROTOCOL_VERSION} {status_line}{CRLF}{CRLF}");

    response.into_bytes()
}

fn not_found(path: &str) -> Vec<u8> {
    let status_line = "404 NOT_FOUND";
    let response = format!("{PROTOCOL_VERSION} {status_line}{CRLF}{CRLF}{path}");

    response.into_bytes()
}
