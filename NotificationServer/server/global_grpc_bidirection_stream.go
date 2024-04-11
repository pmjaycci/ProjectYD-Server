package server

import (
	"io"
	global_grpc "project_yd/grpc"
	"project_yd/util"
	"sync"

	"google.golang.org/grpc/metadata"
)

var Clients = make(map[string]global_grpc.GlobalGRpcService_GlobalGrpcStreamServer)
var BroadcastClients = make(map[string]global_grpc.GlobalGRpcService_GlobalGrpcStreamBroadcastServer)
var sessionMutex sync.Mutex

// -- 클라이언트 연결 시작 클라이언트의 스트림 저장
func RegisterClient(UUID string, client global_grpc.GlobalGRpcService_GlobalGrpcStreamServer) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	_, ok := Clients[UUID]
	if !ok {
		Clients[UUID] = client
	}
}

// -- 클라이언트 연결 시작 클라이언트의 스트림 저장
func RegisterBroadcastClient(UUID string, client global_grpc.GlobalGRpcService_GlobalGrpcStreamBroadcastServer) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	BroadcastClients[UUID] = client
	println("RegisterBroadcast Client UUID:", UUID, "[Clients Size:", len(BroadcastClients), "]")
}

// -- 클라이언트 연결 종료 클라이언트 스트림 제거
func UnregisterClient(UUID string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	_, ok := Clients[UUID]
	if ok {
		delete(Clients, UUID)
	} else {
		println("UnregisterClient:: Error")
		println("존재하지 않는데 삭제 호출 들어옴")
	}
}

func UnregisterBroadcastClient(UUID string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	_, ok := BroadcastClients[UUID]
	if ok {
		delete(BroadcastClients, UUID)
		println("Delete BroadcastClient UUID:", UUID)
	} else {
		println("UnregisterBroadcastClient:: Error")
		println("존재하지 않는데 삭제 호출 들어옴")
	}
}
func (server *GrpcServer) GlobalGrpcStream(client global_grpc.GlobalGRpcService_GlobalGrpcStreamServer) error {
	println("GlobalGrpcStream")
	return nil
}

// -- Grpc BiStream 통신 처리
func (server *GrpcServer) GlobalGrpcStreamBroadcast(client global_grpc.GlobalGRpcService_GlobalGrpcStreamBroadcastServer) error {
	println("GlobalGrpcStreamBroadcast called")

	//-- 메타 데이터에서 UUID추출
	md, ok := metadata.FromIncomingContext(client.Context())
	if !ok {
		message := "메타데이터에서 클라이언트 ID를 추출할 수 없음"
		println(message)
		// 결과를 클라이언트에게 보냄
		response := &global_grpc.GlobalGrpcBroadcast{
			Opcode:  util.ERROR,
			Message: message,
		}
		if err := client.SendMsg(response); err != nil {
			return err
		}
	}

	customHeaderServerName := md.Get("server_name")
	var serverName string
	if len(customHeaderServerName) > 0 {
		serverName = customHeaderServerName[0]
		println("Received CustomHeader UUID:", serverName)
		RegisterBroadcastClient(serverName, client)
		defer UnregisterBroadcastClient(serverName)
	}

	for {
		//-- 클라이언트로부터 메시지를 받음
		request, err := client.Recv()
		if err == io.EOF {
			println("GlobalGRpcStreamBroadcast:: io.EOF Error")
			println(err.Error())
			return nil
		}
		if err != nil {
			println("GlobalGRpcStreamBroadcast:: Error")
			println(err.Error())
			return err
		}
		println("message", request.Message)
		//-- 유저가 호출했을때만 처리
		if request.RpcKey != "" {
			//RpcBroadcastMessage(UUID, client, request)
		}

	}
}

func RpcBroadcastMessage(UUID string, client global_grpc.GlobalGRpcService_GlobalGrpcStreamBroadcastServer, request *global_grpc.GlobalGrpcRequest) {
	// rpcKey와 message를 사용하여 결과를 생성
	result, opcode := LoadRpcStream(request.RpcKey, UUID, request.Message)

	// 결과를 클라이언트에게 보냄
	response := &global_grpc.GlobalGrpcBroadcast{
		Opcode:  opcode,
		Message: result,
	}
	if err := client.SendMsg(response); err != nil {
		println("RpcResponseMessage:: Error")
		println(err.Error())
		return
	}
}

//-- 연결된 게임서버에 보낼 메세지들
func BroadcastMessages() {
	println("BroadcastMessage Start")
	//-- 하트비트
	go SendHeartBeat()
}
