package shared

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NodeNetworkStateStatus is the status of the NodeNetworkState of a specific node
// +k8s:openapi-gen=true
type NodeNetworkStateStatus struct {
	CurrentState             State       `json:"currentState,omitempty"`
	LastSuccessfulUpdateTime metav1.Time `json:"lastSuccessfulUpdateTime,omitempty"`

	Conditions ConditionList `json:"conditions,omitempty" optional:"true"`
}

const (
	NodeNetworkStateConditionAvailable ConditionType = "Available"
	NodeNetworkStateConditionFailing   ConditionType = "Failing"
)

const (
	NodeNetworkStateConditionFailedToConfigure      ConditionReason = "FailedToConfigure"
	NodeNetworkStateConditionSuccessfullyConfigured ConditionReason = "SuccessfullyConfigured"
)
