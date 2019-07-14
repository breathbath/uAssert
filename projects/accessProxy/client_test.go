package main

import (
	"context"
	"github.com/breathbath/uAssert/projects/accessProxy/protos/access_proxy"
	voltha2 "github.com/breathbath/uAssert/projects/voltha"
	"github.com/breathbath/uAssert/simulation"
	tests "github.com/breathbath/uAssert/test"
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

var testsRuntime *tests.Runtime

func init() {
	testsRuntime = tests.NewRuntime()
	testsRuntime.BeforeAll(setup)
	testsRuntime.AfterAll(cleanup)
	testsRuntime.TestCase(testDeviceIdSnMapping)
}

func setup(r *tests.Runtime) {
	volthaServerSimulator := voltha2.NewVolthaServerSimulator(VOLTHA_SERVER)
	err := volthaServerSimulator.StartAsync(time.Microsecond * 500)
	if err != nil {
		log.Panic(err)
	}
	r.SetState("voltha_server", volthaServerSimulator)

	accessProxyServer := NewAccessProxyServer(ACCESS_PROXY_SERVER, VOLTHA_SERVER)
	err = accessProxyServer.StartAsync(time.Microsecond * 500)
	if err != nil {
		log.Panic(err)
	}
	r.SetState("access_proxy_server", accessProxyServer)
}

func cleanup(r *tests.Runtime) {
	r.GetStateOrFail("voltha_server").(*simulation.GrpcServer).Stop()
	r.GetStateOrFail("access_proxy_server").(*simulation.GrpcServer).Stop()
}

func testDeviceIdSnMapping(t *testing.T, r *tests.Runtime) {
	grpcConn, err := grpc.Dial(ACCESS_PROXY_SERVER, grpc.WithInsecure())
	assert.NoError(t, err)

	accessProxyDevicesClient := access_proxy.NewDevicesClient(grpcConn)
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
			Address:         &voltha.Device_Ipv4Address{"11:111:111:11"},
			SerialNumber:    "sn2",
		}
		assert.Equal(t, expectedDevice, device)
	}
}

func TestAccessProxy(t *testing.T) {
	testsRuntime.Run(t)
}
