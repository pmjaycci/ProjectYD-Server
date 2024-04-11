package server

import (
	"io"
	global_grpc "project_yd/grpc"
	packet "project_yd/server/server_packet"
	"project_yd/util"
	"sync"

	"google.golang.org/grpc/metadata"
)

var BroadcastClients = make(map[string]global_grpc.GlobalGRpcService_GlobalGrpcStreamBroadcastServer)
var sessionMutex sync.Mutex

// -- 클라이언트 연결 시작 클라이언트의 스트림 저장
func RegisterBroadcastClient(UUID string, client global_grpc.GlobalGRpcService_GlobalGrpcStreamBroadcastServer) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()

	_, ok := BroadcastClients[UUID]
	if !ok {
		BroadcastClients[UUID] = client
	}
}

// -- 클라이언트 연결 종료 클라이언트 스트림 제거
func UnregisterBroadcastClient(UUID string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	_, ok := BroadcastClients[UUID]
	if ok {
		println("DeleteBroadcastClient::", UUID)
		delete(BroadcastClients, UUID)
		DeleteHeartBeat(UUID)
		delete(Users, UUID)
	} else {
		println("UnregisterBroadcastClient:: Error")
		println("존재하지 않는데 삭제 호출 들어옴")
	}
}

// -- Grpc BiStream 통신 처리
func (server *GrpcServer) GlobalGRpcStream(client global_grpc.GlobalGRpcService_GlobalGrpcStreamServer) error {
	return nil
}

/*
Stream으로 클라이언트에서 RPC를 호출할때마다 해당 메시지는
Queue(List)에 담고 고루틴으로 돌고있는 ReceiveClientsMessage 함수에서
담긴 모든 메시지들을 처리한다. Rpc키의 따라 연결된 모든 클라이언트들에게 메세지들을 전달할지
메시지를 보낸 유저 혹은 특정 유저에게 메시지를 전달할지의 대해서는 구현한 Rpc Function에서 결정한다.
*/
func (server *GrpcServer) GlobalGrpcStreamBroadcast(client global_grpc.GlobalGRpcService_GlobalGrpcStreamBroadcastServer) error {
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

	customHeaderUUID := md.Get("uuid")
	var UUID string
	if len(customHeaderUUID) > 0 {
		UUID = customHeaderUUID[0]
		println("Received CustomHeader UUID:", UUID)
		RegisterBroadcastClient(UUID, client)
		//defer UnregisterBroadcastClient(UUID)
	}

	//-- 클라이언트로부터 메시지를 받음
	request, err := client.Recv()
	if err == io.EOF {
		println("GlobalGRpcStreamBroadcast:: io.EOF Error")
		println("UUID:", UUID)
		println(err.Error())
		UnregisterBroadcastClient(UUID)
		return nil
	}
	if err != nil {
		println("GlobalGRpcStreamBroadcast:: Error")
		println("UUID:", UUID)
		println(err.Error())
		UnregisterBroadcastClient(UUID)
		return err
	}
	if request.RpcKey != "" {
		mutex.Lock()
		defer mutex.Unlock()
		data := packet.UserMessage{}
		data.UUID = UUID
		data.RpcKey = request.RpcKey
		data.Message = request.Message

		ClientMessageEnqueue(data)
	}

	return nil
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
