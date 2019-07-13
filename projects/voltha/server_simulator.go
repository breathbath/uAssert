package voltha

import (
	"context"
	"errors"
	"fmt"
	"github.com/breathbath/uAssert/simulation"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opencord/voltha-protos/go/common"
	"github.com/opencord/voltha-protos/go/omci"
	"github.com/opencord/voltha-protos/go/openflow_13"
	"github.com/opencord/voltha-protos/go/voltha"
	"google.golang.org/grpc"
)

var simCasesMap simulation.GrpcCasesMap

func init() {
	simCasesMap = simulation.GetSimulationMap(GetStubs())
}

type TestServer struct {
}

func (s *TestServer) UpdateLogLevel(context.Context, *voltha.Logging) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Get the membership group of a Voltha Core
func (s *TestServer) GetMembership(ctx context.Context, in *empty.Empty) (*voltha.Membership, error) {
	return nil, errors.New("UnImplemented")
}

// Set the membership group of a Voltha Core
func (s *TestServer) UpdateMembership(ctx context.Context, in *voltha.Membership) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Get high level information on the Voltha cluster
func (s *TestServer) GetVoltha(ctx context.Context, in *empty.Empty) (*voltha.Voltha, error) {
	return nil, errors.New("UnImplemented")
}

// List all Voltha cluster core instances
func (s *TestServer) ListCoreInstances(ctx context.Context, in *empty.Empty) (*voltha.CoreInstances, error) {
	return nil, errors.New("UnImplemented")
}

// Get details on a Voltha cluster instance
func (s *TestServer) GetCoreInstance(ctx context.Context, in *common.ID) (*voltha.CoreInstance, error) {
	return nil, errors.New("UnImplemented")
}

// List all active adapters (plugins) in the Voltha cluster
func (s *TestServer) ListAdapters(ctx context.Context, in *empty.Empty) (*voltha.Adapters, error) {
	return nil, errors.New("UnImplemented")
}

// List all logical devices managed by the Voltha cluster
func (s *TestServer) ListLogicalDevices(ctx context.Context, in *empty.Empty) (*voltha.LogicalDevices, error) {
	return nil, errors.New("UnImplemented")
}

// Get additional information on a given logical device
func (s *TestServer) GetLogicalDevice(ctx context.Context, in *common.ID) (*voltha.LogicalDevice, error) {
	return nil, errors.New("UnImplemented")
}

// List ports of a logical device
func (s *TestServer) ListLogicalDevicePorts(ctx context.Context, in *common.ID) (*voltha.LogicalPorts, error) {
	return nil, errors.New("UnImplemented")
}

// Gets a logical device port
func (s *TestServer) GetLogicalDevicePort(ctx context.Context, in *voltha.LogicalPortId) (*voltha.LogicalPort, error) {
	return nil, errors.New("UnImplemented")
}

// Enables a logical device port
func (s *TestServer) EnableLogicalDevicePort(ctx context.Context, in *voltha.LogicalPortId) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Disables a logical device port
func (s *TestServer) DisableLogicalDevicePort(ctx context.Context, in *voltha.LogicalPortId) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// List all flows of a logical device
func (s *TestServer) ListLogicalDeviceFlows(ctx context.Context, in *common.ID) (*openflow_13.Flows, error) {
	return nil, errors.New("UnImplemented")
}

// Update flow table for logical device
func (s *TestServer) UpdateLogicalDeviceFlowTable(ctx context.Context, in *openflow_13.FlowTableUpdate) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Update meter table for logical device
func (s *TestServer) UpdateLogicalDeviceMeterTable(ctx context.Context, in *openflow_13.MeterModUpdate) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Get all meter stats for logical device
func (s *TestServer) GetMeterStatsOfLogicalDevice(ctx context.Context, in *common.ID) (*openflow_13.MeterStatsReply, error) {
	return nil, errors.New("UnImplemented")
}

// List all flow groups of a logical device
func (s *TestServer) ListLogicalDeviceFlowGroups(ctx context.Context, in *common.ID) (*openflow_13.FlowGroups, error) {
	return nil, errors.New("UnImplemented")
}

