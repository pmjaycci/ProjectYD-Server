package server_packet

type UserMessage struct {
	UUID    string
	RpcKey  string
	Message string
}

type BroadcastBase struct {
	Message string `json:"message"`
}
