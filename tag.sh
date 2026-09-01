#!/bin/bash

set -e

help() {
    cat <<- EOF
Usage: tag.sh {module} {tag}

Creates a git tag for one Go module.

ARGUMENTS:
  module     module directory name, for example, buffer
  TAG        git tag, for example, v1.0.0
EOF
    exit 0
}

MODULE="$1"
TAG="$2"

if [ -z "$TAG" ]
then
    printf "TAG is required\n\n"
    help
fi

if [ -z "$MODULE" ]
then
    printf "module argument is required\n\n"
    help
fi

TAG_REGEX="^v(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)(\\-[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?(\\+[0-9A-Za-z-]+(\\.[0-9A-Za-z-]+)*)?$"
if ! [[ "${TAG}" =~ ${TAG_REGEX} ]]; then
    printf "TAG is not valid: ${TAG}\n\n"
    exit 1
fi

if ! [[ "${MODULE}" =~ ^[A-Za-z0-9][A-Za-z0-9_-]*$ ]]; then
    printf "module is not valid: ${MODULE}\n\n"
    exit 1
fi

if ! grep -Fq "module go.x2ox.com/sorbifolia/${MODULE}" "${MODULE}/go.mod" 2>/dev/null
then
    printf "${MODULE} is not a module of this repo\n\n"
    exit 1
fi

TAG_FOUND=`git tag --list "${MODULE}/${TAG}"`
if [[ ${TAG_FOUND} = "${MODULE}/${TAG}" ]] ; then
    printf "tag ${MODULE}/${TAG} already exists\n\n"
    exit 1
fi

if grep -q '^// Deprecated' "${MODULE}/go.mod"
then
    printf "warning: ${MODULE} is marked deprecated in go.mod\n"
fi

printf "tagging ${MODULE}/${TAG}\n"
git tag "${MODULE}/${TAG}"
