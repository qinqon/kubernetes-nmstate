#!/bin/bash
#
# Shared library for the kcli-based local development cluster.
# Replaces the previous cluster/kubevirtci.sh.
#
# The cluster is a set of plain libvirt VMs (node01..node0N) provisioned with
# kcli and bootstrapped into a kubeadm cluster by cluster/up.sh. The generated
# kubeconfig lives under ${KCLI_PATH} so e2e tooling can detect a kcli cluster.

# KUBEVIRT_PROVIDER is kept (now only used to derive the kubernetes version and
# to satisfy the "k8s" substring checks in the Makefile and e2e helpers).
export KUBEVIRT_PROVIDER=${KUBEVIRT_PROVIDER:-'k8s-1.34'}
export KUBEVIRT_NUM_NODES=${KUBEVIRT_NUM_NODES:-3}
export KUBEVIRT_NUM_SECONDARY_NICS=${KUBEVIRT_NUM_SECONDARY_NICS:-2}

# These can arrive with surrounding whitespace (e.g. from Makefile variables);
# strip everything but digits so arithmetic and `seq` behave.
KUBEVIRT_NUM_NODES=${KUBEVIRT_NUM_NODES//[^0-9]/}
KUBEVIRT_NUM_SECONDARY_NICS=${KUBEVIRT_NUM_SECONDARY_NICS//[^0-9]/}

export CLUSTER_NAME=${CLUSTER_NAME:-nmstate}
export KCLI_IMAGE=${KCLI_IMAGE:-centos9stream}
export KCLI_MACHINE=${KCLI_MACHINE:-pc}
export KCLI_NUMCPUS=${KCLI_NUMCPUS:-4}
export KCLI_MEMORY=${KCLI_MEMORY:-6144}

# Host-side image registry exposed to the nodes as registry:5000.
export REGISTRY_NAME=${REGISTRY_NAME:-nmstate-registry}
export REGISTRY_PORT=${REGISTRY_PORT:-5000}
export REGISTRY_IMAGE=${REGISTRY_IMAGE:-docker.io/library/registry:2}

KCLI_PATH="${PWD}/_kcli"
KCLI_PLAN="${PWD}/cluster/kcli/plan.yml"

# kcli client (hypervisor) to target. Defaults to the local libvirt host so the
# development cluster is provisioned on the developer machine regardless of any
# remote client configured as the kcli default in ~/.kcli/config.yml.
export KCLI_CLIENT=${KCLI_CLIENT:-local}

# Wrapper that pins every kcli invocation to ${KCLI_CLIENT}. Scripts call the
# bare `kcli` command and transparently get the right client.
function kcli() {
    command kcli -C "${KCLI_CLIENT}" "$@"
}

function kcli::path() {
    echo -n "${KCLI_PATH}/${CLUSTER_NAME}"
}

function kcli::kubeconfig() {
    echo -n "${KCLI_PATH}/${CLUSTER_NAME}/kubeconfig"
}

function kcli::plan() {
    echo -n "${KCLI_PLAN}"
}

# kubernetes minor version derived from KUBEVIRT_PROVIDER (e.g. k8s-1.34 -> 1.34)
function kcli::k8s_version() {
    echo -n "${KUBEVIRT_PROVIDER}" | sed -E 's/.*[^0-9]([0-9]+\.[0-9]+).*/\1/'
}

# Gateway IP of the libvirt network the nodes attach to. Reachable from the
# guests and used to advertise the host registry as registry:5000.
function kcli::gateway_ip() {
    local ip
    ip=$(virsh -c qemu:///system net-dumpxml default 2>/dev/null \
        | sed -n "s/.*<ip address='\([0-9.]*\)'.*/\1/p" | head -1)
    echo -n "${ip:-192.168.122.1}"
}

# Echo the node names node01..node0N.
function kcli::nodes() {
    local n
    for n in $(seq 1 "${KUBEVIRT_NUM_NODES}"); do
        printf 'node%02d\n' "${n}"
    done
}

function kcli::control_plane() {
    echo -n "node01"
}

# Primary IP address of a node as known to libvirt.
function kcli::node_ip() {
    local node=${1}
    kcli info vm "${node}" -f ip -v 2>/dev/null | tr -d '\r'
}

# Ensure the libvirt daemons are reachable (modular sockets are socket
# activated, so `systemctl is-active libvirtd` may report inactive).
function kcli::ensure_libvirt() {
    if virsh -c qemu:///system list >/dev/null 2>&1; then
        return 0
    fi
    sudo systemctl start virtqemud.socket virtnetworkd.socket virtstoraged.socket 2>/dev/null || \
        sudo systemctl start libvirtd 2>/dev/null || true
    virsh -c qemu:///system list >/dev/null 2>&1
}
