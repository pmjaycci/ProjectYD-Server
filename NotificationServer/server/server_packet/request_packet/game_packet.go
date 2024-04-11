package request_packet

type CheckHeartBeat struct {
	HeartBeat string `json:"heartBeat"`
}

type PingPong struct {
	Ping string `json:"ping"`
}
