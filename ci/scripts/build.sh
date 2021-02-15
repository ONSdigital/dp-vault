#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-vault
  make build
popd