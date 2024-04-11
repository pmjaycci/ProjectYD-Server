package response_packet

type Login struct {
	Code      uint   `json:"code"`
	Message   string `json:"message"`
	UUID      string `json:"UUID"`
	UserName  string `json:"userName"`
	HeartBeat string `json:"heartBeat"`
}

type DuplicateLogin struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}
