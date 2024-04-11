#!/bin/bash

REPOSITORY_NAME="notification_server_8082-notification_server"

echo "------ Start NotificationServer Build ------"
sudo docker-compose down
sudo docker rmi -f $(docker images | grep $REPOSITORY_NAME | awk '{print $3}')
sudo docker-compose up -d
echo "--------------------------------------------"