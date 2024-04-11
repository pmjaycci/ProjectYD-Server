#!/bin/bash

REPOSITORY_NAME="game_server_8080-game_server"

echo "---------- Start GameServer Build ----------"
sudo docker-compose down
sudo docker rmi -f $(docker images | grep $REPOSITORY_NAME | awk '{print $3}')
sudo docker-compose up -d
echo "--------------------------------------------"