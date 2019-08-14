package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"time"
)

const DEFAULT_ADDRESS = ":55501"

type GrpcServer struct {
	grpcServer  *grpc.Server
	registrators []func(*grpc.Server)
	address     string
}

func NewGrpcServer(address string, registrators ...func(*grpc.Server)) *GrpcServer {
	return &GrpcServer{
		grpcServer:  grpc.NewServer(),
		registrators: registrators,
		address:     address,
	}
}

func (vs *GrpcServer) StartAsync(startupTimeout time.Duration) error {
	errChan := make(chan error)
	go func() {
		err := vs.StartSync()
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(startupTimeout):
		return nil
	}
}

func (vs *GrpcServer) StartSync() error {
	err, lis := vs.prepare()
	if err != nil {
		return err
	}

	if err := vs.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (vs *GrpcServer) Stop() {
	vs.grpcServer.GracefulStop()
}

func (vs *GrpcServer) prepare() (error, net.Listener) {
	address := vs.address
	if address == "" {
		address = DEFAULT_ADDRESS
	}

	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err, lis
	}
	fmt.Println("Starting grpc server")

	for _, reg := range vs.registrators {
		reg(vs.grpcServer)
	}
	return nil, lis
}
