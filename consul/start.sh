#!/bin/bash

nohup consul agent -server -bootstrap-expect=1 -ui -data-dir=./ -node=agent-one -bind=0.0.0.0 -enable-script-checks=true > server.log 2>&1 &