package server

import (
	"net"
	global_grpc "project_yd/grpc"
	"sync"

	"google.golang.org/grpc"
)

var GlobalGrpcEvent *sync.Map
var mutex sync.Mutex

// 함수 시그니처: 클라이언트로부터 받은 key 값에 따라 호출할 함수들
type RpcKeyHandlerFunc func(payload string) string

// -- rpc로 쓰일 function 등록
func RegistRpc(rpcKey string, function RpcKeyHandlerFunc) {
	_, ok := GlobalGrpcEvent.Load(rpcKey)
	if ok {
		println("RPC KEY : ", rpcKey, " is Duplicate!")
		return
	}
	GlobalGrpcEvent.Store(rpcKey, function)
	println("RPC KEY : ", rpcKey, " is Regist Success!")
}

// -- 클라이언트로부터 호출된 rpc명으로 호출할 function을 찾고 payload값을 넘겨주어 결과값을 받는다.
func LoadRpc(rpcKey string, payload string) string {
	mutex.Lock()
	defer mutex.Unlock()

	var result string
	rpcFunc, ok := GlobalGrpcEvent.Load(rpcKey)
	if !ok {
		result = "Not Found Rpc Key :" + rpcKey
		println(result)
		return result
	} else {
		if function, ok := rpcFunc.(RpcKeyHandlerFunc); ok {
			return function(payload)
		} else {
			result = "Result is not RpcKeyHandlerFunc Type"
			return result
		}
	}
}

type GrpcServer struct {
	global_grpc.UnimplementedGlobalGRpcServiceServer
}

func BackgroundGrpcServer() {
	grpcServer := grpc.NewServer()
	listen, err := net.Listen("tcp", SERVER_PORT)
	if err != nil {
		println("ListenError!!", err.Error())
		return
	}
	global_grpc.RegisterGlobalGRpcServiceServer(grpcServer, &GrpcServer{})
	if err := grpcServer.Serve(listen); err != nil {
		println("ListenGrpc Error!!::", err)
	}
}

func StartGrpcServer() {
	println("Start Login Server")
	GlobalGrpcEvent = new(sync.Map)

	go BackgroundGrpcServer()
}
