package virtual_adapter

import (
	"context"
	"fmt"
	addon "gateway_addon_golang"
	"time"

)

type VirtualAdapter struct {
	*addon.AdapterProxy
	devices map[string]*addon.AdapterProxy
}



func NewVirtualAdapter() *YeeAdapter {
	adapter := &YeeAdapter{
		AdapterProxy: addon.NewAdapterProxy("virtual-things-adapter", "virtual-things-adapter"),
		devices:      make(map[string]*addon.Device),
	}
	adapter.OnStartPairing = adapter.StartPairing
	return adapter
}

func (adapter *VirtualAdapter) StartPairing(timeout int) {
	fmt.Println("start pairing")

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	var pairing = func(ctx context.Context) {
		select {
		case <-ctx.Done():
			cancelFunc()
			return
		default:
			adapter.Discover()
		}
	}
	go pairing(ctx)
}

func (adapter *VirtualAdapter) Discover() {



		adapter.HandleDeviceAdded()

}
