#! /bin/bash

if [ "$ACTION" = "start" ]; then
  SERVICES=$(docker compose config --services | grep -v "lytemp-dev")
  docker compose up -d $SERVICES
elif [ "$ACTION" = "stop" ]; then
  SERVICES=$(docker compose config --services | grep -v "lytemp-dev")
  docker compose down $SERVICES
elif [ "$ACTION" = "restart" ]; then
  SERVICES=$(docker compose config --services | grep -v "lytemp-dev")
  docker compose restart $SERVICES
fi
