package game

import (
	//"context"
	//"encoding/json"
	"encoding/json"
	"project_yd/server"
	request "project_yd/server/server_packet/request_packet"
	response "project_yd/server/server_packet/response_packet"
	"project_yd/util"
	//"project_yd/util"
)

func RegistSessionRpc() {
	server.RegistRpc("check_heartbeat", CheckHeartBeat)
}

//-- HeartBeat 구조
/*
1. 유저 로그인시 로그인서버에서 Redis에 HeartBeat발급하여 유저에게 리턴
2. 일정 시간마다 푸시서버에서 연결된 게임서버에 HeartBeat갱신 알림 전송
3. 수신한 게임서버는 연결된 클라이언트에 HeartBeat갱신 RPC호출 요청
4. 수신한 클라이언트는 게임서버에 현재 가지고 있는 HeartBeat값 게임서버에 전송
5. 수신한 게임서버는 Redis에 들고있는 HeartBeat값과 클라이언트로부터 수신된 HeartBeat값을 비교
6. 동일한 값일 경우 새로운 값으로 갱신후 리턴
*/
func CheckHeartBeat(UUID string, payload string) string {
	requestPacket := request.HeartBeat{}
	err := json.Unmarshal([]byte(payload), &requestPacket)
	if err != nil {
		return util.ResponseErrorMessage(util.BadRequest, err.Error())
	}

	isEquals, newHeartBeat := server.CheckEqualsHeartBeat(UUID, requestPacket.HeartBeat)

	if !isEquals {
		server.UnregisterBroadcastClient(UUID)
		return util.ResponseBaseMessage(util.BadRequest, "Not Equals HeartBeat")
	}
	responsePacket := response.HeartBeat{}
	responsePacket.Code = util.Success
	responsePacket.Message = "Success"
	responsePacket.HeartBeat = newHeartBeat

	return util.ResponseMessage(responsePacket)
}
