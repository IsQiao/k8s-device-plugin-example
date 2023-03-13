package main

import (
	"context"
	"fmt"

	"math/rand"
	"time"

	"github.com/golang/glog"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

var _ pluginapi.DevicePluginServer = (*Plugin)(nil)

type Plugin struct {
	Devices map[string]Device
}

// Allocate is called during container creation so that the Device
// Plugin can run device specific operations and instruct Kubelet
// of the steps to make the Device available in the container
func (p *Plugin) Allocate(ctx context.Context, r *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
	// allocate active used devices.
	var response pluginapi.AllocateResponse
	var car pluginapi.ContainerAllocateResponse
	var dev *pluginapi.DeviceSpec

	for _, req := range r.ContainerRequests {
		car = pluginapi.ContainerAllocateResponse{}

		for _, id := range req.DevicesIDs {
			// custom your allocate gpu resource logic.
			glog.Infof("Allocating device ID: %s", id)

			demoDev := p.Devices[id]
			dev = new(pluginapi.DeviceSpec)
			devpath := fmt.Sprintf("/dev/dri/%s", demoDev.Id)
			dev.HostPath = devpath
			dev.ContainerPath = devpath
			dev.Permissions = "rw"
			car.Devices = append(car.Devices, dev)
		}

		car.Devices = append(car.Devices, dev)
		response.ContainerResponses = append(response.ContainerResponses, &car)
	}

	return &response, nil
}

// GetDevicePluginOptions returns options to be communicated with Device
func (*Plugin) GetDevicePluginOptions(context.Context, *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
	return &pluginapi.DevicePluginOptions{}, nil
}

func (*Plugin) GetPreferredAllocation(context.Context, *pluginapi.PreferredAllocationRequest) (*pluginapi.PreferredAllocationResponse, error) {
	return &pluginapi.PreferredAllocationResponse{}, nil
}

// ListAndWatch returns a stream of List of Devices
// Whenever a Device state change or a Device disappears, ListAndWatch
// returns the new list
func (p *Plugin) ListAndWatch(e *pluginapi.Empty, s pluginapi.DevicePlugin_ListAndWatchServer) error {
	p.Devices = GetDevices()

	// get all kube devices
	kubeDevices := []*pluginapi.Device{}
	{
		for id := range p.Devices {
			dev := &pluginapi.Device{
				ID:     id,
				Health: pluginapi.Healthy,
			}
			kubeDevices = append(kubeDevices, dev)
		}
	}

	// send all devices to kubelet
	s.Send(&pluginapi.ListAndWatchResponse{Devices: kubeDevices})
	glog.Info("Send all devices info")

	// watch device health status and report
	{
		// mock device unhealth
		func() {
			ticker := time.NewTicker(10 * time.Second)
			for {
				<-ticker.C
				devIndex := rand.Intn(len(kubeDevices))
				isHealth := rand.Intn(2) == 1
				currentHealth := pluginapi.Healthy
				if !isHealth {
					currentHealth = pluginapi.Unhealthy
				}
				kubeDevices[devIndex].Health = currentHealth

				// report
				s.Send(&pluginapi.ListAndWatchResponse{Devices: kubeDevices})
				glog.Infof("Health state changed, device id: %s, status: %s", kubeDevices[devIndex].ID, kubeDevices[devIndex].Health)
			}
		}()

		select {}
	}
}

func (p *Plugin) PreStartContainer(ctx context.Context, r *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
	return &pluginapi.PreStartContainerResponse{}, nil
}
