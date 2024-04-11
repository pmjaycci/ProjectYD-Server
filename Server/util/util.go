package util

import (
	"encoding/json"
	packet "project_yd/server/server_packet"
)

func ResponseErrorMessage(messageCode uint, errMessage string) string {
	println("Response Error!!", errMessage)
	responsePacket := packet.ResponsePacket{}
	responsePacket.Code = messageCode
	responsePacket.Message = errMessage

	result, _ := json.Marshal(responsePacket)
	return string(result)
}

func ResponseBaseMessage(messageCode uint, message string) string {
	responsePacket := packet.ResponsePacket{}
	responsePacket.Code = messageCode
	responsePacket.Message = message

	result, _ := json.Marshal(responsePacket)
	return string(result)
}

func ResponseMessage(responsePacket interface{}) string {
	result, _ := json.Marshal(responsePacket)
	return string(result)
}
