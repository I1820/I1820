#!/bin/bash
# In The Name Of God
# ========================================
# [] File Name : $file.name
#
# [] Creation Date : $time.strftime("%d-%m-%Y")
#
# [] Created By : $user.name ($user.email)
# =======================================
echo "Remove dockers of existing projects"
docker rm -f `docker ps --format '{{.Names}}' | grep el_`

echo "Remove redises of existing projects"
docker rm -f `docker ps --format '{{.Names}}' | grep rd_`

echo "Create ISRC Network"
docker network create isrc
