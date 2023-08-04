package service

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapplambda/goapp/compo"
)

func Entry() {
	compo.Routes()
	app.RunWhenOnBrowser()
}
