package server

import (
	// fs1 "feature-store/api/feature-store/v1"
	pb "demo-temporal-project/api-protos/demo_temporal_project/v1"
	"demo-temporal-project/configs"
	"demo-temporal-project/internal/service"

	//  pb "demo-temporal-project/protos/path"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *configs.Server, logger log.Logger, transactionService *service.TransactionService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			tracing.Server(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Address != "" {
		opts = append(opts, grpc.Address(c.Grpc.Address))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(*c.Grpc.Timeout))
	}
	srv := grpc.NewServer(opts...)
	pb.RegisterTransactionServer(srv, transactionService)
	return srv
}
