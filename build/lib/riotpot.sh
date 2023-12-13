#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Compile the riotpot application, statik, and plugins
function riotpot::compile(){
    # riotpot::compile::statik

    go build \
        -gcflags="all=-N -l" \
        -o "${RIOTPOT_BIN}/riotpot/" "${RIOTPOT_ROOT}/cmd/riotpot"
}

# This function compiles the plugins as <plugin_name>.so
function riotpot::compile::plugins(){

    for plugin in "$RIOTPOT_ROOT/pkg/plugin"/*/; do
        pg=$(basename "$plugin")
        echo "$plugin"
        echo "$pg"
        go build \
            -buildmode="plugin" \
            --mod="mod" \
            -gcflags="all=-N -l" \
            -o "${RIOTPOT_BIN}/riotpot/plugins/${pg}.so" \
            "$plugin"/*.go
    done
}

# This function places the statik files from the API into the application
function riotpot::compile::statik(){
    statik -src="api/swagger"
}