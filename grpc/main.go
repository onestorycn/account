package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"account/grpc/handler"
	"github.com/micro/go-grpc"
	account "account/grpc/proto/account"
	"context"
	"github.com/micro/go-micro/server"
	"account/library/logger"
	"time"
	"go.uber.org/zap"
)

func main() {

	// New Service
	service := grpc.NewService(
		micro.Name("onestory.account.proto"),
		micro.Version("latest"),
		micro.WrapHandler(logWrapper),
	)
	
	// Initialise service
	service.Init()

	// Register Handler
	account.RegisterAccountServiceHandler(service.Server(), new(handler.Account))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		startTime := time.Now().UnixNano()
		res := fn(ctx, req, rsp)
		endTime := time.Now().UnixNano()

		logger.ZapTrace.Info("end request", zap.String("method", req.Method()), zap.Int64("request_time/ms", (endTime-startTime)/1e6))
		return res
	}
}
