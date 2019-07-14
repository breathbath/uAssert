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

var simCasesMap simulation.SimulationCasesMap

func init() {
	simCasesMap = simulation.GetSimulationMap(GetStubs())
}

type VolthaServerSimulator struct {
}

func (s *VolthaServerSimulator) UpdateLogLevel(context.Context, *voltha.Logging) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Get the membership group of a Voltha Core
func (s *VolthaServerSimulator) GetMembership(ctx context.Context, in *empty.Empty) (*voltha.Membership, error) {
	return nil, errors.New("UnImplemented")
}

// Set the membership group of a Voltha Core
func (s *VolthaServerSimulator) UpdateMembership(ctx context.Context, in *voltha.Membership) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Get high level information on the Voltha cluster
func (s *VolthaServerSimulator) GetVoltha(ctx context.Context, in *empty.Empty) (*voltha.Voltha, error) {
	return nil, errors.New("UnImplemented")
}

// List all Voltha cluster core instances
func (s *VolthaServerSimulator) ListCoreInstances(ctx context.Context, in *empty.Empty) (*voltha.CoreInstances, error) {
	return nil, errors.New("UnImplemented")
}

// Get details on a Voltha cluster instance
func (s *VolthaServerSimulator) GetCoreInstance(ctx context.Context, in *common.ID) (*voltha.CoreInstance, error) {
	return nil, errors.New("UnImplemented")
}

// List all active adapters (plugins) in the Voltha cluster
func (s *VolthaServerSimulator) ListAdapters(ctx context.Context, in *empty.Empty) (*voltha.Adapters, error) {
	return nil, errors.New("UnImplemented")
}

// List all logical devices managed by the Voltha cluster
func (s *VolthaServerSimulator) ListLogicalDevices(ctx context.Context, in *empty.Empty) (*voltha.LogicalDevices, error) {
	return nil, errors.New("UnImplemented")
}

// Get additional information on a given logical device
func (s *VolthaServerSimulator) GetLogicalDevice(ctx context.Context, in *common.ID) (*voltha.LogicalDevice, error) {
	return nil, errors.New("UnImplemented")
}

// List ports of a logical device
func (s *VolthaServerSimulator) ListLogicalDevicePorts(ctx context.Context, in *common.ID) (*voltha.LogicalPorts, error) {
	return nil, errors.New("UnImplemented")
}

// Gets a logical device port
func (s *VolthaServerSimulator) GetLogicalDevicePort(ctx context.Context, in *voltha.LogicalPortId) (*voltha.LogicalPort, error) {
	return nil, errors.New("UnImplemented")
}

// Enables a logical device port
func (s *VolthaServerSimulator) EnableLogicalDevicePort(ctx context.Context, in *voltha.LogicalPortId) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Disables a logical device port
func (s *VolthaServerSimulator) DisableLogicalDevicePort(ctx context.Context, in *voltha.LogicalPortId) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// List all flows of a logical device
func (s *VolthaServerSimulator) ListLogicalDeviceFlows(ctx context.Context, in *common.ID) (*openflow_13.Flows, error) {
	return nil, errors.New("UnImplemented")
}

// Update flow table for logical device
func (s *VolthaServerSimulator) UpdateLogicalDeviceFlowTable(ctx context.Context, in *openflow_13.FlowTableUpdate) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Update meter table for logical device
func (s *VolthaServerSimulator) UpdateLogicalDeviceMeterTable(ctx context.Context, in *openflow_13.MeterModUpdate) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Get all meter stats for logical device
func (s *VolthaServerSimulator) GetMeterStatsOfLogicalDevice(ctx context.Context, in *common.ID) (*openflow_13.MeterStatsReply, error) {
	return nil, errors.New("UnImplemented")
}

// List all flow groups of a logical device
func (s *VolthaServerSimulator) ListLogicalDeviceFlowGroups(ctx context.Context, in *common.ID) (*openflow_13.FlowGroups, error) {
	return nil, errors.New("UnImplemented")
}

// Update group table for device
func (s *VolthaServerSimulator) UpdateLogicalDeviceFlowGroupTable(ctx context.Context, in *openflow_13.FlowGroupTableUpdate) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// List all physical devices controlled by the Voltha cluster
func (s *VolthaServerSimulator) ListDevices(ctx context.Context, in *empty.Empty) (*voltha.Devices, error) {
	c, found := simulation.FindSimulatedCaseForRequest(in, "/voltha.VolthaService/ListDevices", simCasesMap)
	if !found {
		return nil, fmt.Errorf("Not found")
	}

	return c.Response.(*voltha.Devices), nil
}

