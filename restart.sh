#!/bin/sh

PROCESS_NAME="service_item_go"

SHELL_FOLDER=$(cd `dirname $0`; pwd)

ps aux | grep "${PROCESS_NAME}" | grep -v "grep"
COUNT=`ps -ef | grep "${PROCESS_NAME}" | grep -v "grep" | wc -l`
echo ""

if [ 0 == $COUNT ]; then
  echo "${PROCESS_NAME} starting..."
  nohup "${SHELL_FOLDER}/${PROCESS_NAME}" &
  echo "${PROCESS_NAME} started"
else
  echo "${PROCESS_NAME} restarting..."
  kill -1 $(ps -ef | grep "${PROCESS_NAME}" | grep -v "grep" | awk '{print $2}')
  echo "${PROCESS_NAME} restarted"
fi

sleep 1
ps aux | grep "${PROCESS_NAME}" | grep -v "grep"