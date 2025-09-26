#! /bin/bash
set -ex

# Hangi app servisini ayağa kaldıracağımıza env'e göre karar ver
APP_ENVIRONMENT="${APP_ENVIRONMENT:-prod}"

if [ "$APP_ENVIRONMENT" = "stage" ]; then
  APP_SERVICE="lytemp-stage"
else
  APP_SERVICE="lytemp-green"
fi

# Uygulama dışındaki servisleri ayağa kaldır (app servislerini filtre dışı bırak)
SERVICES=$(docker compose config --services \
  | grep -v "lytemp-dev" \
  | grep -v "lytemp-green" \
  | grep -v "lytemp-blue" \
  | grep -v "lytemp-stage")

docker compose up -d $SERVICES

# App servisini build + up
docker compose up -d --build "$APP_SERVICE"
