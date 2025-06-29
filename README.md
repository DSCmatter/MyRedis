# MyRedis

Your own miniature redis written in Go.

---

## About

**MyRedis** is a minimal yet functional in-memory key-value store inspired by [Redis](https://redis.io/). Designed for learning and experimentation, MyRedis demonstrates the core concepts of Redis, including its protocol, data structures, and command parsing, all implemented from scratch.

---

## Features

- A Redis clone that lets you store and retrieve strings and hashes and delete them.
- Parse RESP (REdis Serialization Protocol) to handle commands and send responses.
- Handle multiple client connections simultaneously using goroutines.
- Persist data to disk using an Append Only File (AOF) so the server can recover after crashes or restarts.
- Accept and manage client connections with simple networking.
- Easily extendable design for adding new commands or features.

---

## Getting Started

### Prerequisites

- Go (Latest) 
- redis-cli (`sudo apt install redis-tools`)

### Installation

Clone this repository:

```bash
git clone https://github.com/DSCmatter/MyRedis.git
cd MyRedis
cd src/
```

### Running MyRedis

```bash
sudo snap stop redis 
go run *.go // runs all files in the directory 
redis-cli ping // will output with PONG
```

By default, the server runs on `6379`.

---

## Usage

You can interact with MyRedis using the `redis-cli` tool or any compatible Redis client:

```bash
redis-cli 
```

Try basic commands:

```bash
set name leon
get name
del name
```

```bash
hset names v1 ada
hset names v2 leon 
hgetall names

Response:
1) "v1"
2) "Ada"
3) "v2"
4) "leon
```
---

## Project Structure - /src

```bash
main.go           # Initializer
resp.go           # RESP parser
handler.go        # handles basic commands 
aof.go            # data persistence 
```

---

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Thanks for checking out MyRedis! ⭐ Star the repo if you find it useful or inspiring.