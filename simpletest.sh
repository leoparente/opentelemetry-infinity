#!/usr/bin/env bash

CURRENT_DIR=$pwd
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

cd $SCRIPT_DIR

if [ ! -f "build/otlpinf" ]; then
    make build
fi

build/otlpinf run &
PID=$!

sleep 1

ret=$(curl -o /dev/null -s -w "%{http_code}\n" -X POST --location 'localhost:10222/api/v1/policies' -H 'Content-Type: application/x-yaml' --data 'pol_test:
  config:
    receivers:
      otlp:
        protocols:
          http:
          grpc: 
    exporters:
      logging:
        loglevel: debug
    service:
      pipelines:
        metrics:
          receivers:
          - otlp
          exporters:
          - logging')

kill $PID
cd $CURRENT_DIR

if [[ $ret != 201 ]]; then
  exit 1
fi

exit 0

