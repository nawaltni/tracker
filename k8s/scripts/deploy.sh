#!/bin/bash
# Deployment script to Production Cluster
set -e

# Switch based on the environment parameter
case "$1" in
    development)
        NAMESPACE="development"
        # Other variables specific to development environment
        ;;
    staging)
        NAMESPACE="staging"
        # Other variables specific to staging environment
        ;;
    production)
        NAMESPACE="production"
        # Other variables specific to production environment
        ;;
    *)
        echo "Invalid environment parameter. Please specify development, staging, or production."
        exit 1
        ;;
esac


GCLOUD_PROJECT=nawalt
GCLOUD_CLUSTER=nawalt-1
GCLOUD_ZONE=us-east1-b

# provide a version as an argument
commit=$(git rev-parse --short HEAD)
tstamp=$(date +"%Y%m%d%H%M%S")
# combine vertion tag as commit tag and timestamp
version="$commit-$tstamp"

echo -e "\e[1;34mPreparing to Deploy to \e[33m$NAMESPACE\e[0m"
IMAGE=us-east1-docker.pkg.dev/nawalt/nawalt/tracker
IMAGE_NAME=$IMAGE:$version

echo -e "\e[1;34mPreparing image \e[33m$IMAGE_NAME\e[0m"



echo -e "\e[1;34mRunning from Computer\e[0m"
echo -e "\e[1;31mWarning: This should only be used in an emergency.\e[0m"
read -p "Are you sure you want to deploy from your computer?(y/n)" answer
case ${answer:0:1} in
    y|Y )
        DEPLOY_KEY=$(cat ~/.ssh/id_rsa)
    ;;
    * )
        echo -e "\e[1;32mDeployment stopped\e[0m";
        exit
    ;;
esac


gcloud config set project $GCLOUD_PROJECT
gcloud container clusters get-credentials $GCLOUD_CLUSTER --zone=$GCLOUD_ZONE --project $GCLOUD_PROJECT
gcloud auth configure-docker --quiet us-east1-docker.pkg.dev

# build a container
echo -e "\e[1;34mBuilding image\e[0m"
docker build -t $IMAGE_NAME --build-arg DEPLOY_KEY="$DEPLOY_KEY" .
if [ $? -eq 0 ]; then
	echo -e "\e[1;32mBuild complete\e[0m"
else
	echo -e "\e[1;31mBuild failed\e[0m"
	exit 1
fi

# push to GCR
echo -e "\e[1;34mPushing image to GCR\e[0m"
docker push $IMAGE_NAME

if [ $? -eq 0 ]; then
	echo -e "\e[1;32mImage pushed\e[0m"
else
	echo -e "\e[1;31mCould not push Image \e[0m"
	exit 1
fi

echo -e "\e[1;34mDeploying tracker\e[0m"

helm upgrade --install \
    --values ./k8s/helm/tracker/values.$1.yaml \
    --set deployment.image=$IMAGE \
    --set deployment.tag=$version \
    tracker ./k8s/helm/tracker \
    -n $1
echo -e "\e[1;32mtracker deployed\e[0m"

# check 
kubectl get pods -n $1