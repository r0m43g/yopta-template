migrate up:
	goose -dir backend/migrations mysql "${DSN}" up

build-front:
	cd frontend && npm run build && rm -rf ../backend/static/* && mv dist/* ../backend/static/

run:
	cd backend && go run cmd/server/main.go
