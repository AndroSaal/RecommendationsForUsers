#ABOUT HOW CONNECT TO DATABASE
DB_DRIVE = postgres
DB_USER = postgres
DB_PASS = qwerty
DB_IP = localhost
DB_PORT = 5432
DB_NAME = postgres
DB_OPT = sslmode=disable

#MIGRATIONS
MIGRATIONSDIR = migration
MIGR_NAME = init

#CONFIG
CONFIGPATH = ./config
CONFIGNAME = local.yaml

run : clean build
	./.bin/analytics  --config_path=$(CONFIGPATH) --config_file=$(CONFIGNAME)
	
build :
	go build -o ./.bin/analytics ./internal/cmd/main.go

clean :
	rm -rf ./.bin

migr_create:
	mkdir -p $(MIGRATIONSDIR)
	migrate create -ext sql -dir $(MIGRATIONSDIR) -seq $(MIGR_NAME)

migr_up:
	migrate -path $(MIGRATIONSDIR) -database \
	"$(DB_DRIVE)://$(DB_USER):$(DB_PASS)@$(DB_IP):$(DB_PORT)/$(DB_NAME)?$(DB_OPT)" \
	-verbose up

migr_down:
	migrate -path $(MIGRATIONSDIR) -database \
	"$(DB_DRIVE)://$(DB_USER):$(DB_PASS)@$(DB_IP):$(DB_PORT)/$(DB_NAME)?$(DB_OPT)" \
	-verbose down

test:
	rm -rf coverage.out
	go test ./internal/service ./internal/repository \
	./internal/transport/kafka/consumer -v -cover -coverprofile=coverage.out

check_coverage:
	go tool cover -html=coverage.out