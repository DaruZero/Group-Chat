# Group Chat - Go

This is the Golang implementation of the Group Chat project. The goal of this
implementation is to create a simple CLI chat application, adhering to the
defined features and functionality.

## Usage

### Server

The server is a simple WebSocket server that accepts connections from clients
and broadcasts messages to all connected clients. The server is started by
running the following` command:

```bash
go run server/main.go
```

The server can be configured using the following environment variables:

| Variable | Description           | Default |
| -------- | --------------------- | ------- |
| `PORT`   | The port to listen on | `8080`  |

### Client

The client is a simple CLI application that connects to the server and allows
the user to send and receive messages. The client is started by running the
following command:

```bash
go run client/main.go
```

The client can be configured using the following environment variables:

| Variable | Description              | Default               |
| -------- | ------------------------ | --------------------- |
| `SERVER` | The server to connect to | `ws://localhost:8080` |

### Docker

The server and client can be run using Docker. The following commands can be
used to build the server and client:

```bash
docker build -t group-chat-go-server -f build/server.Dockerfile .
docker build -t group-chat-go-client -f build/client.Dockerfile .
```

Before running the server and client, a Docker network must be created:

```bash
docker network create group-chat-go
```

The server and client can then be run using the following commands:

```bash
docker run -d --name group-chat-go-server --network group-chat-go -p 8080:8080 group-chat-go-server
docker run --name group-chat-go-client --network group-chat-go -it group-chat-go-client
```

Multiple clients can be run at the same time.

## License

This project is licensed under the MIT License. See the [LICENSE](../LICENSE) file
for details.
