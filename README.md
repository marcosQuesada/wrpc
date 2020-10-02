# [WIP]wRPC

gRPC on top of Websockets, naive implementation.

Project aim was not to replace HTTP2, just the transport connection, what it means that all HTTP2 headers, trailers and body are handled in the same way, as just "traffic" that is forwarded between connection peers. So that, once transport is available and ready, http2 handshake is done in a transparent way without any header manipulation.

gRPC allows http2 connection replacement, in server side a buffered connection Listener mimics the behaviour of the real listener, in the backgrounds, a piped connection is attached to the grpc Server, while transport connection gets adapted, in that scheme traffic from transport connection gets forwarded to the grpc piped one, on downstream everything flows the same way, received traffic from grpc piped conn gets forwarded to transport connection. 

On client side, client creation needs to use WithContextDialer option, that enable us to use our to adapted connection, from websocket transport to the net.Conn that gRPC uses internally  
 
## Run example
Start server as: 
```
go run main.go server
```
Run client as:
```
go run main.go client
```
### Enable gRPC log detail
```
export GRPC_GO_LOG_SEVERITY_LEVEL=info 
export GRPC_GO_LOG_VERBOSITY_LEVEL=99
``` 