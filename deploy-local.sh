#!/bin/bash

set -xe

sudo cp conf/isuports-app.service /etc/systemd/system && sudo systemctl daemon-reload

make -C go isuports

sudo systemctl restart nginx isuports-app
