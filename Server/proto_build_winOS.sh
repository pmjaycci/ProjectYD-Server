#!/bin/bash
echo "---------- Start Convert To Language From Protofile ----------"
echo "Protofile Convert To C#"
./protofiles/win_grpc_protoc_x64-1.59.0/protoc --csharp_out=grpc --grpc_out=grpc --plugin=protoc-gen-grpc=./protofiles/win_grpc_protoc_x64-1.59.0/grpc_csharp_plugin.exe  ./protos/global_grpc.proto
echo "C# Convert Success From Protofile"
echo "-------------------------"
echo "Protofile Convert To Golang"
./protofiles/win_grpc_protoc_x64-1.59.0/protoc --go_out=grpc --go-grpc_out=grpc ./protos/global_grpc.proto
echo "Golang Convert Success From Protofile" 
echo "---------- END ----------"