#!/bin/bash
# Clone this repo locally
sudo docker-compose down
sudo docker rm -f `sudo docker ps -aq`
sudo docker rmi -f `sudo docker images -q`
sudo docker swarm leave --force
sudo docker system prune -af
# Load the IPVS kernel module. Because swarms are created in dind,
# the daemon won't load it automatically
sudo modprobe xt_ipvs

# Ensure the Docker daemon is running in swarm mode
docker swarm init

# Get the latest shashwatpp/dind image
docker pull shashwatpp/dind

# Optional (with go1.14): pre-fetch module requirements into vendor
# so that no network requests are required within the containers.
# The module cache is retained in the pwd and l2 containers so the
# download is a one-off if you omit this step.
# go mod vendor

# Start PWD as a container
docker-compose up
