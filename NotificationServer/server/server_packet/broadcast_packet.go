package server_packet

type BroadcastHeartBeat struct {
	HeartBeat string `json:"heartBeat"`
}

type BroadcastDuplicateLogin struct {
	UUID string `json:"uuid"`
}
