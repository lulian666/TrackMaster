#!/usr/bin/env bash
ENV=${ENV:-MASTER}

if [ $ENV == "PRODUCER" ]; then
    echo "start producer"
    exec /app/Producer
fi


if [ $ENV == "CONSUMER" ]; then
    echo "start consumer"
    exec /app/Consumer
fi

if [ $ENV == "MASTER" ]; then
    echo "start main service"
    exec /app/TrackMaster
fi
