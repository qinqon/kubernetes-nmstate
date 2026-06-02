#!/bin/bash
#
# Join a worker node to the kubeadm control-plane.
# Runs on node02..node0N over `kcli ssh`.
#
# Usage: join-worker.sh <join-command...>
#   The join command is produced on node01 with:
#     kubeadm token create --print-join-command
#
set -ex

if [[ $# -eq 0 ]]; then
    echo "ERROR: missing kubeadm join command" >&2
    exit 1
fi

sudo "$@" --node-name="$(hostname)"

echo "### Joined $(hostname) to the cluster"
