#!/bin/bash
APP_NAME="public-ir-controller"
id=$(sudo docker restart $APP_NAME)
sudo docker logs -f $id