# WS-RPC

Protocol Buffer Plugin for WebSocket based RPC generated code.

Similar to gRPC, WsRPC will generate server and client code for your protobuf defined RPC.

Based on templates, you can change the output depending on your needs.

Default supported templates:

[Golang WebSocket Service](./templates/_wsrpc.pb.go.template)

[Unity/C# WebSocket Client](./templates/_client.pb.cs.template)

[API Doc Generator](./templates/_doc.pb.js.template)


## Build
```
go build -o protoc-gen-ws-rpc .
```
It will generate the binary file needed to run the protobuf plugin

## Run the plugin
Make sure protoc-gen-ws-rpc is in the current path

The output file is based on the template name (ex.: _service.cs.template using leaderboard.proto will generate leaderboard_service.cs)

```
protoc --plugin protoc-gen-ws-rpc  --ws-rpc_out=./ --ws-rpc_opt=template=service.cs.template Example.proto
```
--ws-rpc_out    path for the generated file (output)

--ws-rpc_opt    template= path of your template 

## Work in Progress

I need to add the necessary components to make use of the generated code.

Like with gRPC it will be language specific libraries to be able to run the WsRPC files without additional steps.
