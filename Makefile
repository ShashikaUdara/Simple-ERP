.PHONY: build-mysql build-app run

all: help

build-mysql:
	docker build -t erp-db -f Dockerfile.mysql .

build-app:
	docker build -t erp-app .

run:
	docker-compose up

stop:
	docker-compose down

import:
	@./nodb/import-db.sh

drop:
	@./nodb/drop-db.sh

db-user:
	@./nodb/create-db-user.sh

prep:
	sudo systemctl stop mysql

.PHONY: help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Infrastructure Targets:"
	@echo "  make build-app ---------- Build app image."
	@echo "  make build-mysql -------- Build mysql image."
	@echo ""
	@echo "Application Targets:"
	@echo "  make run ---------------- Run the app."
	@echo "  make stop --------------- Stop the app."
	@echo ""
	@echo "Database Targets:"
	@echo "  make import ------------- Create erp db."
	@echo "  make drop --------------- Drop the db."
	@echo ""
	@echo "General Targets:"
	@echo "  make -------------------- Display this help message."