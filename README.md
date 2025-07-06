# GOREDIS  

A simplified Redis-like server implemented in Go, focusing on core data structures and persistence mechanisms. This project serves as a learning exercise to understand the internals of a key-value store and the Redis Serialization Protocol (RESP).  


## Table of Contents
- Features
- Project Structure
- Getting Started
- Prerequisites
- Setup
- Running the Server
- Testing with redis-cli
- Future Enhancements


## Features
This GoRedis implementation currently supports:  

 ### Core Data Structures:
**Strings:** SET, GET, DEL, EXISTS, TYPE  

**Lists:** LPUSH, RPUSH, LPOP, RPOP, LLEN  

**Hashes:** HSET, HGET, HDEL, HLEN  

**Sets:** SADD, SREM, SISMEMBER, SCARD  

**Sorted Sets:** ZADD, ZSCORE, ZREM, ZCARD  

**Basic Commands:** PING, ECHO  

**Persistence:**  

**Append-Only File (AOF):** Commands that modify the database are appended to appendonly.aof.  

**AOF Loading:** Upon server startup, the appendonly.aof file is replayed to restore the database state.  

**Networking:** Simple TCP server listening on 0.0.0.0:6379 (IPv4).  


<pre> ```## Project Structure

goredis/
├── main.go               # Main entry point, server initialization, AOF setup, graceful shutdown.
├── server/
│   ├── server.go         # Handles TCP connections, client management, AOF integration.
│   └── client.go         # Represents a connected client, handles RESP I/O and command dispatch.
├── database/
│   ├── database.go       # In-memory data store, handles key-value storage and TTL.
│   ├── value.go          # Interface for different Redis data types.
│   ├── string.go         # Implementation of Redis String type.
│   ├── list.go           # Implementation of Redis List type.
│   ├── hash.go           # Implementation of Redis Hash type.
│   ├── set.go            # Implementation of Redis Set type.
│   └── zset.go           # Implementation of Redis Sorted Set type (simplified).
├── resp/
│   └── resp.go           # Handles encoding and decoding of Redis Serialization Protocol (RESP).
├── command/
│   ├── processor.go      # Dispatches commands to handlers, integrates with AOF.
│   └── handlers.go       # Contains implementations for various Redis commands.
├── persistence/
│   └── aof.go            # Manages Append-Only File operations (write and load).
├── pubsub/
│   └── pubsub.go         # Placeholder for Publish/Subscribe functionality.
├── transaction/
│   └── transaction.go    # Placeholder for Redis Transactions (MULTI/EXEC/DISCARD).
├── utils/
│   └── utils.go          # Utility functions (e.g., panic recovery).
└── go.mod                # Go module definition and dependencies.
```</pre>



**Getting Started**
**Prerequisites**

**Go Language:** Ensure Go is installed on your system (version 1.22 or higher recommended). Download from golang.org/dl.
**redis-cli (Optional, but Recommended):** The official Redis command-line interface is invaluable for testing.
Windows: Download the Redis Stack Windows Installer and ensure "Add Redis to the PATH environment variable" is checked during installation.
macOS/Linux: Install via your system's package manager (e.g., brew install redis on macOS, sudo apt-get install redis-tools on Ubuntu).

**Setup**
Clone the Repository (or create the structure):
If you're starting from scratch, ensure your project structure matches the one above. If you're using a Git repository, clone it:
git clone https://github.com/your-username/goredis.git
cd goredis

