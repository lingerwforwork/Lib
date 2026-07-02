package device

import (
	"ascend-common/devmanager/common"
	"ascend-common/devmanager/dcmi"
	"fmt"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

//manager := &dcmi.DcManager{}
//	manager.DcInit()
//	//获取card id
//	manager.DcGetCardList()
//	//获取设备id
//	manager.DcGetDeviceNumInCard()
//	manager.DcGetMemoryInfo()

// 华为显卡存在一个cardId 多个chipId(deviceId) 共享一张卡的显存
//type AscendNpuClient interface {
//	GetCount() (int, error)
//	//GetCardList() ([]int32, error)
//	GetMemoryInfo(cardId int32) (*nvml.Memory, error)
//	//GetDeviceCount(cardId int32) (int32, error)
//	GetTotalMemorySize() (uint64, error)
//	IsAllocatable(memoryMB uint64) error
//	Init() error
//	Shutdown() error
//}

type ascendNpuClient struct {
	manager *dcmi.DcManager
}

func NewAscendNpuClient() GpuClient {
	return &ascendNpuClient{
		manager: &dcmi.DcManager{},
	}
}

// 获取card数量
func (d *ascendNpuClient) GetCount() (int, error) {
	//count, err := d.manager.DcGetAllDeviceCount()
	total, _, err := d.manager.DcGetCardList()
	if err != nil {
		return common.RetError, fmt.Errorf("get device count failed, error: %v", err)
	}
	return int(total), err
}

func (d *ascendNpuClient) GetMemoryInfo(cardId int) (*nvml.Memory, error) {
	info, err := dcmi.FuncDcmiGetDeviceHbmInfo(int32(cardId), 0)
	if err != nil {
		return nil, fmt.Errorf("get device memory failed, error: %v", err)
	}
	//注意单位 nvidia 是字节
	return &nvml.Memory{
		Total: info.MemorySize * 1024 * 1024,
		Free:  (info.MemorySize - info.Usage) * 1024 * 1024,
		Used:  info.Usage * 1024 * 1024,
	}, nil
}

func (d *ascendNpuClient) GetTotalMemorySize() (uint64, error) {
	_, cardIds, err := d.manager.DcGetCardList()
	if err != nil {
		return 0, err
	}
	var mb uint64 = 0
	for _, cardId := range cardIds {
		//deviceNum, err := d.manager.DcGetDeviceNumInCard(cardId)
		//if err != nil {
		//	return 0, err
		//}
		//var i int32
		//for i = 0; i < deviceNum; i++ {
		info, err := dcmi.FuncDcmiGetDeviceHbmInfo(cardId, 0)
		if err != nil {
			return 0, err
		}
		mb += info.MemorySize
		//}
	}
	return mb, nil
}

func (d *ascendNpuClient) IsAllocatable(memoryMB uint64) error {
	_, cardIds, err := d.manager.DcGetCardList()
	if err != nil {
		return err
	}
	var freeMB uint64 = 0
	for _, cardId := range cardIds {
		//deviceNum, err := d.manager.DcGetDeviceNumInCard(cardId)
		//if err != nil {
		//	return err
		//}
		//var i int32
		//for i = 0; i < deviceNum; i++ {
		info, err := dcmi.FuncDcmiGetDeviceHbmInfo(cardId, 0)
		if err != nil {
			return err
		}
		freeMB = freeMB + (info.MemorySize - info.Usage)
		//}
	}
	if freeMB < uint64(float64(memoryMB)) {
		return fmt.Errorf("unable to allocate memory. free: %dMB, requested: %dMB", freeMB, memoryMB)
	}
	return nil
}

func (d *ascendNpuClient) Init() error {
	return d.manager.DcInit()
}

func (d *ascendNpuClient) Shutdown() error {
	return d.manager.DcShutDown()
}

func (d *ascendNpuClient) GetCardList() ([]int32, error) {
	_, cardIds, err := d.manager.DcGetCardList()
	return cardIds, err
}

func (d *ascendNpuClient) GetDeviceCount(cardId int32) (int32, error) {
	return d.manager.DcGetDeviceNumInCard(cardId)
}
