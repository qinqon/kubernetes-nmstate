#!/bin/bash
#
# Per-node bootstrap for the kcli-based kubernetes-nmstate dev cluster.
# Runs on each CentOS Stream 9 node (node01..node0N) over `kcli ssh`.
#
# Installs the container runtime, kubeadm/kubelet, the latest NetworkManager
# with the OVS plugin and Open vSwitch, and configures the cluster to trust the
# host-side image registry exposed as registry:5000.
#
# Usage: bootstrap-common.sh <k8s_version> <nm_version> <registry_host_ip>
#   k8s_version       e.g. 1.34 (minor) used to pick the pkgs.k8s.io repo
#   nm_version        'latest' to pull NetworkManager from copr, else anything
#   registry_host_ip  host IP reachable from the node for registry:5000
#
set -ex

K8S_VERSION="${1:?missing k8s version}"
NM_VERSION="${2:-}"
REGISTRY_HOST_IP="${3:?missing registry host ip}"

# pkgs.k8s.io repos are keyed by minor version (vMAJOR.MINOR)
K8S_MINOR="v$(echo "${K8S_VERSION}" | awk -F. '{print $1"."$2}')"

echo "### Disabling swap and SELinux enforcement"
sudo swapoff -a || true
sudo sed -i '/\sswap\s/d' /etc/fstab || true
sudo setenforce 0 || true
sudo sed -i 's/^SELINUX=enforcing/SELINUX=permissive/' /etc/selinux/config || true

echo "### Kernel modules and sysctl for kubernetes networking"
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF
sudo modprobe overlay
sudo modprobe br_netfilter
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF
sudo sysctl --system

echo "### NetworkManager + OVS"
if [[ "${NM_VERSION}" == "latest" ]]; then
    echo "Installing NetworkManager from copr networkmanager/NetworkManager-main"
    sudo dnf install -y dnf-plugins-core
    sudo dnf copr enable -y networkmanager/NetworkManager-main
fi
# Open vSwitch is shipped by the CentOS NFV SIG (versioned openvswitchX.Y
# packages), not by BaseOS/AppStream. Enable the repo and pick the newest
# available openvswitch package.
sudo dnf install -y centos-release-nfv-openvswitch
ovs_pkg=$(sudo dnf -q list available 'openvswitch3.*' 2>/dev/null \
    | awk '/^openvswitch3\./ {print $1}' | sort -V | tail -1)
ovs_pkg=${ovs_pkg:-openvswitch}
echo "Using Open vSwitch package: ${ovs_pkg}"
sudo dnf install -y --allowerasing NetworkManager NetworkManager-ovs "${ovs_pkg}"
sudo systemctl daemon-reload
sudo systemctl enable --now openvswitch
sudo systemctl restart openvswitch
# Keep using NetworkManager's internal DHCP client (never dhclient), matching
# the previous kubevirtci behaviour.
sudo rm -f /etc/NetworkManager/conf.d/002-dhclient.conf
sudo systemctl restart NetworkManager
# Persistent journal so logs survive node reboots (e2e reboot tests).
sudo mkdir -p /var/log/journal
sudo systemctl restart systemd-journald

echo "### Configuring host registry trust (registry:5000 -> ${REGISTRY_HOST_IP})"
if ! grep -q "registry$" /etc/hosts; then
    echo "${REGISTRY_HOST_IP} registry" | sudo tee -a /etc/hosts
fi

echo "### Installing containerd"
sudo dnf install -y dnf-plugins-core
sudo dnf config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo dnf install -y containerd.io
sudo mkdir -p /etc/containerd
containerd config default | sudo tee /etc/containerd/config.toml >/dev/null
# Use the systemd cgroup driver (required for kubelet) and an external hosts.d.
sudo sed -i 's/SystemdCgroup = false/SystemdCgroup = true/' /etc/containerd/config.toml
# containerd may emit config_path with single or double quotes depending on
# version; normalise both to our certs.d directory.
sudo sed -i -E "s#config_path = (''|\"\")#config_path = '/etc/containerd/certs.d'#g" /etc/containerd/config.toml
sudo mkdir -p /etc/containerd/certs.d/registry:5000
cat <<EOF | sudo tee /etc/containerd/certs.d/registry:5000/hosts.toml
server = "http://registry:5000"

[host."http://registry:5000"]
  capabilities = ["pull", "resolve"]
  skip_verify = true
EOF
sudo systemctl enable --now containerd
sudo systemctl restart containerd

echo "### Installing kubeadm/kubelet/kubectl ${K8S_MINOR}"
cat <<EOF | sudo tee /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://pkgs.k8s.io/core:/stable:/${K8S_MINOR}/rpm/
enabled=1
gpgcheck=1
gpgkey=https://pkgs.k8s.io/core:/stable:/${K8S_MINOR}/rpm/repodata/repomd.xml.key
exclude=kubelet kubeadm kubectl cri-tools kubernetes-cni
EOF
sudo dnf install -y --disableexcludes=kubernetes kubelet kubeadm kubectl cri-tools
sudo systemctl enable kubelet

echo "### crictl runtime endpoint"
cat <<EOF | sudo tee /etc/crictl.yaml
runtime-endpoint: unix:///run/containerd/containerd.sock
image-endpoint: unix:///run/containerd/containerd.sock
timeout: 10
EOF

echo "### Bootstrap complete on $(hostname)"
