# thorium-go
Thorium.NET in Golang

A game server cluster management tool using HTTP/REST.

Can run as a one-node or multi-node cluster. Some knowledge of network management and configuration is required to setup and run. I.e. setting up machines, exposing addresses and ports, etc. 

Before Starting
- install Golang (1.4) and setup $GOPATH to work with this repository
- configure and run a Postgres instance
- configure and run a Redis instance
- modify the address and port info inside database/thordb.go script to point to your datastore addresses (from above)
- generate your own RSA keys and replace the ones in keys/ as outlined here: https://gist.github.com/cryptix/45c33ecf0ae54828e63b
- compile your own bolt-server.exe and implement the same startup requests as found in the reference implementation cmd/bolt-server/bolt-server.go
- please note: HTTP/REST can be implemented in Unity and not just Golang!

Getting Started - New Cluster

1. Start Postgres and Redis

2. Start a Master node
> go run cmd/master-server/master-server.go

3. Start a Game Server node
> go run cmd/game-server/game-server.go

Thats it!

Getting Started - New Client

- for tips on implementing a new client, see the reference implementation at cmd/client/client.go and the cmd/test/ scripts for demonstrations of different use-cases
