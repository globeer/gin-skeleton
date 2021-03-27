package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/hyperjiang/gin-skeleton/config"
	"github.com/hyperjiang/gin-skeleton/router"
)

func main() {

	addr := flag.String("addr", config.Server.Addr, "Address to listen and serve")
	flag.Parse()

	gin.SetMode(config.Server.Mode)
	if config.Server.Mode == gin.ReleaseMode {
		gin.DisableConsoleColor()
	}
	// log to file
	logfile, _ := os.Create("app/log/request.log")
	gin.DefaultWriter = io.MultiWriter(logfile)
	errlogfile, _ := os.Create("app/log/error.log")
	gin.DefaultErrorWriter = io.MultiWriter(errlogfile)

	app := gin.Default()

	app.Static("/images", filepath.Join(config.Server.StaticDir, "images"))
	app.StaticFile("/favicon.ico", filepath.Join(config.Server.StaticDir, "images/favicon.ico"))
	app.LoadHTMLGlob(config.Server.ViewDir + "/*")
	app.MaxMultipartMemory = config.Server.MaxMultipartMemory << 20

	router.Route(app)

	// Listen and Serve
	app.Run(*addr)
}
