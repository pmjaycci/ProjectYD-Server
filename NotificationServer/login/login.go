package login

import (
	"encoding/json"
	server "project_yd/server"
	request "project_yd/server/server_packet/request_packet"
	response "project_yd/server/server_packet/response_packet"
	"project_yd/util"
)

func RegistLoginRpc() {
	server.RegistRpc("duplicate_login", DuplicateLoginRpc)
}

func DuplicateLoginRpc(payload string) string {
	println("DuplicateLoginRpc")
	println("payload:", payload)
	requestPacket := request.DuplicateLogin{}
	err := json.Unmarshal([]byte(payload), &requestPacket)
	if err != nil {
		return util.ResponseErrorMessage(util.BadRequest, err.Error())
	}

	//-- Stream연결된 서버들에게 해당 유저 UUID 전송
	server.SendDuplicateLogin(requestPacket.UUID)

	responsePacket := response.DuplicateLogin{}
	responsePacket.Code = util.Success
	responsePacket.Message = "Sueccess"

	return util.ResponseMessage(responsePacket)
}
