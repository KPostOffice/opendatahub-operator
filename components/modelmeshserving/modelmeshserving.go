package modelmeshserving

import (
	"github.com/opendatahub-io/opendatahub-operator/v2/components"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/common"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/deploy"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ComponentName  = "model-mesh"
	Path           = deploy.DefaultManifestPath + "/" + ComponentName + "/base"
	monitoringPath = deploy.DefaultManifestPath + "/" + "modelmesh-monitoring/base"
)

type ModelMeshServing struct {
	components.Component `json:""`
}

func (m *ModelMeshServing) GetComponentName() string {
	return ComponentName
}

// Verifies that Dashboard implements ComponentInterface
var _ components.ComponentInterface = (*ModelMeshServing)(nil)

func (m *ModelMeshServing) IsEnabled() bool {
	return m.Enabled
}

func (m *ModelMeshServing) SetEnabled(enabled bool) {
	m.Enabled = enabled
}

func (m *ModelMeshServing) ReconcileComponent(owner metav1.Object, cli client.Client, scheme *runtime.Scheme, enabled bool, namespace string) error {

	// Update Default rolebinding
	err := common.UpdatePodSecurityRolebinding(cli, []string{"modelmesh", "modelmesh-controller", "odh-model-controller", "odh-prometheus-operator", "prometheus-custom"}, namespace)
	if err != nil {
		return err
	}
	err = deploy.DeployManifestsFromPath(owner, cli, ComponentName,
		Path,
		namespace,
		scheme, enabled)

	if err != nil {
		return err
	}

	// If modelmesh is deployed successfully, deploy modelmesh-monitoring
	err = deploy.DeployManifestsFromPath(owner, cli, ComponentName,
		monitoringPath,
		namespace,
		scheme, enabled)

	return err
}

func (in *ModelMeshServing) DeepCopyInto(out *ModelMeshServing) {
	*out = *in
	out.Component = in.Component
}