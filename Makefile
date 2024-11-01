PROTO_SRC=internal/lockserver/lockserver.proto
GO_OUT_DIR=.

# Commands
PROTOC=protoc
PROTOC_GEN_GO=$(shell which protoc-gen-go)
PROTOC_GEN_GO_GRPC=$(shell which protoc-gen-go-grpc)

# Check if required tools are installed
.PHONY: check_tools
check_tools:
ifndef PROTOC_GEN_GO
	$(error "protoc-gen-go not found. Please install it by running 'go install google.golang.org/protobuf/cmd/protoc-gen-go@latest'")
endif
ifndef PROTOC_GEN_GO_GRPC
	$(error "protoc-gen-go-grpc not found. Please install it by running 'go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest'")
endif

# Generate Go files from .proto
.PHONY: generate
generate: check_tools
	@mkdir -p $(GO_OUT_DIR)
	$(PROTOC) --go_out=$(GO_OUT_DIR) --go-grpc_out=$(GO_OUT_DIR) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $(PROTO_SRC)

# Clean generated files
.PHONY: clean
clean:
	@rm -rf $(GO_OUT_DIR)/*.pb.go

# Run all targets
.PHONY: all
all: generate

