package server

import (
	"encoding/json"
	global_grpc "project_yd/grpc"
	packet "project_yd/server/server_packet"
	"project_yd/util"
	"time"
)

const SERVER_PORT = ":8082"

func SendBroadcastMessageToClient(client global_grpc.GlobalGRpcService_GlobalGrpcStreamBroadcastServer, opcode int32, message string) error {
	response := &global_grpc.GlobalGrpcBroadcast{
		Opcode:  opcode,
		Message: message,
	}
	if err := client.SendMsg(response); err != nil {
		return err
	}
	return nil
}

func BroadcastBaseMessage(client global_grpc.GlobalGRpcService_GlobalGrpcStreamBroadcastServer, opcode int32, broadcastPacket interface{}) {
	message, _ := json.Marshal(broadcastPacket)
	SendBroadcastMessageToClient(client, opcode, string(message))
}

func SendDuplicateLogin(UUID string) {
	for _, client := range BroadcastClients {
		broadcastPacket := packet.BroadcastDuplicateLogin{}
		broadcastPacket.UUID = UUID
		BroadcastBaseMessage(client, util.DUPLICATE_LOGIN, broadcastPacket)
	}
}

//-- 연결된 클라이언트(게임서버)의 수가 0 ㅇ
func SendHeartBeat() {
	for {
		println("Broadcast -- SendHeartBeat")
		if len(BroadcastClients) > 0 {
			for _, client := range BroadcastClients {
				broadcastPacket := packet.BroadcastPacket{}
				broadcastPacket.Message = "Send Heart Beat"
				BroadcastBaseMessage(client, util.HEARTBEAT, broadcastPacket)
			}
		}
		// 60초 대기
		time.Sleep(12 * (3600 * time.Second))
	}
}
