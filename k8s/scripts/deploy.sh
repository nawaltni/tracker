#!/bin/bash
# Deployment script to Production Cluster
set -e

ACCESS_TOKEN=$DO_PAT # expected to be set in the environment
if [[ -z "${ACCESS_TOKEN}" ]]; then
	echo -e "\e[1;31mDigitalOcean Access Token not found\e[0m"
	echo -e "\e[1;31mPlease authenticate doctl for use with your DigitalOcean account. You can generate a token in the control panel at https://cloud.digitalocean.com/account/api/tokens\e[0m"
	exit 1
fi

# Switch based on the environment parameter
case "$1" in
    development)
        NAMESPACE="dev"
        PROJECT=nawalt-dev
        K8S_CLUSTER=nawalt-dev
        K8S_ZONE=nyc1
        VALUES_FILE=values.yaml
        echo -e "\e[1;34mPreparing to Deploy to \e[33mDevelopment\e[0m"
        ;;
    staging)
        NAMESPACE="staging"
        PROJECT=nawalt-staging
        K8S_CLUSTER=nawalt-staging
        K8S_ZONE=nyc1
        VALUES_FILE=values.staging.yaml
        echo -e "\e[1;34mPreparing to Deploy to \e[33mStaging\e[0m"
        ;;
    production)
        NAMESPACE="production"
        PROJECT=nawalt
        K8S_CLUSTER=nawalt
        K8S_ZONE=nyc1
        VALUES_FILE=values.production.yaml
        echo -e "\e[1;34mPreparing to Deploy to \e[33mProduction\e[0m"
        ;;
    *)
        echo "Invalid environment parameter. Please specify development, staging, or production."
        exit 1
        ;;
esac


# provide a version as an argument
commit=$(git rev-parse --short HEAD)
tstamp=$(date +"%Y%m%d%H%M%S")
# combine vertion tag as commit tag and timestamp
version="$commit-$tstamp"
# if $1 NAMESPACE is not production, append the $1 to the version
if [ "$1" != "production" ]; then
    version="$version-$1"
fi

echo -e "\e[1;34mPreparing to Deploy to \e[33m$NAMESPACE\e[0m"
IMAGE=registry.digitalocean.com/nawalt/tracker
IMAGE_NAME=$IMAGE:$version

echo -e "\e[1;34mPreparing image \e[33m$IMAGE_NAME\e[0m"



echo -e "\e[1;34mRunning from Computer\e[0m"
echo -e "\e[1;31mWarning: This should only be used in an emergency.\e[0m"
read -p "Are you sure you want to deploy from your computer?(y/n)" answer
case ${answer:0:1} in
    y|Y )
        :
    ;;
    * )
        echo -e "\e[1;32mDeployment stopped\e[0m";
        exit
    ;;
esac


# Authenticate with DigitalOcean
echo -e "\e[1;34mAuthenticating with DigitalOcean\e[0m"
doctl kubernetes cluster kubeconfig save $K8S_CLUSTER -t $ACCESS_TOKEN

# Setting up Docker config
echo -e "\e[1;34mSetting up Docker config\e[0m"
doctl registry login -t $ACCESS_TOKEN --expiry-seconds 300

# build a container
echo -e "\e[1;34mBuilding image\e[0m"

build_command="docker build --ssh -t $IMAGE_NAME --build-arg DEPLOY_USER=\"$DEPLOY_USER\" ."
push_command="docker push $IMAGE_NAME"

# Check architecture and adjust build command if necessary
if [[ $(uname -m) == 'arm64' ]]; then
    build_command="docker buildx build --ssh default --platform linux/amd64 -t $IMAGE_NAME  --push ."
    push_command="" # No need to push separately for buildx as --push is used
fi

echo $build_command
# Execute build
eval $build_command
echo -e "\e[1;32mBuild complete\e[0m"

# Execute push if not using buildx (as buildx already includes push if --push is used)
if [[ -n $push_command ]]; then
    eval $push_command
fi


if [ $? -eq 0 ]; then
	echo -e "\e[1;32mImage pushed\e[0m"
else
	echo -e "\e[1;31mCould not push Image \e[0m"
	exit 1
fi

echo -e "\e[1;34mDeploying tracker\e[0m"

helm upgrade --install \
    --values k8s/helm/tracker/$VALUES_FILE \
    --set deployment.image=$IMAGE \
    --set deployment.tag=$version \
    tracker ./k8s/helm/tracker \
    -n $NAMESPACE
echo -e "\e[1;32mtracker deployed\e[0m"

# check 
kubectl get pods -n $NAMESPACE