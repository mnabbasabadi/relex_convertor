#!/bin/bash

echo "Generating Server API"
go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.0 --config=server.config.yaml relex.openapi3.yaml
