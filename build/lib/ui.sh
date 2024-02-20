#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

function riotpot::ui::build() {
    cd ui
    npm install
    npm run build
}

function riotpot::ui::serve_dev() {
    serve -s ./ui/build
}
