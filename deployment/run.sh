#!/bin/bash

sudo docker network create -d bridge redisnet
sudo docker run --name=redis --publish=6379:6379 --hostname=redis --restart=on-failure --detach --network redisnet redis:latest
sudo docker run --env-file .env --network redisnet xefiji/lbc:latest
