#!/bin/bash

set -e

help() {
    cat <<- EOF
Usage: TAG=tag $0

Creates git tags for public Go packages.

VARIABLES:
  TAG        git tag, for example, v1.0.0
EOF
    exit 0
}

if [ -z "$TAG" ]
then
    printf "TAG env var is required\n\n";
    help
fi

if ! git grep -Fq "${TAG}" -- '*/go.mod'
then
    printf "no go.mod references ${TAG}; merge the release branch first\n"
    exit 1
fi

# Modules marked with a `// Deprecated` comment in go.mod are not tagged again.
PACKAGE_DIRS=$(find . -mindepth 2 -type f -name 'go.mod' -exec dirname {} \; \
  | sed 's/^\.\///' \
  | sort \
  | while read -r dir
    do
        grep -q '^// Deprecated' "${dir}/go.mod" || echo "${dir}"
    done)

git tag ${TAG}
git push origin ${TAG}

for dir in $PACKAGE_DIRS
do
    printf "tagging ${dir}/${TAG}\n"
    git tag ${dir}/${TAG}
    git push origin ${dir}/${TAG}
done
