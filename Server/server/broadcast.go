package server

import (
	packet "project_yd/server/server_packet"
	util "project_yd/util"
)

func BroadcastHeartBeat() {
	if len(BroadcastClients) <= 0 {
		println("BroadcastClients Len size:", len(BroadcastClients))
		return
	}

	println("Broadcast HeartBeat:: Server->Clients")
	for _, client := range BroadcastClients {
		broadcastPacket := packet.BroadcastBase{}
		broadcastPacket.Message = "Please Request HeartBeat"

		BroadcastBaseMessage(client, util.HEARTBEAT, broadcastPacket)
	}
}

func BroadcastDuplicateLogin(UUID string) {
	client, isExist := BroadcastClients[UUID]
	if !isExist {
		println("BroadcastDuplicateLogin Not Exist UUID:", UUID)
		DeleteHeartBeat(UUID)
		return
	}
	broadcastPacket := packet.BroadcastBase{}
	broadcastPacket.Message = "Duplicate Login"
	BroadcastBaseMessage(client, util.DUPLICATE_LOGIN, broadcastPacket)
}
