/*
Copyright 2026 Open Cluster Management.

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

// Package v1beta1 holds the part of the argocd-operator ArgoCD API
// (argoproj.io/v1beta1) that this repo builds and reads.
//
// Type names, field names and JSON tags match
// github.com/argoproj-labs/argocd-operator/api/v1beta1 at v0.18.0, the operator
// version this repo ships (ARGOCD_OPERATOR_IMAGE in the Makefile). Only the
// fields this repo uses are declared, so a marshalled ArgoCD carries exactly
// those fields and nothing else.
//
// The operator's own package is not imported because v0.18.0 does not compile
// with controller-runtime v0.23 or later. TestTypesMatchOperatorCRD checks every
// field here against the bundled operator CRD, so bump both together
// (see internal/addon/OPERATOR_SYNC.md).
package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupVersion is the API group and version of the ArgoCD CR.
var GroupVersion = schema.GroupVersion{Group: "argoproj.io", Version: "v1beta1"}

// ArgoCDKind is the kind of the ArgoCD CR.
const ArgoCDKind = "ArgoCD"

// ArgoCD is the argocd-operator ArgoCD custom resource.
type ArgoCD struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ArgoCDSpec `json:"spec,omitempty"`
}

// ArgoCDSpec is the subset of the ArgoCD spec this repo uses.
type ArgoCDSpec struct {
	// Server defines the options for the Argo CD server component.
	Server ArgoCDServerSpec `json:"server,omitempty"`

	// SourceNamespaces lists the namespaces Applications may be created in.
	SourceNamespaces []string `json:"sourceNamespaces,omitempty"`

	// ArgoCDAgent defines the argocd-agent principal and agent components.
	ArgoCDAgent *ArgoCDAgentSpec `json:"argoCDAgent,omitempty"`
}

// ArgoCDServerSpec is the subset of the Argo CD server options this repo uses.
type ArgoCDServerSpec struct {
	// Enabled turns the Argo CD server component on or off.
	Enabled *bool `json:"enabled,omitempty"`
}

// ArgoCDAgentSpec defines the argocd-agent components.
type ArgoCDAgentSpec struct {
	// Principal is the principal component, which runs on the hub.
	Principal *PrincipalSpec `json:"principal,omitempty"`

	// Agent is the agent component, which runs on a managed cluster.
	Agent *AgentSpec `json:"agent,omitempty"`
}

// PrincipalSpec is the subset of the principal options this repo reads.
type PrincipalSpec struct {
	// Image is the container image of the principal.
	Image string `json:"image,omitempty"`
}

// AgentSpec is the subset of the agent options this repo sets.
type AgentSpec struct {
	// Enabled turns the agent component on or off.
	Enabled *bool `json:"enabled,omitempty"`

	// AllowedNamespaces lists the namespaces the agent may manage Applications in.
	AllowedNamespaces []string `json:"allowedNamespaces,omitempty"`

	// Creds is the credential identifier the agent authenticates with.
	Creds string `json:"creds,omitempty"`

	// LogLevel is the agent log level.
	LogLevel string `json:"logLevel,omitempty"`

	// LogFormat is the agent log format.
	LogFormat string `json:"logFormat,omitempty"`

	// Image is the container image of the agent.
	Image string `json:"image,omitempty"`

	// Client defines how the agent connects to the principal.
	Client *AgentClientSpec `json:"client,omitempty"`

	// Redis defines the Redis the agent uses.
	Redis *AgentRedisSpec `json:"redis,omitempty"`

	// TLS defines the agent's client certificate and CA.
	TLS *AgentTLSSpec `json:"tls,omitempty"`

	// DestinationBasedMapping routes Applications to agents by spec.destination.name.
	DestinationBasedMapping *DestinationBasedMappingSpec `json:"destinationBasedMapping,omitempty"`
}

// AgentClientSpec defines how the agent connects to the principal.
type AgentClientSpec struct {
	// PrincipalServerAddress is the address of the principal.
	PrincipalServerAddress string `json:"principalServerAddress,omitempty"`

	// PrincipalServerPort is the port of the principal.
	PrincipalServerPort string `json:"principalServerPort,omitempty"`

	// Mode is the agent mode, "managed" or "autonomous".
	Mode string `json:"mode,omitempty"`

	// EnableWebSocket makes the agent use WebSocket for streaming.
	EnableWebSocket *bool `json:"enableWebSocket,omitempty"`

	// EnableCompression turns on compression of the agent connection.
	EnableCompression *bool `json:"enableCompression,omitempty"`

	// KeepAliveInterval is the keep-alive ping interval of the agent connection.
	KeepAliveInterval string `json:"keepAliveInterval,omitempty"`
}

// AgentRedisSpec defines the Redis the agent uses.
type AgentRedisSpec struct {
	// ServerAddress is the address of the Redis server.
	ServerAddress string `json:"serverAddress,omitempty"`
}

// AgentTLSSpec defines the agent's client certificate and CA.
type AgentTLSSpec struct {
	// SecretName is the Secret holding the agent client certificate.
	SecretName string `json:"secretName,omitempty"`

	// RootCASecretName is the Secret holding the CA certificate.
	RootCASecretName string `json:"rootCASecretName,omitempty"`

	// Insecure skips verification of the principal's certificate.
	Insecure *bool `json:"insecure,omitempty"`
}

// DestinationBasedMappingSpec defines destination-based mapping on the agent.
type DestinationBasedMappingSpec struct {
	// Enabled turns destination-based mapping on.
	Enabled *bool `json:"enabled,omitempty"`

	// CreateNamespace creates missing target namespaces on the managed cluster.
	CreateNamespace *bool `json:"createNamespace,omitempty"`
}
