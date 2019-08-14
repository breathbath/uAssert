package voltha

import (
	"github.com/breathbath/uAssert/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opencord/voltha-protos/go/voltha"
)

func GetStubs() grpc.SimulationCases {
	return grpc.SimulationCases{
		{
			Request: &empty.Empty{},
			Response: &voltha.Devices{
				Items: []*voltha.Device{
					{
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
					},
				},
			},
			GrpcNamespace: "/voltha.VolthaService/ListDevices",
		},
	}
}
