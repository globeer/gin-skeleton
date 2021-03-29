package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	"bower.co.kr/c4bapi/config"
	"bower.co.kr/c4bapi/routers"
	"github.com/gin-gonic/gin"
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

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CreateFile(name string) error {
	fo, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func() {
		fo.Close()
	}()
	return nil
}

func main() {

	// we have to ensure the directory we want to write to
	// already exists
	// err := os.MkdirAll(logDir, 0755)
	// if err != nil {
	// 	log.Printf("os.MkdirAll error is %v\n", err.Error())
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

	if !FileExists("log/request.log") {
		CreateFile("log/request.log")
	}
	logFile, err := os.OpenFile("log/request.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		log.Fatalf("os.OpenFile error is %v\n", err.Error())
	}
	// log.SetOutput(io.Writer(logFile))
	gin.DefaultWriter = io.MultiWriter(logFile)
	defer logFile.Close()

	if !FileExists("log/error.log") {
		CreateFile("log/error.log")
	}
	errLogFile, err := os.OpenFile("log/error.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		log.Fatalf("os.OpenFile error is %v\n", err.Error())
	}
	// log.SetOutput(io.Writer(logFile))
	gin.DefaultErrorWriter = io.MultiWriter(errLogFile)
	defer errLogFile.Close()

	// logDir := "log"
	// logFile, err := os.OpenFile(logDir+"/request.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	// if err != nil {
	// 	log.Printf("os.OpenFile error is %v\n", err.Error())
	// 	logFile, _ = os.Create(logDir + "/request.log")
	// }
	// gin.DefaultWriter = io.MultiWriter(logFile)
	// log.Printf("logFile is %v\n", logFile.Name())
	// defer logFile.Close()
	// errlogfile, err := os.OpenFile(logDir+"/error.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	// if err != nil {
	// 	log.Printf("os.OpenFile error is %v\n", err.Error())
	// 	errlogfile, _ = os.Create(logDir + "/error.log")
	// }
	// gin.DefaultErrorWriter = io.MultiWriter(errlogfile)
	// log.Printf("errlogfile is %v\n", errlogfile.Name())
	// defer errlogfile.Close()

	// dailylog.initRotatedFileMust()

	router := gin.Default()
	router.Static("/images", filepath.Join(config.Server.StaticDir, "images"))
	router.Static("/js", filepath.Join(config.Server.StaticDir, "js"))
	router.Static("/css", filepath.Join(config.Server.StaticDir, "css"))
	router.StaticFile("/favicon.ico", filepath.Join(config.Server.StaticDir, "favicon.ico"))
	router.LoadHTMLGlob(config.Server.ViewDir + "/*")
	router.MaxMultipartMemory = config.Server.MaxMultipartMemory << 20
	routers.Route(router)

	// Listen and Serve
	router.Run(*addr)
}