// List all physical devices IDs controlled by the Voltha cluster
func (s *VolthaServerSimulator) ListDeviceIds(ctx context.Context, in *empty.Empty) (*common.IDs, error) {
	return nil, errors.New("UnImplemented")
}

// Request to a voltha Core to reconcile a set of devices based on their IDs
func (s *VolthaServerSimulator) ReconcileDevices(ctx context.Context, in *common.IDs) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Get more information on a given physical device
func (s *VolthaServerSimulator) GetDevice(ctx context.Context, in *common.ID) (*voltha.Device, error) {
	return nil, errors.New("UnImplemented")
}

// Pre-provision a new physical device
func (s *VolthaServerSimulator) CreateDevice(ctx context.Context, in *voltha.Device) (*voltha.Device, error) {
	return nil, errors.New("UnImplemented")
}

// Enable a device.  If the device was in pre-provisioned state then it
// will transition to ENABLED state.  If it was is DISABLED state then it
// will transition to ENABLED state as well.
func (s *VolthaServerSimulator) EnableDevice(ctx context.Context, in *common.ID) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Disable a device
func (s *VolthaServerSimulator) DisableDevice(ctx context.Context, in *common.ID) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Reboot a device
func (s *VolthaServerSimulator) RebootDevice(ctx context.Context, in *common.ID) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Delete a device
func (s *VolthaServerSimulator) DeleteDevice(ctx context.Context, in *common.ID) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// Request an image download to the standby partition
// of a device.
// Note that the call is expected to be non-blocking.
func (s *VolthaServerSimulator) DownloadImage(ctx context.Context, in *voltha.ImageDownload) (*common.OperationResp, error) {
	return nil, errors.New("UnImplemented")
}

// Get image download status on a device
// The request retrieves progress on device and updates db record
func (s *VolthaServerSimulator) GetImageDownloadStatus(ctx context.Context, in *voltha.ImageDownload) (*voltha.ImageDownload, error) {
	return nil, errors.New("UnImplemented")
}

// Get image download db record
func (s *VolthaServerSimulator) GetImageDownload(ctx context.Context, in *voltha.ImageDownload) (*voltha.ImageDownload, error) {
	return nil, errors.New("UnImplemented")
}

// List image download db records for a given device
func (s *VolthaServerSimulator) ListImageDownloads(ctx context.Context, in *common.ID) (*voltha.ImageDownloads, error) {
	return nil, errors.New("UnImplemented")
}

// Cancel an existing image download process on a device
func (s *VolthaServerSimulator) CancelImageDownload(ctx context.Context, in *voltha.ImageDownload) (*common.OperationResp, error) {
	return nil, errors.New("UnImplemented")
}

// Activate the specified image at a standby partition
// to active partition.
// Depending on the device implementation, this call
// may or may not cause device reboot.
// If no reboot, then a reboot is required to make the
// activated image running on device
// Note that the call is expected to be non-blocking.
func (s *VolthaServerSimulator) ActivateImageUpdate(ctx context.Context, in *voltha.ImageDownload) (*common.OperationResp, error) {
	return nil, errors.New("UnImplemented")
}

// Revert the specified image at standby partition
// to active partition, and revert to previous image
// Depending on the device implementation, this call
// may or may not cause device reboot.
// If no reboot, then a reboot is required to make the
// previous image running on device
// Note that the call is expected to be non-blocking.
func (s *VolthaServerSimulator) RevertImageUpdate(ctx context.Context, in *voltha.ImageDownload) (*common.OperationResp, error) {
	return nil, errors.New("UnImplemented")
}

// List ports of a device
func (s *VolthaServerSimulator) ListDevicePorts(ctx context.Context, in *common.ID) (*voltha.Ports, error) {
	return nil, errors.New("UnImplemented")
}

// List pm config of a device
func (s *VolthaServerSimulator) ListDevicePmConfigs(ctx context.Context, in *common.ID) (*voltha.PmConfigs, error) {
	return nil, errors.New("UnImplemented")
}

// Update the pm config of a device
func (s *VolthaServerSimulator) UpdateDevicePmConfigs(ctx context.Context, in *voltha.PmConfigs) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}

