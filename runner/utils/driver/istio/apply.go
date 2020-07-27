package istio

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/spaceuptech/space-cloud/runner/model"
)

// ApplyService deploys the service on istio
func (i *Istio) ApplyService(ctx context.Context, service *model.Service) error {
	// TODO: do we need to rollback on failure? rollback to previous version if it existed else remove. We also need to rollback the cache in this case

	ns := service.ProjectID

	// Set the default concurrency value to 50
	if service.Scale.Concurrency == 0 {
		service.Scale.Concurrency = 50
	}

	// Adjust the min scale in case of tcp based services. Min scale for tcp services need to be at least 1.
	adjustMinScale(service)

	// TODO: remove artifact store related code
	token, err := i.auth.GenerateTokenForArtifactStore(service.ID, service.ProjectID, service.Version)
	if err != nil {
		return err
	}

	// Get the list of secrets required for this service
	listOfSecrets, err := i.getSecrets(ctx, service)
	if err != nil {
		return err
	}

	// Get the scale info of all versions of the service
	prevVirtualService, err := i.getPreviousVirtualServiceIfExists(ctx, ns, service.ID)
	if err != nil {
		return err
	}

	// Create the appropriate kubernetes and istio objects
	kubeServiceAccount := generateServiceAccount(service)
	kubeDeployment := i.generateDeployment(service, token, listOfSecrets)
	kubeGeneralService := generateGeneralService(service)
	kubeInternalService := generateInternalService(service)
	istioVirtualService := i.updateVirtualService(service, prevVirtualService)
	istioGeneralDestRule := generateGeneralDestinationRule(service)
	istioInternalDestRule := generateInternalDestinationRule(service)
	istioAuthPolicy := generateAuthPolicy(service)
	istioSidecar := generateSidecarConfig(service)

	// Create a service account if it doesn't already exist
	logrus.Debugf("Create service account (%s) in %s", kubeServiceAccount.Name, ns)
	if err := i.createServiceAccountIfNotExist(ctx, ns, kubeServiceAccount); err != nil {
		return err
	}

	// Apply the deployment config
	logrus.Debugf("Applying deployment (%s) in %s", kubeDeployment.Name, ns)
	if err := i.applyDeployment(ctx, ns, kubeDeployment); err != nil {
		return err
	}

	// Create a global service if not exists. This is required for service discovery purposes only.
	logrus.Debugf("Creating general service service (%s) in %s if it doesn't already exists", kubeGeneralService.Name, ns)
	if err := i.createServiceIfNotExist(ctx, ns, kubeGeneralService); err != nil {
		return err
	}

	// Apply the internal service config
	logrus.Debugf("Applying internal service (%s) in %s", kubeInternalService.Name, ns)
	if err := i.applyService(ctx, ns, kubeInternalService); err != nil {
		return err
	}

	// Create the virtual service
	logrus.Debugf("Creating virtual service (%s) in %s if it doesn't already exist", istioVirtualService.Name, ns)
	if err := i.createVirtualServiceIfNotExist(ctx, ns, istioVirtualService); err != nil {
		return err
	}

	// Create the general destination rule config
	logrus.Debugf("Creating general destination rules (%s) in %s", istioGeneralDestRule.Name, ns)
	if err := i.createDestinationRulesIfNotExist(ctx, ns, istioGeneralDestRule); err != nil {
		return err
	}

	// Apply the internal destination rule config
	logrus.Debugf("Applying internal destination rules (%s) in %s", istioInternalDestRule.Name, ns)
	if err := i.applyDestinationRules(ctx, ns, istioInternalDestRule); err != nil {
		return err
	}

	// Apply the authorization policy config
	logrus.Debugf("Applying authorization policy (%s) in %s", istioAuthPolicy.Name, ns)
	if err := i.applyAuthorizationPolicy(ctx, ns, istioAuthPolicy); err != nil {
		return err
	}

	// Apply the sidecar config
	logrus.Debugf("Applying sidecar config (%s) in %s", istioSidecar.Name, ns)
	if err := i.applySidecar(ctx, ns, istioSidecar); err != nil {
		return err
	}

	logrus.Infof("Service (%s:%s) applied successfully", service.ProjectID, service.ID)
	return nil
}

// ApplyServiceRoutes sets the traffic splitting logic of each service
func (i *Istio) ApplyServiceRoutes(ctx context.Context, projectID, serviceID string, routes model.Routes) error {
	ns := projectID

	// Get the scale info of all versions of the service
	prevVirtualService, err := i.getPreviousVirtualServiceIfExists(ctx, ns, serviceID)
	if err != nil {
		return err
	}

	scaleConfig, err := i.getAllVersionScaleConfig(ctx, ns, serviceID)
	if err != nil {
		return err
	}

	virtualService, err := i.generateVirtualServiceBasedOnRoutes(projectID, serviceID, scaleConfig, routes, prevVirtualService)
	if err != nil {
		return err
	}

	return i.applyVirtualService(ctx, ns, virtualService)
}
