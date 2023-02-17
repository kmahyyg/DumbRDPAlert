#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd ${SCRIPT_DIR}
VERSIONSTR=$(git describe --long --dirty --tags --always | tr -d '\n')

cat << EOF > ./gosrc/cmd/version.go
package main

var CurVersionStr = "${VERSIONSTR}"
EOF