#!/usr/bin/env bash

CURRENT_DIR=$pwd
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

cd $SCRIPT_DIR

if [ ! -f "../build/otlpinf" ]; then
    make -C ../ build
fi

../build/otlpinf run &
PID=$!

sleep 1

ret=$(curl -o /dev/null -s -w "%{http_code}\n" -X POST -H "Content-Type: application/x-yaml" --data @policy.yaml http://localhost:10222/api/v1/policies)

echo $ret

kill $PID

cd $CURRENT_DIR
exit 0

