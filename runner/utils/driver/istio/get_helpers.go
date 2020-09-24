package istio

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kedacore/keda/api/v1alpha1"
	"github.com/segmentio/ksuid"
	"github.com/spaceuptech/helpers"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spaceuptech/space-cloud/runner/model"
)

func (i *Istio) getPreviousVirtualServiceIfExists(ctx context.Context, ns, service string) (*v1alpha3.VirtualService, error) {
	prevVirtualService, err := i.istio.NetworkingV1alpha3().VirtualServices(ns).Get(ctx, getVirtualServiceName(service), metav1.GetOptions{})
	if kubeErrors.IsNotFound(err) {
		// We'll simple send `nil` if the virtual service did not actually exist. This is important since it indicates that
		// a virtual service needs to be created
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return prevVirtualService, nil
}

func (i *Istio) getServiceDeployments(ctx context.Context, ns, serviceID string) (*appsv1.DeploymentList, error) {
	return i.kube.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("app=%s", serviceID)})
}

func (i *Istio) getServiceDeploymentsCount(ctx context.Context, ns, serviceID string) (int, error) {
	deployments, err := i.getServiceDeployments(ctx, ns, serviceID)
	if err != nil {
		return 0, err
	}

	return len(deployments.Items), nil
}

func (i *Istio) getVirtualServices(ctx context.Context, ns string) (*v1alpha3.VirtualServiceList, error) {
	return i.istio.NetworkingV1alpha3().VirtualServices(ns).List(ctx, metav1.ListOptions{})
}

func (i *Istio) getServiceScaleConfig(ctx context.Context, ns, serviceID string) (*v1alpha1.ScaledObjectList, error) {
	return i.keda.KedaV1alpha1().ScaledObjects(ns).List(ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("app.kubernetes.io/name=%s", serviceID)})
}

func (i *Istio) getKedaTriggerAuthsForVersion(ctx context.Context, ns, serviceID, version string) (*v1alpha1.TriggerAuthenticationList, error) {
	labelSelector := fmt.Sprintf("app.kubernetes.io/name=%s,app.kubernetes.io/version=%s,app.kubernetes.io/managed-by=space-cloud", serviceID, version)
	return i.keda.KedaV1alpha1().TriggerAuthentications(ns).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
}

func (i *Istio) getAllVersionScaleConfig(ctx context.Context, ns, serviceID string) (map[string]model.AutoScaleConfig, error) {
	// Get all deployments of the provided service
	scaledOjectList, err := i.getServiceScaleConfig(ctx, ns, serviceID)
	if err != nil {
		return nil, err
	}

	// Throw error if the deployment contains no config at all
	if len(scaledOjectList.Items) == 0 {
		return nil, helpers.Logger.LogError(helpers.GetRequestID(ctx), fmt.Sprintf("no versions of service (%s) has been deployed", serviceID), nil, nil)
	}

	// Load the scale config of each version
	c := make(map[string]model.AutoScaleConfig, len(scaledOjectList.Items))
	for _, scaledObject := range scaledOjectList.Items {
		c[scaledObject.Labels["app.kubernetes.io/version"]] = model.AutoScaleConfig{
			MaxReplicas:      *scaledObject.Spec.MaxReplicaCount,
			MinReplicas:      *scaledObject.Spec.MinReplicaCount,
			PollingInterval:  *scaledObject.Spec.PollingInterval,
			CoolDownInterval: *scaledObject.Spec.CooldownPeriod,
		}
	}

	return c, nil
}

func getScaleConfigFromDeployment(deployment appsv1.Deployment) *model.AutoScaleConfig {
	autoscale := getDefaultAutoScaleConfig()

	minReplicas, err := strconv.Atoi(deployment.Annotations["minReplicas"])
	if err == nil {
		autoscale.MinReplicas = int32(minReplicas)
	}
	maxReplicas, err := strconv.Atoi(deployment.Annotations["maxReplicas"])
	if err == nil {
		autoscale.MaxReplicas = int32(maxReplicas)
	}

	trigger := model.AutoScaleTrigger{}

	mode := deployment.Annotations["mode"]
	if mode == "" {
		mode = "per-second"
	}
	switch mode {
	case "", "per-second":
		trigger.Name = "requests-per-second"
		trigger.Type = "requests-per-second"
	case "parallel":
		trigger.Name = "active-requests"
		trigger.Type = "active-requests"
	}

	trigger.MetaData = map[string]string{"target": deployment.Annotations["concurrency"]}

	autoscale.Triggers = []model.AutoScaleTrigger{trigger}
	return autoscale
}

