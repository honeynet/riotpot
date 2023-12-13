#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Installs dependencies
function riotpot::requirements::install() {
    riotpot::log::status "Installing requirements..."

    local -a pkt=($(riotpot::requirements::get_requirements))

    riotpot::dependencies::with_apt ${pkt[@]}
    riotpot::golang::install || return 1

    riotpot::log::status "Finished installing requirements"
}

# Installs third party applications that we depend on
function riotpot::dependencies::with_apt() {
    local -a deps=()
    deps=($@)

    riotpot::log::status "Updating environment..."
    apt-get update -y

    riotpot::log::status "Installing third-party packages..."

    for dep in "${deps[@]}"; do
        if [[ -z $(which "$dep") ]]; then
            apt-get --no-install-recommends -y install "$dep"
        fi
    done
}

# Return the list of required packages
function riotpot::requirements::get_requirements() {
    local -a reqs=($(cat "${RIOTPOT_ROOT}/build/requirements.txt"))
    echo ${reqs[@]}
}
