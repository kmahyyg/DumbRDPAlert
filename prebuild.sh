#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd ${SCRIPT_DIR}
VERSIONSTR=$(git describe --long --dirty --tags --always | tr -d '\n')

cat << EOF > ./gosrc/embedded/version.go
package embedded

var CurVersionStr = "${VERSIONSTR}"
EOF
