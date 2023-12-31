package compo

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var _ app.AppUpdater = (*Root)(nil)
var _ app.Mounter = (*Root)(nil)

type Root struct {
	app.Compo
}

func (r *Root) Render() app.UI {
	return app.Div().Body(
		app.P().Text("goapplambda demo application"),
		app.P().Body(
			app.A().Href("https://github.com/mlctrez/goapplambda").Text("goapplambda on github"),
		),
		app.P().Text("GOAPP_VERSION = "+app.Getenv("GOAPP_VERSION")),
		app.P().Text("DEV = "+app.Getenv("DEV")),
		&Headers{},
	)
}

func (r *Root) OnAppUpdate(ctx app.Context) {
	if app.Getenv("DEV") != "" && ctx.AppUpdateAvailable() {
		ctx.Reload()
	}
}

func (r *Root) OnMount(ctx app.Context) {

}
