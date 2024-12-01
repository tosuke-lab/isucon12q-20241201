#!/bin/bash

set -xe

export GIT_REV=$(git rev-parse --short HEAD)

sudo cp conf/isuports-app.service /etc/systemd/system
sudo systemctl daemon-reload

envsubst '$GIT_REV' <conf/nginx.conf | sudo tee /etc/nginx/nginx.conf >/dev/null
sudo cp conf/nginx-isuports.conf /etc/nginx/sites-available/isuports.conf
envsubst '$GIT_REV' <conf/my.conf | sudo tee /etc/mysql/mysql.conf.d/isuports.conf >/dev/null

make -C go isuports

sudo systemctl restart nginx isuports-app
