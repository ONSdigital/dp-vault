#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-vault
  make test
popd