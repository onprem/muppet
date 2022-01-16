#!/usr/bin/env bash

trap 'kill 0' SIGTERM

SERVICE_CONTAINER=${SERVICE_CONTAINER:-"onprem/muppet-service"}
AGENT_CONTAINER=${AGENT_CONTAINER:-"onprem/command-agent"}

# Muppet Service

echo "starting the muppet service on http://localhost:8080"

docker run -p 8080:8080 ${SERVICE_CONTAINER} &

sleep 0.5

# Running 5 command agents

for i in $(seq 0 4); do
  echo "starting command agent with hostname: host00${i}"

  docker run \
    --network host \
    ${AGENT_CONTAINER} \
    --hostname host00"${i}" &
  
  sleep 2
done

echo "everything started; waiting for signal"

wait
