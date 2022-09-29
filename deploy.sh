#!/bin/sh

set -e
set -x

#go build

ps -ef \
| grep homeflix \
| grep -v grep \
| grep movies \
| awk '{ print $2 }' \
| xargs kill \
|| true

cp homeflix $HOME/bin/homeflix

cd $HOME/movies

./start_homeflix.sh
