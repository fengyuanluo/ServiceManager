// @title Example Service API
// @version 0.1.0
// @description Example Service 的 API 文档
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
//
//go:generate go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/swagger.go -o ../docs -d ../ --parseInternal --parseDependency
package cmd
