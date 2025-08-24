# Go-Redis:
##  A minimal implementation of the redis server and CLI

This project is based on the coding challenges by John Crickett

[Build your own redis](https://codingchallenges.fyi/challenges/challenge-redis/)

## Purpose
My purpose in building this project is to learn more about Go and to challenge myself to build a tool from scratch. It's been a lot of fun to try and figure this out and taught me more about golang. 

## Running locally 

To run locally, pull down the code into your own repository. 
Start the server with `go run .` or build into your own binary and start with `./<your-binary>`

This includes an interactive shell mode similar to redis-cli so you could run `go run .` in one terminal session and `go run . start` in another. 
This will interact similar to the redis-cli. You can also connect to this server with the redis-cli command. 

## Benchmarks
The redis-benchmark utility is a way to test how efficient your redis implementation is with concurrent connections. Here's how mine did:

```
$ redis-benchmark -t set,get, -n 100000 -q
SET: 101522.84 requests per second, p50=0.327 msec
GET: 122699.39 requests per second, p50=0.215 msec
```

