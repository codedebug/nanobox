package code

import (
	"net"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ipControl"
)

// processCodeDestroy ...
type processCodeDestroy struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("code_destroy", codeDestroyFunc)
}

//
func codeDestroyFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	if control.Meta["name"] == "" {
		return nil, errMissingImageOrName
	}
	return &processCodeDestroy{control: control}, nil
}

//
func (codeDestroy processCodeDestroy) Results() processor.ProcessControl {
	return codeDestroy.control
}

//
func (codeDestroy *processCodeDestroy) Process() error {

	// get the service from the database
	service := models.Service{}

	//
	if err := data.Get(util.AppName(), codeDestroy.control.Meta["name"], &service); err != nil {
		return err
	}

	//
	if err := docker.ContainerRemove(service.ID); err != nil {
		return err
	}

	//
	if err := provider.RemoveNat(service.ExternalIP, service.InternalIP); err != nil {
		return err
	}

	//
	if err := provider.RemoveIP(service.ExternalIP); err != nil {
		return err
	}

	//
	if err := ipControl.ReturnIP(net.ParseIP(service.ExternalIP)); err != nil {
		return err
	}

	//
	if err := ipControl.ReturnIP(net.ParseIP(service.InternalIP)); err != nil {
		return err
	}

	// save the service
	return data.Delete(util.AppName(), codeDestroy.control.Meta["name"])
}
