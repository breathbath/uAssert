package server

import (
	"context"
	"github.com/breathbath/uAssert/projects/accessProxy/protos/access_proxy"
	"github.com/breathbath/uAssert/simulation"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opencord/voltha-protos/go/voltha"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	ACCESS_PROXY_ADDRESS = "locahost:52051"
)

type ApDevicesServer struct{
	volthaAddress string
}

func (as ApDevicesServer) GetDeviceBySn(ctx context.Context, sn *access_proxy.SerialNumber) (*voltha.Device, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(as.volthaAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	volthaClient := voltha.NewVolthaServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	devices, err := volthaClient.ListDevices(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	for _, device := range devices.Items {
		if sn.Sn == device.SerialNumber {
			return device, nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "%s", sn.Sn)
}

func GetAccessProxyServer() *simulation.GrpcServer {
	return simulation.NewGrpcServer(
		func(server *grpc.Server) {
			access_proxy.RegisterDevicesServer(server, &ApDevicesServer{})
		},
		ACCESS_PROXY_ADDRESS,
	)
}