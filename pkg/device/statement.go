package device

import "github.com/NVIDIA/go-nvml/pkg/nvml"

type GpuClient interface {
	GetCount() (int, error)
	GetMemoryInfo(gpuId int) (*nvml.Memory, error)
	GetTotalMemorySize() (uint64, error)
	IsAllocatable(memoryMB uint64) error
	Init() error
	Shutdown() error
}
