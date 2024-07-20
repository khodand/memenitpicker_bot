#!/bin/bash

cd "$(dirname "$0")" || exit

docker-compose down
docker-compose pull

mkdir -p logs
touch bot.log  # Creates an empty bot.log if it doesn't exist

cp bot.log "logs/meme_bot_$(date +%F).log"
docker-compose up -d
nohup docker-compose logs -f > bot.log 2>&1 &
