//go:build !wasm

package service

import (
	"encoding/json"
	"fmt"
	"github.com/abihf/delta/v2"
	"github.com/gin-gonic/gin"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapplambda/goapp"
	"github.com/mlctrez/goapplambda/goapp/compo"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Entry() {
	compo.Routes()

	engine, err := buildGinEngine()
	if err != nil {
		log.Fatal(err)
	}
	err = delta.ServeOrStartLambda(":8080", engine, delta.WithLambdaURL())
	if err != nil {
		log.Fatal(err)
	}

}

var DevEnv = os.Getenv("DEV")
var IsDev = DevEnv != ""

type engineSetup func(*gin.Engine) error

func jsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency.String()

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}
func buildGinEngine() (engine *gin.Engine, err error) {

	if !IsDev {
		gin.SetMode(gin.ReleaseMode)
	}

	engine = gin.New()

	// required for go-app to work correctly
	engine.RedirectTrailingSlash = false

	if IsDev {
		// omit some common paths to reduce startup logging noise
		skipPaths := []string{
			"/app.css", "/app.js", "/app-worker.js", "/manifest.webmanifest", "/wasm_exec.js", "/web/app.wasm"}
		engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: skipPaths}))
	} else {
		engine.Use(jsonLoggerMiddleware())
	}
	engine.Use(gin.Recovery())

	engine.GET("/app.css", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/web/app.css")
	})

	// setupStaticHandlers
	setups := []engineSetup{setupApiEndpoints, setupGoAppHandler}

	for _, setup := range setups {
		if err = setup(engine); err != nil {
			return nil, err
		}
	}

	return
}

func setupApiEndpoints(engine *gin.Engine) error {
	// setup other api endpoints here

	engine.GET("/api/headers", func(context *gin.Context) {
		m := map[string]string{}
		for k, v := range context.Request.Header {
			if strings.ToLower(k) == "host" || strings.ToLower(k) == "via" {
				continue
			}
			m[k] = v[0]
		}
		context.JSON(http.StatusOK, m)
	})

	return nil
}

func setupGoAppHandler(engine *gin.Engine) (err error) {

	handler := &app.Handler{
		Name:                    "goapplambda",
		ShortName:               "goapplambda",
		BackgroundColor:         "#222",
		ThemeColor:              "#000",
		Styles:                  []string{"/web/style.css"},
		Title:                   "go-app lambda demo",
		Description:             "demonstrates deployment of go-app on aws lambda and s3",
		Author:                  "mlctrez@gmail.com",
		Keywords:                []string{"go-app", "lambda", "aws"},
		Env:                     app.Environment{},
		WasmContentLengthHeader: "Wasm-Content-Length",
	}

	handler.Env["DEV"] = os.Getenv("DEV")

	if IsDev {
		handler.AutoUpdateInterval = time.Second * 3
		handler.Version = ""
	} else {
		handler.AutoUpdateInterval = time.Hour
		handler.Version = fmt.Sprintf("%s@%s", goapp.Version, goapp.Commit)
	}

	goAppHandler := gin.WrapH(handler)

	goAppUrls := []string{
		"/", "/web/:path", "/app.js", "/app-worker.js", "/manifest.webmanifest", "/wasm_exec.js",
	}
	for _, url := range goAppUrls {
		engine.GET(url, goAppHandler)
	}

	return nil
}
