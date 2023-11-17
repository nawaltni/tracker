#!/bin/bash

export DEPLOY_KEY=$(cat ~/.ssh/id_rsa) 
docker-compose up $1 $@