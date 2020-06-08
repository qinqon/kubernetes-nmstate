#!/bin/bash

set -xe

old_version=$(hack/versions.sh -2)
new_version=$(hack/versions.sh -1)
gh_organization=nmstate
gh_repo=kubernetes-nmstate

function upload() {
    resource=$1
    $GITHUB_RELEASE upload \
        -u $gh_organization \
        -r $gh_repo \
        --name $(basename $resource) \
	    --tag $new_version \
		--file $resource
}

function create_github_release() {
    # Create the release
    $GITHUB_RELEASE release \
        -u $gh_organization \
        -r $gh_repo \
        --tag $new_version \
        --name $new_version \
        --description "$(hack/render-release-notes.sh $old_version $new_version)"


    # Upload operator CRDs
    for manifest in $(ls deploy/crds/nmstate.io_*nmstate*); do
        upload $manifest
    done

    # Upload operator manifests
    for manifest in $(find $MANIFESTS_DIR -type f); do
        upload $manifest
    done
}

make OPERATOR_IMAGE_TAG=$new_version HANDLER_IMAGE_TAG=$new_version \
    manifests \
    push-handler \
    push-operator

create_github_release