func getScaleConfigFromKedaConfig(service, version string, scaledObjList []v1alpha1.ScaledObject, triggerAuthList []v1alpha1.TriggerAuthentication) *model.AutoScaleConfig {
	// See if a valid scaled object exists for given deployment
	var scaledObject *v1alpha1.ScaledObject
	for _, obj := range scaledObjList {
		if obj.Labels["app.kubernetes.io/name"] == service && obj.Labels["app.kubernetes.io/version"] == version {
			scaledObject = &obj
			break
		}
	}

	// Return nil value if no keda scaled object found
	if scaledObject == nil {
		return nil
	}

	// Prepare auto scaling triggers
	autoScaleTriggers := make([]model.AutoScaleTrigger, len(scaledObject.Spec.Triggers))
	for i, trigger := range scaledObject.Spec.Triggers {
		// Prepare trigger auths
		triggerType, metadata := getTriggerTypeFromTrigger(trigger)

		t := model.AutoScaleTrigger{
			Name:     trigger.Name,
			Type:     triggerType,
			MetaData: metadata,
		}

		// Set authorisation ref if its provided
		if trigger.AuthenticationRef != nil && trigger.AuthenticationRef.Name != "" {
			t.AuthenticatedRef = getTriggerAuthForTrigger(trigger.AuthenticationRef.Name, triggerAuthList)
		}

		autoScaleTriggers[i] = t
	}

	// Get triggers for cpu and memory usage as well
	autoScaleTriggers = append(autoScaleTriggers, getTriggersFromAdvancedConfig(scaledObject.Spec.Advanced)...)

	// Prepare an auto scale object
	autoScale := &model.AutoScaleConfig{
		CoolDownInterval: *scaledObject.Spec.CooldownPeriod,
		PollingInterval:  *scaledObject.Spec.PollingInterval,
		MinReplicas:      *scaledObject.Spec.MinReplicaCount,
		MaxReplicas:      *scaledObject.Spec.MaxReplicaCount,
		Triggers:         autoScaleTriggers,
	}

	return autoScale
}

func getTriggerTypeFromTrigger(trigger v1alpha1.ScaleTriggers) (string, map[string]string) {
	if trigger.Type == "external-push" {
		if scaler, p := trigger.Metadata["scaler"]; p && scaler == "space-cloud.io/scaler" {
			if triggerType, p := trigger.Metadata["type"]; p {
				return triggerType, map[string]string{"target": trigger.Metadata["target"]}
			}
		}
	}

	return trigger.Type, trigger.Metadata
}

func getTriggerAuthForTrigger(trigger string, refs []v1alpha1.TriggerAuthentication) *model.AutoScaleAuthRef {
	var triggerAuth *v1alpha1.TriggerAuthentication
	for _, ref := range refs {
		if ref.Name == trigger {
			triggerAuth = &ref
			break
		}
	}

	// Return nil if no trigger auth object is found
	if triggerAuth == nil {
		return nil
	}

	// Prepare mapping
	var secretName string
	mapping := make([]model.AutoScaleAuthRefMapping, len(triggerAuth.Spec.SecretTargetRef))
	for i, ref := range triggerAuth.Spec.SecretTargetRef {
		mapping[i] = model.AutoScaleAuthRefMapping{
			Parameter: ref.Parameter,
			Key:       ref.Key,
		}

		// Extract secret name as well. We will have a consistent secret name for each parameter
		secretName = ref.Name
	}

	// Finally prepare and return the auto scaling auth reference
	return &model.AutoScaleAuthRef{
		SecretName:    secretName,
		SecretMapping: mapping,
	}
}

