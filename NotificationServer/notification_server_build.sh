#!/bin/bash


current_time=$(date "+%m월%d일 %H:%M")

echo "------ Start NotificationServer Build ------"
GOOS=linux GOARCH=amd64 go build -o notification_server
echo "$current_time Build Success"
echo "--------------------------------------------"