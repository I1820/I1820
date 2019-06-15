#!/bin/bash
# In The Name Of God
# ========================================
# [] File Name : me.sh
#
# [] Creation Date : 28-09-2018
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
# Creates basis for pm component of I1820 platform.
# please run this script once and for all.
docker network create -d bridge --subnet 192.168.72.0/24 --gateway 192.168.72.1 i1820_projects
docker pull i1820/elrunner
docker pull redis:alpine
echo "please do not run this script more than once"