func getTriggersFromAdvancedConfig(advancedConfig *v1alpha1.AdvancedConfig) []model.AutoScaleTrigger {
	triggers := make([]model.AutoScaleTrigger, 0)

	// Return if advanced config is nil
	if advancedConfig == nil {
		return triggers
	}

	// Return if hpa config is nil
	hpa := advancedConfig.HorizontalPodAutoscalerConfig
	if hpa == nil {
		return triggers
	}

	// Iterate over resource metrics to get resource scalers
	for _, config := range hpa.ResourceMetrics {
		triggers = append(triggers, model.AutoScaleTrigger{
			Name: string(config.Name),
			Type: string(config.Name),
			MetaData: map[string]string{
				"target": strconv.Itoa(int(*config.Target.AverageUtilization)),
			},
		})
	}

	return triggers
}

func extractPreferredServiceAffinityObject(arr []v1.WeightedPodAffinityTerm, multiplier int32) []model.Affinity {
	affinities := []model.Affinity{}
	for _, preferredSchedulingTerm := range arr {
		matchExpression := []model.MatchExpressions{}
		for _, expression := range preferredSchedulingTerm.PodAffinityTerm.LabelSelector.MatchExpressions {
			matchExpression = append(matchExpression, model.MatchExpressions{
				Key:       expression.Key,
				Values:    expression.Values,
				Attribute: "label",
				Operator:  string(expression.Operator),
			})
		}
		affinities = append(affinities, model.Affinity{
			ID:               ksuid.New().String(),
			Type:             model.AffinityTypeService,
			Weight:           preferredSchedulingTerm.Weight * multiplier,
			Operator:         model.AffinityOperatorPreferred,
			TopologyKey:      preferredSchedulingTerm.PodAffinityTerm.TopologyKey,
			Projects:         preferredSchedulingTerm.PodAffinityTerm.Namespaces,
			MatchExpressions: matchExpression,
		})
	}
	return affinities
}

func extractRequiredServiceAffinityObject(arr []v1.PodAffinityTerm, multiplier int32) []model.Affinity {
	affinities := []model.Affinity{}
	for _, preferredSchedulingTerm := range arr {
		matchExpression := []model.MatchExpressions{}
		for _, expression := range preferredSchedulingTerm.LabelSelector.MatchExpressions {
			matchExpression = append(matchExpression, model.MatchExpressions{
				Key:       expression.Key,
				Values:    expression.Values,
				Attribute: "label",
				Operator:  string(expression.Operator),
			})
		}
		affinities = append(affinities, model.Affinity{
			ID:               ksuid.New().String(),
			Type:             model.AffinityTypeService,
			Weight:           100 * multiplier,
			Operator:         model.AffinityOperatorRequired,
			TopologyKey:      preferredSchedulingTerm.TopologyKey,
			Projects:         preferredSchedulingTerm.Namespaces,
			MatchExpressions: matchExpression,
		})
	}
	return affinities
}

func extractPreferredNodeAffinityObject(arr []v1.PreferredSchedulingTerm) []model.Affinity {
	affinities := []model.Affinity{}
	for _, preferredSchedulingTerm := range arr {
		matchExpression := []model.MatchExpressions{}
		for _, expression := range preferredSchedulingTerm.Preference.MatchExpressions {
			matchExpression = append(matchExpression, model.MatchExpressions{
				Key:       expression.Key,
				Values:    expression.Values,
				Attribute: "label",
				Operator:  string(expression.Operator),
			})
		}
		affinities = append(affinities, model.Affinity{
			ID:               ksuid.New().String(),
			Type:             model.AffinityTypeNode,
			Weight:           preferredSchedulingTerm.Weight,
			Operator:         model.AffinityOperatorPreferred,
			MatchExpressions: matchExpression,
		})
	}
	return affinities
}

func extractRequiredNodeAffinityObject(arr []v1.NodeSelectorTerm) []model.Affinity {
	affinities := []model.Affinity{}
	for _, nodeSelectorTerm := range arr {
		matchExpression := []model.MatchExpressions{}
		for _, expression := range nodeSelectorTerm.MatchExpressions {
			matchExpression = append(matchExpression, model.MatchExpressions{
				Key:       expression.Key,
				Values:    expression.Values,
				Attribute: "label",
				Operator:  string(expression.Operator),
			})
		}
		affinities = append(affinities, model.Affinity{
			ID:               ksuid.New().String(),
			Type:             model.AffinityTypeNode,
			Operator:         model.AffinityOperatorRequired,
			MatchExpressions: matchExpression,
		})
	}
	return affinities
}
