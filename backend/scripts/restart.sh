docker compose stop app
docker rm -f app
docker compose build --no-cache app
docker compose up -d app
