package server

import (
	"context"
	global_grpc "project_yd/grpc"
	"project_yd/util"

	"google.golang.org/grpc/metadata"
)

// -- 클라이언트로부터 호출이 들어오면 rpcKey를 기준으로 등록된 function을 호출할것을 찾는다.
func (server *GrpcServer) GlobalGRpc(ctx context.Context, request *global_grpc.GlobalGrpcRequest) (*global_grpc.GlobalGrpcResponse, error) {
	//-- 메타 데이터에서 UUID추출
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		message := "메타데이터에서 클라이언트 ID를 추출할 수 없음"
		println(message)
		// 결과를 클라이언트에게 보냄
		response := &global_grpc.GlobalGrpcResponse{
			Message: message,
		}
		return response, nil
	}
	result := &global_grpc.GlobalGrpcResponse{}

	metaData := md.Get("UUID")
	var UUID string
	if len(metaData) <= 0 {
		message := util.ResponseErrorMessage(util.BadRequest, "UUID Error")
		result.Message = message
		return result, nil
	}

	UUID = metaData[0]
	println("Rpc Key:", request.RpcKey, " / Request UUID:", UUID)

	result.Message = LoadRpc(request.RpcKey, UUID, request.Message)

	return result, nil
}
