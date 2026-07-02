package device

import (
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

type nvidiaGpuClient struct {
	devices map[int]nvml.Device
}

func NewNvidiaGpuClient() GpuClient {
	return &nvidiaGpuClient{
		devices: make(map[int]nvml.Device, 0),
	}
}

func (c *nvidiaGpuClient) getDevice(gpuId int) (nvml.Device, error) {
	device, exist := c.devices[gpuId]
	if exist {
		return device, nil
	}
	newDevice, ret := nvml.DeviceGetHandleByIndex(gpuId)
	if ret != nvml.SUCCESS {
		return nil, fmt.Errorf("unable to get device at index %d: %v", gpuId, nvml.ErrorString(ret))
	}
	c.devices[gpuId] = newDevice
	return newDevice, nil
}

func (c *nvidiaGpuClient) GetCount() (int, error) {
	count, ret := nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		return -1, fmt.Errorf("unable to get device count: %v", nvml.ErrorString(ret))
	}
	return count, nil
}

func (c *nvidiaGpuClient) GetMemoryInfo(gpuId int) (*nvml.Memory, error) {
	device, err := c.getDevice(gpuId)
	if err != nil {
		return nil, err
	}
	memory, ret := device.GetMemoryInfo()
	if ret != nvml.SUCCESS {
		return nil, fmt.Errorf("unable to shutdown nvml: %v", nvml.ErrorString(ret))
	}
	return &memory, nil
}

func (c *nvidiaGpuClient) GetTotalMemorySize() (uint64, error) {
	deviceCount, err := c.GetCount()
	if err != nil {
		return 0, err
	}
	var mb uint64 = 0
	for i := range deviceCount {
		memory, err := c.GetMemoryInfo(i)
		if err != nil {
			return 0, err
		}
		mb += memory.Total / 1024 / 1024
	}
	return mb, nil
}

func (c *nvidiaGpuClient) IsAllocatable(memoryMB uint64) error {
	deviceCount, err := c.GetCount()
	if err != nil {
		return err
	}
	var freeMB uint64 = 0
	for id := range deviceCount {
		memory, err := c.GetMemoryInfo(id)
		if err != nil {
			return err
		}
		freeMB += memory.Free / 1024 / 1024
	}
	if freeMB < uint64(float64(memoryMB)) {
		return fmt.Errorf("unable to allocate memory. free: %dMB, requested: %dMB", freeMB, memoryMB)
	}
	return nil
}

func (c *nvidiaGpuClient) Init() error {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		return fmt.Errorf("unable to initialize nvml: %v", nvml.ErrorString(ret))
	}
	return nil
}

func (c *nvidiaGpuClient) Shutdown() error {
	ret := nvml.Shutdown()
	if ret != nvml.SUCCESS {
		return fmt.Errorf("unable to shutdown nvml: %v", nvml.ErrorString(ret))
	}
	return nil
}
