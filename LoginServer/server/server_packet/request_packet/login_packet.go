package request_packet

type Login struct {
	Id string `json:"id"`
}

type DuplicateLogin struct {
	UUID string `json:"uuid"`
}
