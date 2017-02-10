package broker

import (
	"code.cloudfoundry.org/lager"
	"context"
	"errors"
	"github.com/Altoros/cf-broker-boilerplate/cmd"
	"github.com/Altoros/cf-broker-boilerplate/config"
	"github.com/Altoros/cf-broker-boilerplate/model"
	"github.com/jinzhu/gorm"
	"github.com/pivotal-cf/brokerapi"
)

type serviceBroker struct {
	Config config.Config
	Db     *gorm.DB
	Opts   cmd.CommandOpts
	Logger lager.Logger
}

func NewServiceBroker(
	opts cmd.CommandOpts,
	config config.Config,
	db *gorm.DB,
	logger lager.Logger) *serviceBroker {

	return &serviceBroker{
		Config: config,
		Opts:   opts,
		Db:     db,
		Logger: logger,
	}
}

func (b *serviceBroker) Services(context context.Context) []brokerapi.Service {
	planList := []brokerapi.ServicePlan{}

	for _, plan := range b.Config.Plans {
		planList = append(planList, brokerapi.ServicePlan{
			ID:          plan.Name,
			Name:        plan.Name,
			Description: plan.Description,
		})
	}
	b.Logger.Info("Serving a catalog request")
	return []brokerapi.Service{
		{
			ID:            b.Opts.ServiceID,
			Name:          b.Opts.Name,
			Description:   b.Opts.Description,
			Bindable:      true,
			Tags:          []string{"boilerplate", "test"},
			Plans:         planList,
			PlanUpdatable: true,
		},
	}
}

func (b *serviceBroker) Provision(context context.Context, instanceId string, provisionDetails brokerapi.ProvisionDetails, asyncAllowed bool) (brokerapi.ProvisionedServiceSpec, error) {
	planName := provisionDetails.PlanID
	plan, err := b.Config.PlanByName(planName)
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{IsAsync: false}, err
	}
	b.Logger.Info("Starting provisioning a service instance", lager.Data{
		"instance-id":       instanceId,
		"plan-id":           plan.Name,
		"plan-desctiption":  plan.Description,
		"organization-guid": provisionDetails.OrganizationGUID,
		"space-guid":        provisionDetails.SpaceGUID,
	})

	// create a service instance here

	serivceInstance := model.ServiceInstance{
		InstanceId: instanceId,
	}

	err = b.Db.Create(&serivceInstance).Error
	if err != nil {
		return brokerapi.ProvisionedServiceSpec{IsAsync: false}, err
	}

	b.Logger.Info("service instance created", lager.Data{
		"instance-id": instanceId,
	})

	return brokerapi.ProvisionedServiceSpec{IsAsync: false}, err
}

func (b *serviceBroker) Update(context context.Context, instanceId string, updateDetails brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{IsAsync: false}, errors.New("Not implemented")
}

func (b *serviceBroker) Deprovision(context context.Context, instanceId string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	var serivceInstance model.ServiceInstance
	err := b.Db.First(&serivceInstance, "instance_id = ?", instanceId).Error

	// delete a service instance here

	err = b.Db.Delete(&serivceInstance).Error
	b.Logger.Info("service instance is removed", lager.Data{
		"instance-id": instanceId,
	})
	return brokerapi.DeprovisionServiceSpec{IsAsync: false}, err
}

func (b *serviceBroker) Bind(context context.Context, instanceId, bindingId string, details brokerapi.BindDetails) (brokerapi.Binding, error) {
	var serivceInstance model.ServiceInstance
	b.Logger.Info("Searching for a service instance", lager.Data{"instance_id": instanceId})
	err := b.Db.First(&serivceInstance, "instance_id = ?", instanceId).Error

	b.Logger.Info("Starting binding a service instance", lager.Data{
		"instance-id":                instanceId,
		"serivceInstance.instanceId": serivceInstance.InstanceId,
		"app-guid":                   details.AppGUID,
		"plan-id":                    details.PlanID,
	})

	if err != nil {
		return brokerapi.Binding{}, err
	}

	// do your service binding here

	serviceBinding := model.ServiceBinding{
		BindingId:  bindingId,
		InstanceId: instanceId,
	}
	err = b.Db.Create(&serviceBinding).Error
	if err != nil {
		return brokerapi.Binding{}, err
	}
	b.Logger.Info("service binding created", lager.Data{
		"instance-id": instanceId,
		"binding-id":  bindingId,
	})
	return brokerapi.Binding{}, err
}

func (b *serviceBroker) Unbind(context context.Context, instanceId, bindingId string, details brokerapi.UnbindDetails) error {
	var serivceBinding model.ServiceBinding
	err := b.Db.First(&serivceBinding, "binding_id = ?", bindingId).Error
	if err != nil {
		return err
	}

	// remove your service binding here

	err = b.Db.Delete(&serivceBinding).Error

	return err
}

func (b *serviceBroker) LastOperation(context context.Context, instanceID, operationData string) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, nil
}
