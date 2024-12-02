#!/bin/bash

set -xe

export GIT_REV=$(git rev-parse --short HEAD)

sudo cp conf/isuports-app.service /etc/systemd/system
sudo systemctl daemon-reload

envsubst '$GIT_REV' <conf/nginx.conf | sudo tee /etc/nginx/nginx.conf >/dev/null
sudo cp conf/nginx-isuports.conf /etc/nginx/sites-available/isuports.conf
envsubst '$GIT_REV' <conf/my.cnf | sudo tee /etc/mysql/mysql.conf.d/mysql.cnf >/dev/null

make -C go isuports dbclean

sudo systemctl restart nginx mysql isuports-app
