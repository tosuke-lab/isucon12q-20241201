#!/bin/bash

set -xe

sudo cp conf/isuports-app.service /etc/systemd/system && sudo systemctl daemon-reload
sudo cp conf/nginx.conf /etc/nginx/nginx.conf
sudo cp conf/nginx-isuports.conf /etc/nginx/sites-available/isuports.conf

make -C go isuports

sudo systemctl restart nginx isuports-app
