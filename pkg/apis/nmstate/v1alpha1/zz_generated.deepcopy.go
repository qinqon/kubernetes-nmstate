// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkConfigurationPolicy) DeepCopyInto(out *NodeNetworkConfigurationPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkConfigurationPolicy.
func (in *NodeNetworkConfigurationPolicy) DeepCopy() *NodeNetworkConfigurationPolicy {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkConfigurationPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeNetworkConfigurationPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkConfigurationPolicyList) DeepCopyInto(out *NodeNetworkConfigurationPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NodeNetworkConfigurationPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkConfigurationPolicyList.
func (in *NodeNetworkConfigurationPolicyList) DeepCopy() *NodeNetworkConfigurationPolicyList {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkConfigurationPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeNetworkConfigurationPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkConfigurationPolicyMatch) DeepCopyInto(out *NodeNetworkConfigurationPolicyMatch) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkConfigurationPolicyMatch.
func (in *NodeNetworkConfigurationPolicyMatch) DeepCopy() *NodeNetworkConfigurationPolicyMatch {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkConfigurationPolicyMatch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkConfigurationPolicySpec) DeepCopyInto(out *NodeNetworkConfigurationPolicySpec) {
	*out = *in
	in.Match.DeepCopyInto(&out.Match)
	if in.DesiredState != nil {
		in, out := &in.DesiredState, &out.DesiredState
		*out = make(State, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkConfigurationPolicySpec.
func (in *NodeNetworkConfigurationPolicySpec) DeepCopy() *NodeNetworkConfigurationPolicySpec {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkConfigurationPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkConfigurationPolicyStatus) DeepCopyInto(out *NodeNetworkConfigurationPolicyStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkConfigurationPolicyStatus.
func (in *NodeNetworkConfigurationPolicyStatus) DeepCopy() *NodeNetworkConfigurationPolicyStatus {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkConfigurationPolicyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkState) DeepCopyInto(out *NodeNetworkState) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkState.
func (in *NodeNetworkState) DeepCopy() *NodeNetworkState {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeNetworkState) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkStateList) DeepCopyInto(out *NodeNetworkStateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NodeNetworkState, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkStateList.
func (in *NodeNetworkStateList) DeepCopy() *NodeNetworkStateList {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkStateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeNetworkStateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkStateSpec) DeepCopyInto(out *NodeNetworkStateSpec) {
	*out = *in
	if in.DesiredState != nil {
		in, out := &in.DesiredState, &out.DesiredState
		*out = make(State, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkStateSpec.
func (in *NodeNetworkStateSpec) DeepCopy() *NodeNetworkStateSpec {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkStateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkStateStatus) DeepCopyInto(out *NodeNetworkStateStatus) {
	*out = *in
	if in.CurrentState != nil {
		in, out := &in.CurrentState, &out.CurrentState
		*out = make(State, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkStateStatus.
func (in *NodeNetworkStateStatus) DeepCopy() *NodeNetworkStateStatus {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkStateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in State) DeepCopyInto(out *State) {
	{
		in := &in
		*out = make(State, len(*in))
		copy(*out, *in)
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new State.
func (in State) DeepCopy() State {
	if in == nil {
		return nil
	}
	out := new(State)
	in.DeepCopyInto(out)
	return *out
}
