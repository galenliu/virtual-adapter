package virtual_adapter

import (
	"addon/devices"
	"time"

	"addon"
	"context"
	"fmt"
)

type VirtualAdapter struct {
	*addon.AdapterProxy
}

func NewVirtualAdapter() *VirtualAdapter {
	adapter := &VirtualAdapter{
		AdapterProxy: addon.NewAdapterProxy("Virtual-adapter", "Virtual-adapter", "yeelight-adapter"),
	}
	adapter.OnPairing = adapter.StartPairing
	return adapter
}

func (adapter *VirtualAdapter) StartPairing(timeout float64) {
	fmt.Printf("adapter(%s)start pairing", adapter.ID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	var pairing = func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				cancel()
				fmt.Printf("adapter(%s) pair over -------------- \t\n", adapter.ID)
				return
			default:
				adapter.Discover()
			}

		}
	}
	go pairing(ctx)
}

var id int = 0

func (adapter *VirtualAdapter) Discover() {
	id++
	time.Sleep(time.Duration(1) * time.Second)
	var virtualLight = devices.NewLightBulb(fmt.Sprintf("virtual-light%d", id), fmt.Sprintf("virtual-light%d", id))
	adapter.AddDevice(virtualLight.Device)

}
