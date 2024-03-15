#!/bin/bash

if [ -f ~/.ssh/id_rsa ]; then
    export DEPLOY_KEY=$(cat ~/.ssh/id_rsa)
else
    export DEPLOY_KEY=$(cat ~/.ssh/id_ed25519.pub)
fi

docker-compose up $1 $@