#!/bin/bash -ilex

set -xe

DOCKER_BASE_REPO=$1
DOCKER_BUILD_TAG=$2
APP_NAME=$3
module="./App.go"

######################
# build package #
######################
echo "delete old App"
# rm ./App

echo "build ${module}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o App $module

echo "upx App"
upx App
#


#####################
# build docker image #
#####################

docker build --build-arg exePath=App -t ${DOCKER_BASE_REPO}/${APP_NAME}:${DOCKER_BUILD_TAG} -f deploy/Dockerfile .
docker tag ${DOCKER_BASE_REPO}/${APP_NAME}:${DOCKER_BUILD_TAG} ${DOCKER_BASE_REPO}/${APP_NAME}:latest
docker push ${DOCKER_BASE_REPO}/${APP_NAME}:${DOCKER_BUILD_TAG}
docker push ${DOCKER_BASE_REPO}/${APP_NAME}:latest

# sh ./deploy/docker-build.sh ${docker_resp} 0.0.1 ${APP_NAME}

# docker stop ${APP_NAME} && docker rm ${APP_NAME} && docker pull ${docker_resp}/${APP_NAME}:latest
# docker restart ${APP_NAME}

# docker run -d --name ${APP_NAME} --restart=always -m 512M --memory-swap -1 -p 8080:8080 ${docker_resp}/${APP_NAME}:latest


rm ./App
