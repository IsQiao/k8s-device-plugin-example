package main

import (
	"github.com/kubevirt/device-plugin-manager/pkg/dpm"
)

var _ dpm.ListerInterface = (*Lister)(nil)

// Lister serves as an interface between imlementation and Manager machinery. User passes
// implementation of this interface to NewManager function. Manager will use it to obtain resource
// namespace, monitor available resources and instantate a new plugin for them.
type Lister struct {
	ResUpdateChan chan dpm.PluginNameList
}

// Discover implements dpm.ListerInterface
func (l *Lister) Discover(pluginNameListCh chan dpm.PluginNameList) {
	for {
		select {
		case newResourcesList := <-l.ResUpdateChan: // New resources found
			pluginNameListCh <- newResourcesList
		case <-pluginNameListCh: // Stop message received
			// Stop resourceUpdateCh
			return
		}
	}
}

// GetResourceNamespace implements dpm.ListerInterface
func (*Lister) GetResourceNamespace() string {
	return "demo.com"
}

// NewPlugin instantiates a plugin implementation. It is given the last name of the resource,
// e.g. for resource name "color.example.com/red" that would be "red". It must return valid
// implementation of a PluginInterface.
func (*Lister) NewPlugin(string) dpm.PluginInterface {
	return &Plugin{}
}
