#!/bin/bash

REPOSITORY_NAME="login_server_8081-login_server"

echo "---------- Start LoginServer Build ----------"
sudo docker-compose down
sudo docker rmi -f $(docker images | grep $REPOSITORY_NAME | awk '{print $3}')
sudo docker-compose up -d
echo "---------------------------------------------"