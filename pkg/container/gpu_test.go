package container

import (
	"testing"

	libContainer "github.com/moby/moby/api/types/container"
)

func TestApplyGPUDeviceRequests(t *testing.T) {
	t.Run("empty gpu ids", func(t *testing.T) {
		hostConfig := &libContainer.HostConfig{}
		applyGPUDeviceRequests(hostConfig, nil)
		if len(hostConfig.Resources.DeviceRequests) != 0 {
			t.Fatalf("expected no device requests, got %d", len(hostConfig.Resources.DeviceRequests))
		}
	})

	t.Run("specific gpu ids", func(t *testing.T) {
		hostConfig := &libContainer.HostConfig{}
		applyGPUDeviceRequests(hostConfig, []uint32{0, 2})
		reqs := hostConfig.Resources.DeviceRequests
		if len(reqs) != 1 {
			t.Fatalf("expected 1 device request, got %d", len(reqs))
		}
		req := reqs[0]
		if req.Driver != "nvidia" || req.Count != 0 {
			t.Fatalf("unexpected request: %+v", req)
		}
		if len(req.DeviceIDs) != 2 || req.DeviceIDs[0] != "0" || req.DeviceIDs[1] != "2" {
			t.Fatalf("unexpected device ids: %v", req.DeviceIDs)
		}
		if len(req.Capabilities) != 1 || len(req.Capabilities[0]) != 1 || req.Capabilities[0][0] != "gpu" {
			t.Fatalf("unexpected capabilities: %v", req.Capabilities)
		}
	})
}
