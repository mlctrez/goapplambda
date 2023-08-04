//go:build !wasm

package service

import (
	"bytes"
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
	"sort"
	"time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("helloHandler")
	fmt.Println(r.URL.Path)
	_, _ = w.Write([]byte("hello world!"))
}

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
			"/app.css", "/app.js", "/app-worker.js", "/manifest.webmanifest", "/wasm_exec.js",
			"/web/logo-192.png", "/web/logo-512.png", "/web/logo.svg", "/web/app.wasm"}
		engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: skipPaths}))
	}
	engine.Use(gin.Recovery())

	engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Wasm-Content-Length", goapp.WasmSize)
		c.Next()
	})

	engine.Use(func(c *gin.Context) {
		m := map[string]string{}
		for k, v := range c.Request.Header {
			m[k] = v[0]
		}
		c.Next()
	})
	// https://www.google.com/maps/@38.74720,-90.72470,15z?entry=ttu
	// "Cloudfront-Viewer-Latitude": "38.74720",	//    "Cloudfront-Viewer-Longitude": "-90.72470"

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

		var keys []string
		for k := range context.Request.Header {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var buf bytes.Buffer
		buf.WriteString("{")
		for i, key := range keys {
			if i != 0 {
				buf.WriteString(",")
			}
			k, _ := json.Marshal(key)
			buf.Write(k)
			buf.WriteString(":")
			val, _ := json.Marshal(context.Request.Header[key][0])
			buf.Write(val)
		}

		buf.WriteString("}")

		context.Header("Content-Type", "application/json")
		context.String(http.StatusOK, buf.String())
	})

	return nil
}

func setupGoAppHandler(engine *gin.Engine) (err error) {

	handler := &app.Handler{
		Name:      "goapplambda",
		ShortName: "goapplambda",
		Icon: app.Icon{
			Default:    "/web/logo-192.png",
			Large:      "/web/logo-512.png",
			SVG:        "/web/logo.svg",
			AppleTouch: "/web/logo-192.png",
		},
		BackgroundColor:         "#222",
		ThemeColor:              "#000",
		Styles:                  []string{"/web/style.css"},
		Title:                   "go-app lambda demo",
		Description:             "demonstrates deployment of go-app on aws lambda url and s3",
		Author:                  "mlctrez@gmail.com",
		Keywords:                []string{"go-app", "lambda", "aws"},
		HTML:                    nil,
		Body:                    nil,
		AutoUpdateInterval:      0,
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

	engine.NoRoute(gin.WrapH(handler))
	return nil
}
