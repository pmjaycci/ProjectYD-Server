package login

import (
	"context"
	"database/sql"
	"encoding/json"
	server "project_yd/server"
	"project_yd/server/server_packet/request_packet"
	request "project_yd/server/server_packet/request_packet"
	"project_yd/server/server_packet/response_packet"
	response "project_yd/server/server_packet/response_packet"
	"project_yd/util"

	"github.com/google/uuid"
)

func RegistLoginRpc() {
	server.RegistRpc("login", LoginRpc)
}

func LoginRpc(payload string) string {
	requestPacket := request.Login{}
	err := json.Unmarshal([]byte(payload), &requestPacket)
	if err != nil {
		return util.ResponseErrorMessage(util.BadRequest, err.Error())
	}
	println("TEST:: payload:", payload)
	var UUID string
	var userName string
	ctx := context.Background()
	//-- 로그인 정보 가져오기
	db := server.DBManager.Login
	loginSql := `SELECT uid, user_name FROM account WHERE user_id = ?`
	err = db.QueryRowContext(ctx, loginSql, requestPacket.Id).Scan(&UUID, &userName)
	//-- 등록된 정보가 없으므로 신규 유저로 처리
	responsePacket := response.Login{}
	if err == sql.ErrNoRows {
		//-- uuid 생성
		uid := CreateNameUUID(requestPacket.Id)
		createUUIDSql := `INSERT INTO account (uid, user_id, money) VALUES(?,?, 10000)`
		_, err := db.ExecContext(ctx, createUUIDSql, uid, requestPacket.Id)
		if err != nil {
			return util.ResponseErrorMessage(util.ServerError, err.Error())
		}

		heartBeat, err := server.SetHeartBeat(UUID)
		if err != nil {
			return util.ResponseErrorMessage(util.ServerError, err.Error())
		}

		//-- 인벤토리 기본템 추가
		gameDB := server.DBManager.Game
		invenSql := `INSERT INTO inventory (uid, item_id, item_count, enchant_level) VALUES(?, 0, 1, 0)`
		_, invenErr := gameDB.ExecContext(ctx, invenSql, uid)
		if invenErr != nil {
			return util.ResponseErrorMessage(util.ServerError, err.Error())
		}

		responsePacket.UUID = uid
		responsePacket.HeartBeat = heartBeat
		responsePacket.Message = "Success"
		responsePacket.Code = util.Success

		return util.ResponseMessage(responsePacket)
	}
	if err != nil {
		return util.ResponseErrorMessage(util.ServerError, err.Error())
	}
	println("DB Get UUID:", UUID)

	//CheckDuplicateLogin(UUID)
	//-- Redis에서 HeartBeat Key가 존재하는지 체크후 존재할경우 중복로그인 처리
	if server.HasHeartBeat(UUID) {
		println("DuplicateLogin UUID:", UUID)
		DuplicateLogin(UUID)
		return util.ResponseErrorMessage(util.Conflict, "Duplicate Login")
	}

	heartBeat, err := server.SetHeartBeat(UUID)
	println("UUID::", UUID, "/SetHeartBeat::", heartBeat)
	if err != nil {
		return util.ResponseErrorMessage(util.ServerError, err.Error())
	}

	responsePacket.UUID = UUID
	responsePacket.UserName = userName
	responsePacket.HeartBeat = heartBeat
	responsePacket.Code = util.Success
	responsePacket.Message = "Sueccess"

	return util.ResponseMessage(responsePacket)
}

// -- 랜덤 uuid 생성
func CreateRandomUUID() string {
	uid := uuid.New().String()
	println("CreatRandomUUID:: uuid:", uid)
	return uid
}

// -- 이름기반 uuid 생성
func CreateNameUUID(name string) string {
	baseUUID := uuid.New()
	nameByte := []byte(name)
	nameUUID := uuid.NewSHA1(baseUUID, nameByte)
	println("CreateNameUUID:: name:", name, "/uuid:", nameUUID.String())
	return nameUUID.String()
}

func DuplicateLogin(UUID string) {
	requestPacket := request_packet.DuplicateLogin{}
	requestPacket.UUID = UUID

	payload := server.GlobalGrpcMessageToNotificationServer("duplicate_login", requestPacket)
	println("DuplicateLogin payload:", payload)
	responsePacket := response_packet.DuplicateLogin{}

	err := json.Unmarshal([]byte(payload), &responsePacket)
	if err != nil {
		println("DuplicationLogin Json Unmarshal Error")
		println(err.Error())
	}
	if responsePacket.Code != util.Success {
		println("DuplicateLogin Error")
		println(responsePacket.Message)
	}
}
