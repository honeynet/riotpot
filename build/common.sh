#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Every script that loads this one will be set in the same path
RIOTPOT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd -P)

source "${RIOTPOT_ROOT}/build/log.sh"

readonly RIOTPOT_BIN="${OUT_DIR:-bin}"
readonly RIOTPOT_PLUGINS="${RIOTPOT_ROOT}/plugin"