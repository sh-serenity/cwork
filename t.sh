#!/bin/sh

replace '#DOCKER_HOST_ADDRESS=192.168.1.1' DOCKER_HOST_ADDRESS=$1.1  -- .env 