//go:generate goagen bootstrap -d github.com/goadesign/goa/design/apidsl

package main

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"goa/controllers"

	//"goa/controllers"
	"goa/utils/database"
	"os"
)

func main() {
	// Create service
	service := goa.New("test_project")

	db, err := database.Connect(os.Getenv("PGUSER"), os.Getenv("PGPASS"), os.Getenv("PGDB"), os.Getenv("PGHOST"), os.Getenv("PGPORT"))
	if err != nil {
		service.LogError("db error", "err", err)
	}

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	_ = controllers.NewAuthenticationController(service, db) //TODO fix this

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}
}
