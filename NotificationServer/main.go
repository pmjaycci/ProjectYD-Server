package main

import (
	"project_yd/login"
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
		//-- Rpc 등록
		RegistRpc()
	}()

	waitGroup.Wait()
	select {}
}
