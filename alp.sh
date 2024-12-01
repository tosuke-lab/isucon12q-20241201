#!/bin/bash

set -xe

head_rev=$(git rev-parse HEAD)
rev=${1:-$head_rev}
ssh isucon@isu01.i12q.tosuke.dev "cd webapp && ./alp-local.sh $rev"
