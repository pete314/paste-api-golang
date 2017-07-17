Build [![CircleCI](https://circleci.com/gh/pete314/paste-api-golang/tree/master.svg?style=svg)](https://circleci.com/gh/pete314/paste-api-golang/tree/master)

# Introduction
This repository contains a REST API written in golang, with docker development and deployment support. The purpose of this is nothing more than to give an example of where to start with docker and go.

# Dependencies
- [Docker](https://www.docker.com/)
- [Redis](https://redis.io) - can be installed on host or as/in a container
- [go](https://golang.org/) - can be installed on host or build in container

# Running in docker
The configuration ```config/server-config.local.json``` should be updated accordingly with the service hosts/ports.
To run this project in a docker container:
```bash
cd $PROJECT_ROOT
## Build the container
sudo docker build  -t paste-api-dev .

## Run source
sudo docker run -p 6379 -p 8080:8080 -it paste-api-dev
```

The above commands will build and run the container, which is ok for development. For production there is actually less needed, so the script deploy script uses  the [scratch](https://hub.docker.com/_/scratch/) container.
To build a container for production, the host needs golang installed and the environment setup. To build run

```bash
cd $PROJECT_ROOT/scripts

## Build the container
./deploy-api.sh

## Run the container
sudo docker run -p 6379 -p 8080:8080 -it paste-api-scratch
```

# Run tests
In the docker container or on the host the tests can be run directly 
```bash
cd $PROJECT_ROOT/api

## Run tests
go test -v ./src/tests/modules/v0/unit/...
```

It is also part of the build script (Makefile) and can be executed with ```make test```

# HTTP Benchmarks 

Benchmark [wrk](https://github.com/wg/wrk)

|request | http method | throughput (/sec) | latency|
|---|---|---|---|
|/v0/paste/{ID}| GET | 25000 | 1.21 ms|
|/v0/paste| PUT | 18000| 1.6 ms|


```bash
## Run the API
sudo docker run -p 6379 -p 8080:8080 -it paste-api-scratch

## Running GET benchmarks
# Insert an entry
ID=$(curl -s -X PUT -d '{"note": "Todo: benchmark go-paste"}' 127.0.0.1:8080/v0/paste | awk -F":" '{print $4}' | awk -F'"' '{print $2}')

# Run GET requests
wrk -t5 -c 25 -d5s http://127.0.0.1:8080/v0/paste/$ID

Running 5s test @ http://127.0.0.1:8080/v0/paste/c1ff2459-b58c-4c70-bf77-ea6d52f15936
  5 threads and 25 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.21ms    1.41ms  23.26ms   93.09%
    Req/Sec     5.13k     0.99k   15.04k    82.07%
  128219 requests in 5.10s, 46.83MB read
Requests/sec:  25155.74
Transfer/sec:      9.19MB

# Run PUT requests
echo 'wrk.method = "PUT" wrk.body   = "{\"note\":\"test\"}" wrk.headers["Content-Type"] = "application/json"' > /tmp/paste_put.lua

wrk -t5 -c 25 -d5s -s /tmp/paste_put.lua http://127.0.0.1:8080/v0/paste

Running 5s test @ http://127.0.0.1:8080/v0/paste
  5 threads and 25 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.64ms    1.80ms  36.26ms   92.01%
    Req/Sec     3.77k   576.99     5.89k    70.40%
  93994 requests in 5.03s, 48.12MB read
Requests/sec:  18698.24
Transfer/sec:      9.57MB

```
