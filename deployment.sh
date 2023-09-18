#!/bin/bash

# Check if the Docker image tag argument is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <docker_image_tag>"
    exit 1
fi

repo_url=$1

containers=$(docker ps | awk '{print $1}')
images=$(docker ps | awk '{print $2}')

# Pull Docker images
for image in $images; do
    docker pull $repo_url
done

# Stop and remove Docker containers
for container in $containers; do
    docker kill $container
h fdone

# Start Docker containers
docker run --env-file=.env $repo_url 