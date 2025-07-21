package sysinfo

import (
	"github.com/dullkingsman/go-pkg/roga/internal/json"
	"github.com/dullkingsman/go-pkg/utils"
	"log"
	"os"
	"runtime"
)

type SystemSpecifications struct {
	Product    ProductIdentifier `json:"product"`
	InstanceId *string           `json:"instanceId,omitempty"`
	MachineId  *string           `json:"machineId,omitempty"`
	MacAddress *string           `json:"macAddress,omitempty"`
	Os         string            `json:"os"`
	Arch       string            `json:"arch"`
	CpuCores   int               `json:"cpuCores"`
	Memory     uint64            `json:"memory"`
	SwapSize   uint64            `json:"swapSize"`
	DiskSize   uint64            `json:"diskSize"`
	PageSize   int               `json:"pageSize"`
}

func GetSystemSpecifications() SystemSpecifications {
	var (
		product    = getProductIdentifier()
		provider   = getCloudProvider(product)
		instanceId = getCloudInstanceID(provider)
		machineId  = getMachineID()
		macAddress = getMacAddress()
	)

	var totalMemory, _, err = GetMemoryStats()

	if err != nil {
		log.Fatalf("failed to read memory stats: " + err.Error())
	}

	swapSize, _, err := readSwapStats()

	if err != nil {
		log.Fatalf("failed to read swap stats: " + err.Error())
	}

	diskSize, _, err := GetDiskStats("/")

	if err != nil {
		log.Fatalf("failed to read memory stats: " + err.Error())
	}

	return SystemSpecifications{
		Product:    product,
		InstanceId: instanceId,
		MachineId:  machineId,
		MacAddress: macAddress,
		Os:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		CpuCores:   runtime.NumCPU(),
		Memory:     totalMemory,
		SwapSize:   swapSize,
		DiskSize:   diskSize,
		PageSize:   os.Getpagesize(),
	}
}

func (c *SystemSpecifications) Json() ([]byte, error) {
	if c == nil {
		return nil, nil
	}

	var mj = json.NewManualJson()

	product, err := c.Product.Json()

	if err != nil {
		return nil, err
	}

	mj.WriteMarshalledJsonField("product", product)
	mj.WriteStringField("instanceId", utils.ValueOr(c.InstanceId, ""), true)
	mj.WriteStringField("machineId", utils.ValueOr(c.MachineId, ""), true)
	mj.WriteStringField("macAddress", utils.ValueOr(c.MacAddress, ""), true)
	mj.WriteStringField("os", c.Os)
	mj.WriteStringField("arch", c.Arch)
	mj.WriteInt64Field("cpuCores", int64(c.CpuCores))
	mj.WriteUint64Field("memory", c.Memory)
	mj.WriteUint64Field("swapSize", c.SwapSize)
	mj.WriteUint64Field("diskSize", c.DiskSize)
	mj.WriteInt64Field("pageSize", int64(c.PageSize))

	return mj.End(), nil
}
