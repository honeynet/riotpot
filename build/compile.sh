#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

RIOTPOT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..

source "${RIOTPOT_ROOT}/build/common.sh"
source "${RIOTPOT_ROOT}/build/lib/riotpot.sh"

riotpot::compile
riotpot::compile::plugins