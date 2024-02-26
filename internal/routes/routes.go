package routes

import (
	"github.com/doublehops/dhapi-example/internal/middleware"
	"github.com/doublehops/dhapi-example/internal/service"
	group "github.com/mythrnr/httprouter-group"

	"github.com/doublehops/dhapi-example/internal/handlers/author"
	"github.com/doublehops/dhapi-example/internal/handlers/mynewtable"
)

func GetV1Routes(app *service.App) *group.RouteGroup {
	authorHandle := author.New(app)

	authorGroup := group.New("/author")
	authorGroup.GET(authorHandle.GetAll).Middleware(middleware.AuthMiddleware)
	authorGroup.Children(
		group.New("/:id").GET(authorHandle.GetByID),
		group.New("").POST(authorHandle.Create),
		group.New("/:id").PUT(authorHandle.UpdateByID),
		group.New("/:id").DELETE(authorHandle.DeleteByID),
	)

	// New routes created by scaffolding can be added here.

	myNewTableHandle := mynewtable.New(app)

	myNewTableGroup := group.New("/my-new-table")
	myNewTableGroup.GET(myNewTableHandle.GetAll).Middleware(middleware.AuthMiddleware)
	myNewTableGroup.Children(
		group.New("/:id").GET(myNewTableHandle.GetByID),
		group.New("").POST(myNewTableHandle.Create),
		group.New("/:id").PUT(myNewTableHandle.UpdateByID),
		group.New("/:id").DELETE(myNewTableHandle.DeleteByID),
	)

	g := group.New("/v1").Children(
		authorGroup,
		myNewTableGroup,
	)

	return g
}
