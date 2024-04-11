package response_packet

type DuplicateLogin struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}
