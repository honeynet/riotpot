
#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

readonly GO_VERSION="1.18"
readonly RIOTPOT_GOPATH="${RIOTPOT_GOPATH:-"/usr/bin/go"}"
export RIOTPOT_GOPATH

# This exports environment variables necessary to install golang apps
function riotpot::golang::setup_env() {
    export GOPATH="${RIOTPOT_GOPATH}"
    export PATH=$PATH:$GOPATH/bin
    export GOCACHE="${GOCACHE:-"${RIOTPOT_GOPATH}/cache/build"}"
    export GOMODCACHE="${GOMODCACHE:-"${RIOTPOT_GOPATH}/cache/mod"}"

    unset GOBIN
}

# This function ensures that some version is installed
#
# $1: minimum golang version expected
function riotpot::golang::ensure_installed() {
    # If it is NOT installed, return 1
    local -r gopath=$(which go)
    if [[ -z $gopath ]]; then
        riotpot::log::error "Could not find go in PATH"
        return 1
    fi

    # Check if the version was provided
    [[ ${#1} -gt 0 ]] || return
    
    local -r version=$1
    # If the version is NOT the same, return 1
    local -r v=$(go version 2>&1 | awk '{print $3}')
    if [[ "$v" < "$version" ]]; then
        riotpot::log::error "Found go version '$v', but expected '$version'"
        return 1
    fi
}

function riotpot::golang::install() {
    local -r arch="linux-amd64"
    local -r version="go${GO_VERSION}"
    local -r file_name="${version}.${arch}.tar.gz"

    # If go is already installed in the right version, return
    riotpot::golang::ensure_installed $version && return

    riotpot::log::status "Downloading and installing golang..."

    curl  -OL "https://go.dev/dl/$file_name"

    # From go.dev/doc/install
    rm -rf /usr/bin/go && tar -C /usr/bin -xf "${file_name}" && rm -rf "${file_name}"
    source $HOME/.profile
}

function riotpot::golang::install_requirements() {
    # statik
    go get github.com/rakyll/statik
}