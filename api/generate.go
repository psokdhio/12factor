//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -generate types,chi-server,strict-server,spec,client -package api -o openapi.gen.go openapi.yaml

package api
