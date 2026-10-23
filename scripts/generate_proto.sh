#!/bin/bash
# Regenerates proto/listentogether.pb.go from proto/listentogether.proto.
# Needs protoc and protoc-gen-go (go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11).
set -e
cd "$(dirname "${BASH_SOURCE[0]}")/.."
protoc --go_out=proto --go_opt=paths=source_relative -Iproto proto/listentogether.proto
echo "Generated proto/listentogether.pb.go"
