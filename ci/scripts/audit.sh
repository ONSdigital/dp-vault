#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-vault
  make audit
popd