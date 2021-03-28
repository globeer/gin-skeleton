package main

import (
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/globeer/gin-skeleton/config"
	"github.com/globeer/gin-skeleton/router"
	// "github.com/kjk/dailyrotate"
)

// func onLogClose(path string, didRotate bool) {
// 	fmt.Printf("we just closed a file '%s', didRotate: %v\n", path, didRotate)
// 	if !didRotate {
// 		return
// 	}
// 	// process just closed file e.g. upload to backblaze storage for backup
// 	go func() {
// 		// if processing takes a long time, do it in background
// 	}()
// }

// var (
// 	logFile *dailyrotate.File
// )

// func openLogFile(pathFormat string, onClose func(string, bool)) error {
// 	w, err := dailyrotate.NewFile(pathFormat, onLogClose)
// 	if err != nil {
// 		return err
// 	}
// 	logFile = w
// 	return nil
// }

// func closeLogFile() error {
// 	return logFile.Close()
// }

// func writeToLog(msg string) error {
// 	_, err := logFile.Write([]byte(msg))
// 	return err
// }

func main() {

	logDir := "app/log"

	// we have to ensure the directory we want to write to
	// already exists
	// err := os.MkdirAll(logDir, 0755)
	// if err != nil {
	// 	log.Fatalf("os.MkdirAll()(")
	// }
	// // only for the purpose of the demo, cleanup the directory
	// defer os.RemoveAll(logDir)

	// pathFormat := filepath.Join(logDir, "2006-01-02.txt")
	// err = openLogFile(pathFormat, onLogClose)
	// if err != nil {
	// 	log.Fatalf("openLogFile failed with '%s'\n", err)
	// }
	// defer closeLogFile()

	// err = writeToLog("hello\n")
	// if err != nil {
	// 	log.Fatalf("writeToLog() failed with '%s'\n", err)
	// }

	// // this is the actual path of the log file
	// // we're currently writing to
	// path := logFile.Path()

	// err = closeLogFile()
	// if err != nil {
	// 	log.Fatalf("closeLogFile() failed with '%s'\n", err)
	// }

	// st, err := os.Stat(path)
	// if err != nil {
	// 	log.Fatalf("os.Stat('%s') failed with '%s'\n", path, err)
	// }
	// fmt.Printf("We wrote %d bytes to log file %s\n", st.Size(), path)

	addr := flag.String("addr", config.Server.Addr, "Address to listen and serve")
	flag.Parse()

	gin.SetMode(config.Server.Mode)
	// if config.Server.Mode == gin.ReleaseMode {
	// 	gin.DisableConsoleColor()
	// }
	gin.DisableConsoleColor()

	// t := time.Now()
	// startTime := t.Format("2006-01-02 15:04:05")

	logFile, err := os.OpenFile(logDir+"/request.log", os.O_APPEND, 0666)
	if err != nil {
		logFile, _ = os.Create(logDir + "/request.log")
	}
	gin.DefaultWriter = io.MultiWriter(logFile)
	errlogfile, err := os.OpenFile(logDir+"/error.log", os.O_APPEND, 0666)
	if err != nil {
		errlogfile, _ = os.Create(logDir + "/error.log")
	}
	gin.DefaultErrorWriter = io.MultiWriter(errlogfile)

	// dailylog.initRotatedFileMust()

	app := gin.Default()

	app.Static("/images", filepath.Join(config.Server.StaticDir, "images"))
	app.StaticFile("/favicon.ico", filepath.Join(config.Server.StaticDir, "images/favicon.ico"))
	app.LoadHTMLGlob(config.Server.ViewDir + "/*")
	app.MaxMultipartMemory = config.Server.MaxMultipartMemory << 20

	router.Route(app)

	// Listen and Serve
	app.Run(*addr)
}
