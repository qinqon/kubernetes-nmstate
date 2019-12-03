package conditions

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	nmstatev1alpha1 "github.com/nmstate/kubernetes-nmstate/pkg/apis/nmstate/v1alpha1"
)

type Enactment struct {
	Conditions nmstatev1alpha1.ConditionList
	Policy     string
}

type Policy struct {
	Conditions   nmstatev1alpha1.ConditionList
	NodeSelector map[string]string
}

func enactment(policy string, conditionsSetter func(*nmstatev1alpha1.ConditionList, string)) Enactment {
	conditions := nmstatev1alpha1.ConditionList{}
	conditionsSetter(&conditions, "")
	return Enactment{
		Policy:     policy,
		Conditions: conditions,
	}
}

func policy(conditionsSetter func(*nmstatev1alpha1.ConditionList, string), message string, nodeSelector map[string]string) Policy {
	conditions := nmstatev1alpha1.ConditionList{}
	conditionsSetter(&conditions, message)
	return Policy{
		Conditions:   conditions,
		NodeSelector: nodeSelector,
	}
}

func cleanTimestamps(conditions nmstatev1alpha1.ConditionList) nmstatev1alpha1.ConditionList {
	dummyTime := metav1.Time{time.Unix(0, 0)}
	for i, _ := range conditions {
		conditions[i].LastHeartbeatTime = dummyTime
		conditions[i].LastTransitionTime = dummyTime
	}
	return conditions
}

func allNodes() map[string]string {
	return map[string]string{}
}

func forNode(node string) map[string]string {
	return map[string]string{
		"kubernetes.io/hostname": node,
	}
}

var _ = Describe("Conditions manager", func() {
	type ConditionsCase struct {
		Enactments    []Enactment
		NumberOfNodes int
		Policy        Policy
	}
	DescribeTable("the policy overall condition",
		func(c ConditionsCase) {
			objs := []runtime.Object{}
			s := scheme.Scheme
			s.AddKnownTypes(nmstatev1alpha1.SchemeGroupVersion,
				&nmstatev1alpha1.NodeNetworkConfigurationPolicy{},
				&nmstatev1alpha1.NodeNetworkConfigurationEnactment{},
				&nmstatev1alpha1.NodeNetworkConfigurationEnactmentList{},
			)
			policy := nmstatev1alpha1.NodeNetworkConfigurationPolicy{}
			policy.Name = "policy1"
			policy.Spec.NodeSelector = c.Policy.NodeSelector

			for i, enactment := range c.Enactments {
				nnce := nmstatev1alpha1.NodeNetworkConfigurationEnactment{}
				nodeName := fmt.Sprintf("node%d", i)
				nnce.Name = nmstatev1alpha1.EnactmentKey(nodeName, enactment.Policy).Name
				nnce.Status.Conditions = enactment.Conditions
				nnce.Labels = map[string]string{
					"policy": enactment.Policy,
				}
				objs = append(objs, &nnce)
			}
			for i := 1; i <= c.NumberOfNodes; i++ {
				node := corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: fmt.Sprintf("node%d", i),
					},
				}
				node.Labels = map[string]string{
					"kubernetes.io/hostname": fmt.Sprintf("node%d", i),
				}
				objs = append(objs, &node)
			}

			objs = append(objs, &policy)

			manager := Manager{
				client: fake.NewFakeClientWithScheme(s, objs...),
				policy: &policy,
			}
			err := manager.refreshPolicyConditions()
			Expect(err).ToNot(HaveOccurred())
			Expect(cleanTimestamps(policy.Status.Conditions)).To(ConsistOf(cleanTimestamps(c.Policy.Conditions)))
		},
		Entry("when all enactments are progressing then policy is progressing", ConditionsCase{
			Enactments: []Enactment{
				enactment("policy1", setEnactmentProgressing),
				enactment("policy1", setEnactmentProgressing),
				enactment("policy1", setEnactmentProgressing),
			},
			NumberOfNodes: 3,
			Policy:        policy(setPolicyProgressing, "3/3 nodes configuring", allNodes()),
		}),
		Entry("when all enactments are success then policy is success", ConditionsCase{
			Enactments: []Enactment{
				enactment("policy1", setEnactmentSuccess),
				enactment("policy1", setEnactmentSuccess),
				enactment("policy1", setEnactmentSuccess),
			},
			NumberOfNodes: 3,
			Policy:        policy(setPolicySuccess, "All 3 nodes successfully configured", allNodes()),
		}),
		Entry("when partial enactments are success then policy is progressing", ConditionsCase{
			Enactments: []Enactment{
				enactment("policy1", setEnactmentSuccess),
				enactment("policy1", setEnactmentSuccess),
				enactment("policy1", setEnactmentSuccess),
			},
			NumberOfNodes: 4,
			Policy:        policy(setPolicyProgressing, "1/4 nodes not started to configure yet", allNodes()),
		}),
		Entry("when enactments are progressing/success then policy is progressing", ConditionsCase{
			Enactments: []Enactment{
				enactment("policy1", setEnactmentSuccess),
				enactment("policy1", setEnactmentProgressing),
				enactment("policy1", setEnactmentSuccess),
			},
			NumberOfNodes: 3,
			Policy:        policy(setPolicyProgressing, "1/3 nodes configuring", allNodes()),
		}),
		Entry("when enactments are failed/progressing/success then policy is degraded", ConditionsCase{
			Enactments: []Enactment{
				enactment("policy1", setEnactmentSuccess),
				enactment("policy1", setEnactmentProgressing),
				enactment("policy1", setEnactmentFailedToConfigure),
				enactment("policy1", setEnactmentSuccess),
			},
			NumberOfNodes: 4,
			Policy:        policy(setPolicyFailedToConfigure, "1/4 nodes failed to configure", allNodes()),
		}),
		Entry("when neither of enactments are matching then policy is neither at degraded/progressing/success ", ConditionsCase{
			Enactments: []Enactment{
				enactment("policy1", setEnactmentNodeSelectorNotMatching),
				enactment("policy1", setEnactmentNodeSelectorNotMatching),
				enactment("policy1", setEnactmentNodeSelectorNotMatching),
			},
			NumberOfNodes: 3,
			Policy:        policy(setPolicyNotMatching, "None of 3 nodes matches the policy", allNodes()),
		}),
		Entry("when some enactments are from different profile then policy conditions are not affected by them", ConditionsCase{
			Enactments: []Enactment{
				enactment("policy1", setEnactmentSuccess),
				enactment("policy2", setEnactmentProgressing),
				enactment("policy2", setEnactmentProgressing),
			},
			NumberOfNodes: 3,
			Policy:        policy(setPolicySuccess, "All 1 nodes successfully configured", forNode("node1")),
		}),
	)
})
