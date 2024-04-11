package server

import (
	"container/list"
	"net"
	global_grpc "project_yd/grpc"
	"project_yd/util"
	"sync"

	"google.golang.org/grpc"
)

var GlobalGrpcEvent *sync.Map
var GlobalGrpcStreamEvent *sync.Map
var mutex sync.Mutex
var ClientMessageList list.List

// 함수 시그니처: 클라이언트로부터 받은 key 값에 따라 호출할 함수들
type RpcKeyHandlerFunc func(UUID string, payload string) string
type RpcStreamBroadcastKeyHandlerFunc func(UUID string, payload string) (string, int32)

// #region ** UNARY RPC **
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
func LoadRpc(rpcKey string, UUID string, payload string) string {
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
			return function(UUID, payload)
		} else {
			result = "Result is not RpcKeyHandlerFunc Type"
			return result
		}
	}
}

//#endregion

// #region ** BIDIRECT STREAM RPC **
// -- rpc로 쓰일 function 등록
func RegistRpcStream(rpcKey string, function RpcStreamBroadcastKeyHandlerFunc) {
	_, ok := GlobalGrpcStreamEvent.Load(rpcKey)
	if ok {
		println("RPC KEY : ", rpcKey, " is Duplicate!")
		return
	}
	GlobalGrpcStreamEvent.Store(rpcKey, function)
	println("RPC KEY : ", rpcKey, " is Regist Success!")
}

// -- 클라이언트로부터 호출된 rpc명으로 호출할 function을 찾고 payload값을 넘겨주어 결과값을 받는다.
func LoadRpcStream(rpcKey string, UUID string, payload string) (string, int32) {
	mutex.Lock()
	defer mutex.Unlock()

	var result string
	rpcFunc, ok := GlobalGrpcStreamEvent.Load(rpcKey)
	if !ok {
		result = "Not Found RpcStream Key :" + rpcKey
		println(result)
		return result, util.ERROR
	} else {
		if function, ok := rpcFunc.(RpcStreamBroadcastKeyHandlerFunc); ok {
			return function(UUID, payload)
		} else {
			result = "Result is not RpcStreamKeyHandlerFunc Type"
			return result, util.ERROR
		}
	}
}

//#endregion

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

	go ReceiveClientsMessage()
}

func StartGrpcServer() {
	println("Start GRPC Server")
	GlobalGrpcEvent = new(sync.Map)
	GlobalGrpcStreamEvent = new(sync.Map)
	ClientMessageList = *list.New()
	Users = make(map[string]*User)
	go BackgroundGrpcServer()

}
