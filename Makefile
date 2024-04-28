generate-rps:
	@protoc --go_out=./pkg  --go_opt=paths=source_relative     --go-grpc_out=./pkg --go-grpc_opt=paths=source_relative    api/notes_api.proto

migrate-test:
	@go run github.com/pressly/goose/v3/cmd/goose@latest --dir="db/migrations" postgres "postgres://postgres:Test_1tesT@127.0.0.1:5432/happy?sslmode=disable" reset
	@go run github.com/pressly/goose/v3/cmd/goose@latest --dir="db/migrations" postgres "postgres://postgres:Test_1tesT@127.0.0.1:5432/happy?sslmode=disable" up

migrate-run:
	@go run github.com/pressly/goose/v3/cmd/goose@latest --dir="db/migrations" postgres "postgres://postgres:Test_1tesT@127.0.0.1:5432/happy?sslmode=disable" up

migrate-drop:
	@go run github.com/pressly/goose/v3/cmd/goose@latest -dir db/migrations postgres "postgres://postgres:Test_1tesT@127.0.0.1:5432/happy?sslmode=disable" down-to 20240218144349
	
#   go run github.com/pressly/goose/v3/cmd/goose@latest -dir db/migrations/ create new_kek_table sql
# export PATH=$PATH:/usr/local/go/bin
# export PATH=$PATH:$HOME/go/bin