// List all flows of a device
func (s *VolthaServerSimulator) ListDeviceFlows(ctx context.Context, in *common.ID) (*openflow_13.Flows, error) {
	return nil, errors.New("UnImplemented")
}

// List all flow groups of a device
func (s *VolthaServerSimulator) ListDeviceFlowGroups(ctx context.Context, in *common.ID) (*openflow_13.FlowGroups, error) {
	return nil, errors.New("UnImplemented")
}

// List device types known to Voltha
func (s *VolthaServerSimulator) ListDeviceTypes(ctx context.Context, in *empty.Empty) (*voltha.DeviceTypes, error) {
	return nil, errors.New("UnImplemented")
}

// Get additional information on a device type
func (s *VolthaServerSimulator) GetDeviceType(ctx context.Context, in *common.ID) (*voltha.DeviceType, error) {
	return nil, errors.New("UnImplemented")
}

// List all device sharding groups
func (s *VolthaServerSimulator) ListDeviceGroups(ctx context.Context, in *empty.Empty) (*voltha.DeviceGroups, error) {
	return nil, errors.New("UnImplemented")
}

// Stream control packets to the dataplane
func (s *VolthaServerSimulator) StreamPacketsOut(voltha.VolthaService_StreamPacketsOutServer) error {
	return errors.New("UnImplemented")
}

// Receive control packet stream
func (s *VolthaServerSimulator) ReceivePacketsIn(*empty.Empty, voltha.VolthaService_ReceivePacketsInServer) (error) {
	return errors.New("UnImplemented")
}

func (s *VolthaServerSimulator) ReceiveChangeEvents(*empty.Empty, voltha.VolthaService_ReceiveChangeEventsServer) error {
	return errors.New("UnImplemented")
}
func (s *VolthaServerSimulator) GetDeviceGroup(context.Context, *common.ID) (*voltha.DeviceGroup, error) {
	return nil, errors.New("UnImplemented")
}
func (s *VolthaServerSimulator) CreateAlarmFilter(context.Context, *voltha.AlarmFilter) (*voltha.AlarmFilter, error) {
	return nil, errors.New("UnImplemented")
}
func (s *VolthaServerSimulator) GetAlarmFilter(context.Context, *common.ID) (*voltha.AlarmFilter, error) {
	return nil, errors.New("UnImplemented")
}
func (s *VolthaServerSimulator) UpdateAlarmFilter(context.Context, *voltha.AlarmFilter) (*voltha.AlarmFilter, error) {
	return nil, errors.New("UnImplemented")
}
func (s *VolthaServerSimulator) DeleteAlarmFilter(context.Context, *common.ID) (*empty.Empty, error) {
	return nil, errors.New("UnImplemented")
}
func (s *VolthaServerSimulator) ListAlarmFilters(context.Context, *empty.Empty) (*voltha.AlarmFilters, error) {
	return nil, errors.New("UnImplemented")
}
func (s *VolthaServerSimulator) GetImages(context.Context, *common.ID) (*voltha.Images, error) {
	return nil, errors.New("UnImplemented")
}
func (s *VolthaServerSimulator) SelfTest(context.Context, *common.ID) (*voltha.SelfTestResponse, error) {
	return nil, errors.New("UnImplemented")
}
func (s *VolthaServerSimulator) GetMibDeviceData(context.Context, *common.ID) (*omci.MibDeviceData, error) {
	return nil, errors.New("UnImplemented")
}
func (s *VolthaServerSimulator) GetAlarmDeviceData(context.Context, *common.ID) (*omci.AlarmDeviceData, error) {
	return nil, errors.New("UnImplemented")
}
func (s *VolthaServerSimulator) SimulateAlarm(context.Context, *voltha.SimulateAlarmRequest) (*common.OperationResp, error) {
	return nil, errors.New("UnImplemented")
}
func (s *VolthaServerSimulator) Subscribe(context.Context, *voltha.OfAgentSubscriber) (*voltha.OfAgentSubscriber, error) {
	return nil, errors.New("UnImplemented")
}

func NewVolthaServerSimulator(address string) *simulation.GrpcServer {
	return simulation.NewGrpcServer(
		address,
		func(server *grpc.Server) {
			voltha.RegisterVolthaServiceServer(server, &VolthaServerSimulator{})
		},
	)
}