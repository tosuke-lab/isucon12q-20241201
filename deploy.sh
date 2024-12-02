#!/bin/bash

set -xe

ssh isucon@isu01.i12q.tosuke.dev "cd webapp && git pull && ./deploy-local.sh" 2>&1 | sed -E "s/^/[isu01] /" &
ssh isucon@isu03.i12q.tosuke.dev "cd webapp && git pull && ./deploy-local.sh" 2>&1 | sed -E "s/^/[isu03] /" &

wait
