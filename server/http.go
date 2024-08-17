package server

import (
	pb "demo-temporal-project/api-protos/demo_temporal_project/v1"
	"demo-temporal-project/configs"
	"demo-temporal-project/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *configs.Server, logger log.Logger, transactionService *service.TransactionService) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			tracing.Server(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Address != "" {
		opts = append(opts, http.Address(c.Http.Address))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(*c.Http.Timeout))
	}
	srv := http.NewServer(opts...)
	pb.RegisterTransactionHTTPServer(srv, transactionService)
	return srv
}
