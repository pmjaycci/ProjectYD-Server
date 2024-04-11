package main

import (
	login "project_yd/login"
	server "project_yd/server"
	"sync"
	//itemcreateevent "github.com/heroiclabs/nakama/v3/nurhyme_common/ItemCreateEvent"
)

func RegistRpc() {
	login.RegistLoginRpc()
}

func main() {
	var waitGroup sync.WaitGroup

	server.StartDBConnection()
	server.RedisConnection()

	//-- GoRoutine Count +1
	waitGroup.Add(1)
	go func() {
		//-- GoRoutine Count -1
		defer waitGroup.Done()
		//-- Start Grpc server
		server.StartGrpcServer()
		RegistRpc()
		//-- Notification서버 연결
		server.ConnectToNotificationServer()
	}()

	waitGroup.Wait()
	select {}
}