// Update group table for device
func (s *TestServer) UpdateLogicalDeviceFlowGroupTable(ctx context.Context, in *openflow_13.FlowGroupTableUpdate) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// List all physical devices controlled by the Voltha cluster
func (s *TestServer) ListDevices(ctx context.Context, in *empty.Empty) (*voltha.Devices, error) {
	c, found := simulation.FindSimulatedCaseForRequest(in, "/voltha.VolthaService/ListDevices", simCasesMap)
	if !found {
		return nil, fmt.Errorf("Not found")
	}

	return c.Response.(*voltha.Devices), nil
}

// List all physical devices IDs controlled by the Voltha cluster
func (s *TestServer) ListDeviceIds(ctx context.Context, in *empty.Empty) (*common.IDs, error) {
	return nil, errors.New("UnImplemented")
}

// Request to a voltha Core to reconcile a set of devices based on their IDs
func (s *TestServer) ReconcileDevices(ctx context.Context, in *common.IDs) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Get more information on a given physical device
func (s *TestServer) GetDevice(ctx context.Context, in *common.ID) (*voltha.Device, error) {
	return nil, errors.New("UnImplemented")
}

// Pre-provision a new physical device
func (s *TestServer) CreateDevice(ctx context.Context, in *voltha.Device) (*voltha.Device, error) {
	return nil, errors.New("UnImplemented")
}

// Enable a device.  If the device was in pre-provisioned state then it
// will transition to ENABLED state.  If it was is DISABLED state then it
// will transition to ENABLED state as well.
func (s *TestServer) EnableDevice(ctx context.Context, in *common.ID) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Disable a device
func (s *TestServer) DisableDevice(ctx context.Context, in *common.ID) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Reboot a device
func (s *TestServer) RebootDevice(ctx context.Context, in *common.ID) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Delete a device
func (s *TestServer) DeleteDevice(ctx context.Context, in *common.ID) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Request an image download to the standby partition
// of a device.
// Note that the call is expected to be non-blocking.
func (s *TestServer) DownloadImage(ctx context.Context, in *voltha.ImageDownload) (*common.OperationResp, error) {
	return nil, errors.New("UnImplemented")
}

// Get image download status on a device
// The request retrieves progress on device and updates db record
func (s *TestServer) GetImageDownloadStatus(ctx context.Context, in *voltha.ImageDownload) (*voltha.ImageDownload, error) {
	return nil, errors.New("UnImplemented")
}

// Get image download db record
func (s *TestServer) GetImageDownload(ctx context.Context, in *voltha.ImageDownload) (*voltha.ImageDownload, error) {
	return nil, errors.New("UnImplemented")
}

// List image download db records for a given device
func (s *TestServer) ListImageDownloads(ctx context.Context, in *common.ID) (*voltha.ImageDownloads, error) {
	return nil, errors.New("UnImplemented")
}

// Cancel an existing image download process on a device
func (s *TestServer) CancelImageDownload(ctx context.Context, in *voltha.ImageDownload) (*common.OperationResp, error) {
	return nil, errors.New("UnImplemented")
}

// Activate the specified image at a standby partition
// to active partition.
// Depending on the device implementation, this call
// may or may not cause device reboot.
// If no reboot, then a reboot is required to make the
// activated image running on device
// Note that the call is expected to be non-blocking.
func (s *TestServer) ActivateImageUpdate(ctx context.Context, in *voltha.ImageDownload) (*common.OperationResp, error) {
	return nil, errors.New("UnImplemented")
}

// Revert the specified image at standby partition
// to active partition, and revert to previous image
// Depending on the device implementation, this call
// may or may not cause device reboot.
// If no reboot, then a reboot is required to make the
// previous image running on device
// Note that the call is expected to be non-blocking.
func (s *TestServer) RevertImageUpdate(ctx context.Context, in *voltha.ImageDownload) (*common.OperationResp, error) {
	return nil, errors.New("UnImplemented")
}

// List ports of a device
func (s *TestServer) ListDevicePorts(ctx context.Context, in *common.ID) (*voltha.Ports, error) {
	return nil, errors.New("UnImplemented")
}

