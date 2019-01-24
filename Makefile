test:
	@docker-compose up -d --scale test=0
	@docker-compose up --abort-on-container-exit test

clean:
	@docker-compose rm -sfv
