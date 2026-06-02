#!/bin/bash
#
# Initialise the kubeadm control-plane on node01.
# Runs on node01 over `kcli ssh`.
#
# Usage: init-controlplane.sh <apiserver_ip> <pod_cidr> <primary_nic>
#
set -ex

APISERVER_IP="${1:?missing apiserver ip}"
POD_CIDR="${2:-10.244.0.0/16}"
PRIMARY_NIC="${3:?missing primary nic}"

sudo kubeadm init \
    --apiserver-advertise-address="${APISERVER_IP}" \
    --pod-network-cidr="${POD_CIDR}" \
    --node-name="$(hostname)"

# Make kubectl work for the default user and for root.
mkdir -p "${HOME}/.kube"
sudo cp -f /etc/kubernetes/admin.conf "${HOME}/.kube/config"
sudo chown "$(id -u):$(id -g)" "${HOME}/.kube/config"

# Deploy flannel, pinned to the primary NIC so it does not pick one of the
# secondary interfaces that share the same subnet.
curl -sSL -o /tmp/kube-flannel.yml \
    https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml
# Inject --iface=<primary> into the kube-flannel container args.
sudo sed -i "s#- --kube-subnet-mgr#- --kube-subnet-mgr\n        - --iface=${PRIMARY_NIC}#" /tmp/kube-flannel.yml
sudo kubectl --kubeconfig /etc/kubernetes/admin.conf apply -f /tmp/kube-flannel.yml

echo "### Control-plane ready on $(hostname)"