// List pm config of a device
func (s *TestServer) ListDevicePmConfigs(ctx context.Context, in *common.ID) (*voltha.PmConfigs, error) {
	return nil, errors.New("UnImplemented")
}

// Update the pm config of a device
func (s *TestServer) UpdateDevicePmConfigs(ctx context.Context, in *voltha.PmConfigs) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// List all flows of a device
func (s *TestServer) ListDeviceFlows(ctx context.Context, in *common.ID) (*openflow_13.Flows, error) {
	return nil, errors.New("UnImplemented")
}

// List all flow groups of a device
func (s *TestServer) ListDeviceFlowGroups(ctx context.Context, in *common.ID) (*openflow_13.FlowGroups, error) {
	return nil, errors.New("UnImplemented")
}

// List device types known to Voltha
func (s *TestServer) ListDeviceTypes(ctx context.Context, in *empty.Empty) (*voltha.DeviceTypes, error) {
	return nil, errors.New("UnImplemented")
}

// Get additional information on a device type
func (s *TestServer) GetDeviceType(ctx context.Context, in *common.ID) (*voltha.DeviceType, error) {
	return nil, errors.New("UnImplemented")
}

// List all device sharding groups
func (s *TestServer) ListDeviceGroups(ctx context.Context, in *empty.Empty) (*voltha.DeviceGroups, error) {
	return nil, errors.New("UnImplemented")
}

// Stream control packets to the dataplane
func (s *TestServer) StreamPacketsOut(voltha.VolthaService_StreamPacketsOutServer) error {
	return errors.New("UnImplemented")
}

// Receive control packet stream
func (s *TestServer) ReceivePacketsIn(*empty.Empty, voltha.VolthaService_ReceivePacketsInServer) (error) {
	return errors.New("UnImplemented")
}

func (s *TestServer) ReceiveChangeEvents(*empty.Empty, voltha.VolthaService_ReceiveChangeEventsServer) error {
	return errors.New("UnImplemented")
}
func (s *TestServer) GetDeviceGroup(context.Context, *common.ID) (*voltha.DeviceGroup, error) {
	return nil, errors.New("UnImplemented")
}
func (s *TestServer) CreateAlarmFilter(context.Context, *voltha.AlarmFilter) (*voltha.AlarmFilter, error) {
	return nil, errors.New("UnImplemented")
}
func (s *TestServer) GetAlarmFilter(context.Context, *common.ID) (*voltha.AlarmFilter, error) {
	return nil, errors.New("UnImplemented")
}
func (s *TestServer) UpdateAlarmFilter(context.Context, *voltha.AlarmFilter) (*voltha.AlarmFilter, error) {
	return nil, errors.New("UnImplemented")
}
func (s *TestServer) DeleteAlarmFilter(context.Context, *common.ID) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}
func (s *TestServer) ListAlarmFilters(context.Context, *empty.Empty) (*voltha.AlarmFilters, error) {
	return nil, errors.New("UnImplemented")
}
func (s *TestServer) GetImages(context.Context, *common.ID) (*voltha.Images, error) {
	return nil, errors.New("UnImplemented")
}
func (s *TestServer) SelfTest(context.Context, *common.ID) (*voltha.SelfTestResponse, error) {
	return nil, errors.New("UnImplemented")
}
func (s *TestServer) GetMibDeviceData(context.Context, *common.ID) (*omci.MibDeviceData, error) {
	return nil, errors.New("UnImplemented")
}
func (s *TestServer) GetAlarmDeviceData(context.Context, *common.ID) (*omci.AlarmDeviceData, error) {
	return nil, errors.New("UnImplemented")
}
func (s *TestServer) SimulateAlarm(context.Context, *voltha.SimulateAlarmRequest) (*common.OperationResp, error) {
	return nil, errors.New("UnImplemented")
}
func (s *TestServer) Subscribe(context.Context, *voltha.OfAgentSubscriber) (*voltha.OfAgentSubscriber, error) {
	return nil, errors.New("UnImplemented")
}

func GetVolthaServer() *simulation.GrpcServer {
	return simulation.NewGrpcServer(
		func(server *grpc.Server) {
			voltha.RegisterVolthaServiceServer(server, &TestServer{})
		},
		GRPC_ADDRESS,
	)
}