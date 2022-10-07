#!/bin/bash

docker network create -d bridge redisnet
docker run -d \
  -h redis \
  -e REDIS_PASSWORD=**CHANGE_ME** \  
  -p 6379:6379 \
  --name redis \
  --restart on-failure \
  --network redisnet \
  redis:latest /bin/sh -c 'redis-server --appendonly yes --requirepass ${REDIS_PASSWORD}'

docker run --env-file .env --network redisnet xefiji/lbc:latest
