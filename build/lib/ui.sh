#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

function riotpot::ui::build(){
    npm install ui --silent
    npm --prefix="ui" --omit="dev" run build 
}

function riotpot::ui::serve_dev(){
    serve -s ./ui/build
}