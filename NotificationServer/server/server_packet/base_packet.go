package server_packet

type RequestPacket struct {
	Data interface{} `json:"data"`
}

type ResponsePacket struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

type BroadcastPacket struct {
	Message string `json:"message"`
}
