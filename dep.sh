#!/bin/sh

: ${GOPATH:=${HOME}/go}

SRCDIR="${GOPATH}/src/github.com/ProdriveTechnologies/distfile-mirror"
rm -rf "${SRCDIR}"
mkdir -p "$(dirname "${SRCDIR}")"
ln -sf "$(pwd)" "${SRCDIR}"
(cd "${SRCDIR}"; dep "$@")
bazel run //:gazelle
