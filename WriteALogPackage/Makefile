gen_proto:
	protoc internal/api/v1/*.proto \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--proto_path=.

.PHONY: gen_proto