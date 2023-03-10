package main

import (
	"time"

	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
)

func main() {
	l := Lister{
		ResUpdateChan: make(chan dpm.PluginNameList),
	}
	manager := dpm.NewManager(&l)

	// device ready
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		<-ticker.C
		l.ResUpdateChan <- []string{"demo_device"}
	}()
	manager.Run()
}
