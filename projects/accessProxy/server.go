package main

import (
	"context"
	"github.com/breathbath/uAssert/projects/accessProxy/protos/access_proxy"
	"github.com/breathbath/uAssert/simulation"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opencord/voltha-protos/go/voltha"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
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

	defer func(){
		err = conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()

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

func NewAccessProxyServer(accessProxyAddress, volthaAddress string) *simulation.GrpcServer {
	return simulation.NewGrpcServer(
		accessProxyAddress,
		func(server *grpc.Server) {
			access_proxy.RegisterDevicesServer(server, &ApDevicesServer{volthaAddress: volthaAddress})
		},
	)
}