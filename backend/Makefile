compose-restart:
	@sh scripts/restart.sh

local-run:
	docker compose up 	

hot:
	docker compose -p app -f docker-compose.yml up --build

deploy:
	ACTION=start ./scripts/deploy.sh

services-start:
	ACTION=start ./scripts/dev.sh

services-stop:
	ACTION=stop ./scripts/dev.sh

services-restart:
	ACTION=restart ./scripts/dev.sh
	
deploy-stage:
	APP_ENVIRONMENT=stage ACTION=start ./scripts/deploy.sh