(Replace https://github.com/your-username/goredis.git with your actual repository URL if applicable, or just navigate to your existing project directory).
Initialize Go Module Dependencies:
Run go mod tidy in the project root to download and manage dependencies.
go mod tidy


**Running the Server**
<pre>```
Start the GoRedis server:
Open your terminal or command prompt in the goredis project root and run:
go run main.go

You should see output similar to this:
2025/07/05 19:08:34 AOF: Loading data from appendonly.aof...
2025/07/05 19:08:34 AOF: Finished loading appendonly.aof.
2025/07/05 19:08:34 Server listening on 0.0.0.0:6379
2025/07/05 19:08:34 GoRedis server started on address 0.0.0.0:6379

The server will now be listening for connections on port 6379. Keep this terminal window open.
Testing with redis-cli
Open a separate, new terminal or command prompt window to connect to your running GoRedis server.
Connect to the server:
redis-cli

You should see the 127.0.0.1:6379> prompt.
Execute sample commands:
Strings & Persistence:
127.0.0.1:6379> SET mykey "Hello GoRedis!"
OK
127.0.0.1:6379> GET mykey
"Hello GoRedis!"
127.0.0.1:6379> SET anotherkey "This will persist" EX 10
OK
127.0.0.1:6379> EXISTS anotherkey
(integer) 1

Now, go back to the server terminal (Ctrl+C to stop it). Then restart it (go run main.go). Reconnect with redis-cli and check:
127.0.0.1:6379> GET mykey
"Hello GoRedis!"
127.0.0.1:6379> GET anotherkey # This should still be there if less than 10 seconds passed since first SET
"This will persist"

Lists:
127.0.0.1:6379> LPUSH mylist apple banana cherry
(integer) 3
127.0.0.1:6379> LPOP mylist
"cherry"
127.0.0.1:6379> RPUSH mylist date elderberry
(integer) 4
127.0.0.1:6379> LLEN mylist
(integer) 3

Hashes:
127.0.0.1:6379> HSET myuser name "Alice" email "alice@example.com"
(integer) 2
127.0.0.1:6379> HGET myuser name
"Alice"
127.0.0.1:6379> HDEL myuser email
(integer) 1
127.0.0.1:6379> HLEN myuser
(integer) 1

Sets:
127.0.0.1:6379> SADD mysports football basketball tennis
(integer) 3
127.0.0.1:6379> SISMEMBER mysports basketball
(integer) 1
127.0.0.1:6379> SCARD mysports
(integer) 3
127.0.0.1:6379> SREM mysports tennis
(integer) 1

Sorted Sets:
127.0.0.1:6379> ZADD myleaderboard 100 "PlayerA" 50 "PlayerB" 120 "PlayerC"
(integer) 3
127.0.0.1:6379> ZSCORE myleaderboard "PlayerA"
"100"
127.0.0.1:6379> ZCARD myleaderboard
(integer) 3
127.0.0.1:6379> ZREM myleaderboard "PlayerB"
(integer) 1


Exit redis-cli:
127.0.0.1:6379> QUIT


Stop the GoRedis server: Go back to the terminal running go run main.go and press Ctrl+C.```<\pre>



## Future Enhancements

This project is a solid foundation. Here are some key areas for future development to make it more like a full-fledged Redis:  

- AOF Rewrite (BGREWRITEAOF): Implement the logic to optimize the AOF file by rewriting it in the background, removing redundant commands.

- RDB Snapshotting: Add support for saving and loading binary RDB snapshots of the database.

- Comprehensive Command Set: Implement more commands for each data type (e.g., LRANGE, HGETALL, SMEMBERS, ZRANGE).

- Transactions (MULTI/EXEC/DISCARD/WATCH): Fully implement Redis transactions with optimistic locking.

- Publish/Subscribe (Pub/Sub): Add PUBLISH, SUBSCRIBE, PSUBSCRIBE for real-time messaging.

- Advanced TTL Management: Implement background key eviction for expired keys more efficiently.

- Authentication (AUTH): Add a simple password-based authentication mechanism.

- Error Handling & Robustness: Enhance error handling, especially for network issues and malformed commands.

- Performance Optimizations: Explore more performant data structures (e.g., skip lists for sorted sets, specialized list implementations).

- Metrics & Monitoring (INFO): Provide server statistics and information.

