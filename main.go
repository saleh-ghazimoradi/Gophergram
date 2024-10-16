package main

import (
	"github.com/saleh-ghazimoradi/Gophergram/cmd"
)

// @title Gophergram
// @description API for Gophergram, a social platform for gophers
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description
func main() {
	cmd.Execute()
}
