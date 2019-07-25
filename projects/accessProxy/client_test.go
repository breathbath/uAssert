package main

import (
	"context"
	"github.com/breathbath/uAssert/projects/accessProxy/protos/access_proxy"
	voltha2 "github.com/breathbath/uAssert/projects/voltha"
	"github.com/breathbath/uAssert/simulation"
	"github.com/opencord/voltha-protos/go/voltha"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"log"
	"testing"
	"time"
)

const (
	VOLTHA_SERVER       = "localhost:23456"
	ACCESS_PROXY_SERVER = "localhost:23457"
)

var (
	volthaServerSimulator    *simulation.GrpcServer
	accessProxyServer        *simulation.GrpcServer
	accessProxyDevicesClient access_proxy.DevicesClient
)

func setup() {
	volthaServerSimulator = voltha2.NewVolthaServerSimulator(VOLTHA_SERVER)
	err := volthaServerSimulator.StartAsync(time.Microsecond * 500)
	if err != nil {
		log.Panic(err)
	}

	accessProxyServer = NewAccessProxyGrpcServer(ACCESS_PROXY_SERVER, VOLTHA_SERVER)
	err = accessProxyServer.StartAsync(time.Microsecond * 500)
	if err != nil {
		log.Panic(err)
	}

	grpcConn, err := grpc.Dial(ACCESS_PROXY_SERVER, grpc.WithInsecure())
	if err != nil {
		log.Panic(err)
	}
	accessProxyDevicesClient = access_proxy.NewDevicesClient(grpcConn)
}

func cleanup() {
	volthaServerSimulator.Stop()
	accessProxyServer.Stop()
}

func testDeviceIdSnMapping(t *testing.T) {
	device, err := accessProxyDevicesClient.GetDeviceBySn(
		context.Background(),
		&access_proxy.SerialNumber{Sn: "sn2"},
		grpc.WaitForReady(true),
	)
	assert.NoError(t, err)
	if device != nil {
		expectedDevice := &voltha.Device{
			Id:              "id1",
			Type:            "Olt",
			Root:            true,
			ParentId:        "",
			ParentPortNo:    22,
			Vendor:          "Some",
			Model:           "xyw",
			HardwareVersion: "333",
			FirmwareVersion: "333",
			Address:         &voltha.Device_Ipv4Address{Ipv4Address: "11:111:111:11"},
			SerialNumber:    "sn2",
		}
		assert.Equal(t, expectedDevice, device)
	}
}

func testNoDeviceBySnIsFound(t *testing.T) {
	device, err := accessProxyDevicesClient.GetDeviceBySn(
		context.Background(),
		&access_proxy.SerialNumber{Sn: "some_unknown_sn"},
		grpc.WaitForReady(true),
	)
	assert.EqualError(t, err, "rpc error: code = NotFound desc = some_unknown_sn")
	assert.Nil(t, device)
}

func TestGrpcAccessProxy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	setup()
	defer cleanup()
	t.Run("testDeviceIdSnMapping", testDeviceIdSnMapping)
	t.Run("testNoDeviceBySnIsFound", testNoDeviceBySnIsFound)
}
