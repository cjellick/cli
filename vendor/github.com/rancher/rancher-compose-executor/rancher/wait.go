package rancher

import (
	"time"

	"github.com/rancher/go-rancher/v2"
)

func (r *RancherService) WaitFor(resource *client.Resource, output interface{}, transitioning func() string) error {
	for {
		if transitioning() != "yes" {
			return nil
		}

		time.Sleep(150 * time.Millisecond)

		err := r.context.Client.Reload(resource, output)
		if err != nil {
			return err
		}
	}
}

func (r *RancherService) Wait(service *client.Service) error {
	return r.WaitFor(&service.Resource, service, func() string {
		return service.Transitioning
	})
}

func (r *RancherService) WaitState(service *client.Service) error {
	service, err := r.context.Client.Service.ById(service.Id)

	if service.HealthState != r.context.WaitState {
		if err = r.WaitState(service); err != nil {
			return err
		}
	}
	return nil
}

func (r *RancherService) waitInstance(instance *client.Instance) error {
	return r.WaitFor(&instance.Resource, instance, func() string {
		return instance.Transitioning
	})
}
