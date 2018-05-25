#!/bin/bash
# In The Name Of God
# ========================================
# [] File Name : cleanup.sh
#
# [] Creation Date : 25-05-2018
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
for name in $(curl -s "127.0.0.1:8080/api/project" | jq -r '.[].name'); do
	echo $name
	curl -X DELETE -o /dev/null -w "%{http_code}" -s "127.0.0.1:8080/api/project/$name"
done
