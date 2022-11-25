# AuctionSystem-DISYS

## How to run

### Servers

The system will have 3 servers. You therefore need to run:

```shl
go run server/main.go 0
go run server/main.go 1
go run server/main.go 2
```

in seperate terminals.

### Clients

To run a client:

```shl
go run client/main.go [username]
```

With `[username]` a chosen username.
