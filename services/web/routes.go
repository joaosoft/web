package web

import (
	"net/http"

	"github.com/joaosoft/manager"
)

func (controller *Controller) RegisterRoutes(web manager.IWeb) error {
	return web.AddRoutes(
		manager.NewRoute(http.MethodGet, "/api/v1/migrations/:id", controller.GetMigrationHandler),
		manager.NewRoute(http.MethodGet, "/api/v1/migrations", controller.GetMigrationsHandler),
		manager.NewRoute(http.MethodPost, "/api/v1/migrations", controller.CreateMigrationHandler),
		manager.NewRoute(http.MethodDelete, "/api/v1/migrations/:id", controller.DeleteMigrationHandler),
		manager.NewRoute(http.MethodDelete, "/api/v1/migrations", controller.DeleteMigrationsHandler),
	)
}
