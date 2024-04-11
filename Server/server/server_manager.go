package server

import (
	"context"
	"encoding/json"
	global_grpc "project_yd/grpc"
	packet "project_yd/server/server_packet"
	request "project_yd/server/server_packet/request_packet"
	"project_yd/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const SERVER_PORT = ":8080"
const SERVER_NAME = "game_server"

func SendBroadcastMessageToClient(client global_grpc.GlobalGRpcService_GlobalGrpcStreamBroadcastServer, opcode int32, message string) error {
	response := &global_grpc.GlobalGrpcBroadcast{
		Opcode:  opcode,
		Message: message,
	}
	if err := client.SendMsg(response); err != nil {
		return err
	}
	return nil
}

func BroadcastBaseMessage(client global_grpc.GlobalGRpcService_GlobalGrpcStreamBroadcastServer, opcode int32, broadcastPacket interface{}) {
	message, _ := json.Marshal(broadcastPacket)
	SendBroadcastMessageToClient(client, opcode, string(message))
}

var NotificationServer global_grpc.GlobalGRpcServiceClient

// -- NotificationServer 연결
func ConnectToNotificationServer() {
	println("ConnectToNotificationServer Start")
	address := util.NotificationIp + util.NotificationPort
	println("NotificationServer Address:", address)
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		println("Connect To NotificationServer Error!!")
		println(err.Error())
		return
	}
	//defer conn.Close()

	NotificationServer = global_grpc.NewGlobalGRpcServiceClient(conn)

	md := metadata.Pairs("server_name", SERVER_NAME)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := NotificationServer.GlobalGrpcStreamBroadcast(ctx)
	if err != nil {
		println("Connect To NotificationServer Error!!")
		println("stream error")
		println(err.Error())
		return
	}

	go ReceiveNotificationMessage(conn, stream)

	/*
		request := global_grpc.GlobalGrpcRequest{}
		request.RpcKey = ""
		request.Message = "ConnectTest"
		if err := stream.Send(&request); err != nil {
			println("can not send %v", err)
		}
	*/

	go func() {
		//for {
		state := conn.GetState()
		println("연결상태 :", state.String())
		//}
	}()
}

// -- NotificationServer로부터 수신된 메시지 요청 처리
func ReceiveNotificationMessage(conn *grpc.ClientConn, stream global_grpc.GlobalGRpcService_GlobalGrpcStreamBroadcastClient) {
	println("ReceiveNotificationMessage Start")
	for {
		listen, err := stream.Recv()
		if err != nil {
			if status.Code(err) == codes.Canceled {
				println("client disconnected")
			}
			return
		}
		//println(listen.Message)

		var opcode string
		switch listen.Opcode {
		case util.HEARTBEAT:
			opcode = "HeartBeat"
		case util.DUPLICATE_LOGIN:
			opcode = "DuplicateLogin"
		}
		println("Noti->Game:: Notification:: Opcode:", opcode, "/Message::", listen.Message)

		switch listen.Opcode {
		case util.HEARTBEAT:
			BroadcastHeartBeat()
		case util.DUPLICATE_LOGIN:
			requestPacket := request.DuplicateLoginFromNotificationServer{}
			err := json.Unmarshal([]byte(listen.Message), &requestPacket)
			if err != nil {
				println("ReceiveNotificationMessage Error")
				println("Opcode : DUPLICATE_LOGIN")
				println(err.Error())
				return
			}
			UUID := requestPacket.UUID
			BroadcastDuplicateLogin(UUID)
		}

		defer conn.Close()
	}
}

func ReceiveClientsMessage() {
	for {
		if len(BroadcastClients) <= 0 {
			continue
		}
		clientMessage, hasData := ClientMessageDequeue()
		if !hasData {
			continue
		}

		println("ClientMessage", clientMessage.Message)

		/*
			for UUID, client := range BroadcastClients {
				requestPacket, err := client.Recv()
				if err != nil {
					println("UUID:", UUID, "/ ReceiveClientsMessage Error")
					println(err.Error())
				}
				println("UUID:", UUID, "/RpcKey:", requestPacket.RpcKey)
				println("Message:", requestPacket.Message)
			}
		*/
	}
}

func ClientMessageEnqueue(data packet.UserMessage) {
	ClientMessageList.PushBack(data)
}
func ClientMessageDequeue() (packet.UserMessage, bool) {
	element := ClientMessageList.Front()
	if element == nil {
		// 리스트가 비어있을 때 처리할 로직 추가
		return packet.UserMessage{}, false
	}

	data := ClientMessageList.Remove(element)
	// data를 packet.UserMessage로 타입 어설션
	userMessage, ok := data.(packet.UserMessage)
	if !ok {
		// 타입 어설션 실패 시 처리할 로직 추가
		return packet.UserMessage{}, false
	}

	return userMessage, true
}

type Weapon struct {
	Id      int
	Enchant int
}
type Effect struct {
	Id    int
	Count int
}
type User struct {
	CurrentStage int
	Gold         int
	Slot         []Weapon
	Effect       []Effect
}

var Users map[string]*User
