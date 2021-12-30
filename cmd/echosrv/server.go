package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/go-stack/stack"
	"github.com/karlmutch/kv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"

	echo "github.com/karlmutch/platform-services/internal/gen/echosrv"
)

type echoServer struct {
	echo.UnimplementedEchoServer
	health *health.Server
}

func (*echoServer) Echo(ctx context.Context, in *echo.Request) (resp *echo.Response, err error) {
	if in == nil {
		return nil, fmt.Errorf("request is missing a message to echo")
	}

	return &echo.Response{
		Message:  in.Message,
		DateTime: &timestamp.Timestamp{Seconds: time.Now().Unix()}}, nil
}

func (es *echoServer) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (resp *grpc_health_v1.HealthCheckResponse, err error) {
	return es.health.Check(ctx, in)
}

func (*echoServer) Watch(in *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) (err error) {
	return kv.NewError(grpc_health_v1.HealthCheckResponse_UNKNOWN.String())
}

func runServer(ctx context.Context, serviceName string, port int) (errC chan kv.Error) {

	errC = make(chan kv.Error, 3)

	server := grpc.NewServer()
	echoSrv := &echoServer{health: health.NewServer()}

	echo.RegisterEchoServer(server, echoSrv)
	grpc_health_v1.RegisterHealthServer(server, echoSrv)

	reflection.Register(server)

	netListen, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		errC <- kv.Wrap(err).With("stack", stack.Trace().TrimRuntime())
		return
	}

	go func() {
		echoSrv.health.SetServingStatus(serviceName, grpc_health_v1.HealthCheckResponse_SERVING)

		// serve API
		if err := server.Serve(netListen); err != nil {
			errC <- kv.Wrap(err).With("stack", stack.Trace().TrimRuntime())
		}
		echoSrv.health.SetServingStatus(serviceName, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		func() {
			defer recover()
			close(errC)
		}()
	}()

	go func() {
		select {
		case <-ctx.Done():
			server.Stop()
		}
	}()
	return errC
}
