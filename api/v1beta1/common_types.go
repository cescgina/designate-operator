/*
Copyright 2020 Red Hat

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"github.com/openstack-k8s-operators/lib-common/modules/common/util"
	corev1 "k8s.io/api/core/v1"
)

const (
	// Container image fall-back defaults

	// DesignateAPIContainerImage is the fall-back container image for DesignateAPI
	DesignateAPIContainerImage = "quay.io/podified-antelope-centos9/openstack-designate-api:current-podified"
	// DesignateCentralContainerImage is the fall-back container image for DesignateCentral
	DesignateCentralContainerImage = "quay.io/podified-antelope-centos9/openstack-designate-central:current-podified"
	// DesignateMdnsContainerImage is the fall-back container image for DesignateMdns
	DesignateMdnsContainerImage = "quay.io/podified-antelope-centos9/openstack-designate-mdns:current-podified"
	// DesignateProducerContainerImage is the fall-back container image for DesignateProducer
	DesignateProducerContainerImage = "quay.io/podified-antelope-centos9/openstack-designate-producer:current-podified"
	// DesignateWorkerContainerImage is the fall-back container image for DesignateWorker
	DesignateWorkerContainerImage = "quay.io/podified-antelope-centos9/openstack-designate-worker:current-podified"
)

// DesignateTemplate defines common input parameters used by all Designate services
type DesignateTemplate struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=designate
	// ServiceUser - optional username used for this service to register in designate
	ServiceUser string `json:"serviceUser"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=designate
	// DatabaseUser - optional username used for designate DB, defaults to designate
	// TODO: -> implement needs work in mariadb-operator, right now only designate
	DatabaseUser string `json:"databaseUser"`

	// +kubebuilder:validation:Optional
	// Secret containing OpenStack password information for DesignateDatabasePassword
	Secret string `json:"secret"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default={database: DesignateDatabasePassword, service: DesignatePassword}
	// PasswordSelectors - Selectors to identify the DB and ServiceUser password from the Secret
	PasswordSelectors PasswordSelector `json:"passwordSelectors"`
}

// DesignateServiceTemplate defines the input parameters that can be defined for a given
// Designate service
type DesignateServiceTemplate struct {
	// +kubebuilder:validation:Required
	// ContainerImage - Designate Container Image URL (will be set to environmental default if empty)
	ContainerImage string `json:"containerImage"`

	// +kubebuilder:validation:Optional
	// NodeSelector to target subset of worker nodes running this service. Setting here overrides
	// any global NodeSelector settings within the Designate CR.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// +kubebuilder:validation:Optional
	// Debug - enable debug for different deploy stages. If an init container is used, it runs and the
	// actual action pod gets started with sleep infinity
	Debug DesignateServiceDebug `json:"debug,omitempty"`

	// +kubebuilder:validation:Optional
	// CustomServiceConfig - customize the service config using this parameter to change service defaults,
	// or overwrite rendered information using raw OpenStack config format. The content gets added to
	// to /etc/<service>/<service>.conf.d directory as a custom config file.
	CustomServiceConfig string `json:"customServiceConfig,omitempty"`

	// +kubebuilder:validation:Optional
	// CustomServiceConfigSecrets - customize the service config using this parameter to specify Secrets
	// that contain sensitive service config data. The content of each Secret gets added to the
	// /etc/<service>/<service>.conf.d directory as a custom config file.
	CustomServiceConfigSecrets []string `json:"customServiceConfigSecrets,omitempty"`

	// +kubebuilder:validation:Optional
	// ConfigOverwrite - interface to overwrite default config files like e.g. policy.json.
	// But can also be used to add additional files. Those get added to the service config dir in /etc/<service> .
	// TODO: -> implement
	DefaultConfigOverwrite map[string]string `json:"defaultConfigOverwrite,omitempty"`

	// +kubebuilder:validation:Optional
	// Resources - Compute Resources required by this service (Limits/Requests).
	// https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	// NetworkAttachments is a list of NetworkAttachment resource names to expose the services to the given network
	NetworkAttachments []string `json:"networkAttachments,omitempty"`
}

// PasswordSelector to identify the DB and AdminUser password from the Secret
type PasswordSelector struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="DesignateDatabasePassword"
	// Database - Selector to get the designate database user password from the Secret
	// TODO: not used, need change in mariadb-operator
	Database string `json:"database"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="DesignatePassword"
	// Service - Selector to get the designate service password from the Secret
	Service string `json:"service"`
}

// DesignateDebug indicates whether certain stages of Designate deployment should
// pause in debug mode
type DesignateDebug struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	// dbInitContainer enable debug (waits until /tmp/stop-init-container disappears)
	DBInitContainer bool `json:"dbInitContainer"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	// dbSync enable debug
	DBSync bool `json:"dbSync"`
}

// DesignateServiceDebug indicates whether certain stages of Designate service
// deployment should pause in debug mode
type DesignateServiceDebug struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	// initContainer enable debug (waits until /tmp/stop-init-container disappears)
	InitContainer bool `json:"initContainer"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	// service enable debug
	Service bool `json:"service"`
}

// SetupDefaults - initializes any CRD field defaults based on environment variables (the defaulting mechanism itself is implemented via webhooks)
func SetupDefaults() {
	// Acquire environmental defaults and initialize Designate defaults with them
	designateDefaults := DesignateDefaults{
		APIContainerImageURL:      util.GetEnvVar("DESIGNATE_API_IMAGE_URL_DEFAULT", DesignateAPIContainerImage),
		CentralContainerImageURL:  util.GetEnvVar("DESIGNATE_CENTRAL_IMAGE_URL_DEFAULT", DesignateCentralContainerImage),
		MdnsContainerImageURL:     util.GetEnvVar("DESIGNATE_MDNS_IMAGE_URL_DEFAULT", DesignateMdnsContainerImage),
		ProducerContainerImageURL: util.GetEnvVar("DESIGNATE_PRODUCER_IMAGE_URL_DEFAULT", DesignateProducerContainerImage),
		WorkerContainerImageURL:   util.GetEnvVar("DESIGNATE_WORKER_IMAGE_URL_DEFAULT", DesignateWorkerContainerImage),
	}

	SetupDesignateDefaults(designateDefaults)
}
