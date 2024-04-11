package server

import (
	"context"
	global_grpc "project_yd/grpc"
)

// -- 클라이언트로부터 호출이 들어오면 rpcKey를 기준으로 등록된 function을 호출할것을 찾는다.
func (server *GrpcServer) GlobalGRpc(ctx context.Context, request *global_grpc.GlobalGrpcRequest) (*global_grpc.GlobalGrpcResponse, error) {
	result := &global_grpc.GlobalGrpcResponse{}
	result.Message = LoadRpc(request.RpcKey, request.Message)

	return result, nil
}
