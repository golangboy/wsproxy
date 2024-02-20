#!/bin/sh
if [ -z $ws ]; then
    echo "Starting server"
    /app/server/app
else
    echo "Starting client"
    /app/client/app -ws=$ws
fi