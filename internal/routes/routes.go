package routes

import (
	"github.com/doublehops/dh-go-framework/internal/handlers/health"
	"github.com/doublehops/dh-go-framework/internal/handlers/user"
	"github.com/doublehops/dh-go-framework/internal/middleware"
	"github.com/doublehops/dh-go-framework/internal/service"
	group "github.com/mythrnr/httprouter-group"

	"github.com/doublehops/dh-go-framework/internal/handlers/author"
	// "github.com/doublehops/dh-go-framework/internal/handlers/mynewtable"
)

func GetV1Routes(app *service.App) *group.RouteGroup {
	healthHandle := health.New(app)
	healthGroup := group.New("/health").GET(healthHandle.Check)

	authorHandle := author.New(app)
	authorGroup := group.New("/author")
	authorGroup.GET(authorHandle.GetAll).Middleware(middleware.AuthMiddleware)
	authorGroup.Children(
		group.New("/:id").GET(authorHandle.GetByID),
		group.New("").POST(authorHandle.Create),
		group.New("/:id").PUT(authorHandle.UpdateByID),
		group.New("/:id").DELETE(authorHandle.DeleteByID),
	)

	userHandle := user.New(app)

	userGroup := group.New("/user")
	userGroup.GET(userHandle.GetAll).Middleware(middleware.AuthMiddleware)
	userGroup.Children(
		group.New("/:id").GET(userHandle.GetByID),
		group.New("").POST(userHandle.Create),
		group.New("/:id").PUT(userHandle.UpdateByID),
		group.New("/:id").DELETE(userHandle.DeleteByID),
	)

	// New routes created by scaffolding can be added here.

	// myNewTableHandle := mynewtable.New(app)
	//
	// myNewTableGroup := group.New("/my-new-table")
	// myNewTableGroup.GET(myNewTableHandle.GetAll).Middleware(middleware.AuthMiddleware)
	// myNewTableGroup.Children(
	// 	group.New("/:id").GET(myNewTableHandle.GetByID),
	// 	group.New("").POST(myNewTableHandle.Create),
	// 	group.New("/:id").PUT(myNewTableHandle.UpdateByID),
	// 	group.New("/:id").DELETE(myNewTableHandle.DeleteByID),
	// )

	g := group.New("/v1").Children(
		healthGroup,
		authorGroup,
		userGroup,
		// myNewTableGroup,
		// Add new groups here.
	)

	return g
}
