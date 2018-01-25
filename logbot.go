package mergbigfolder 

import (
	log "github.com/Sirupsen/logrus"
	lum "gopkg.in/natefinch/lumberjack.v2"
	"fmt"
	"runtime"
	"strings"
	"strconv"
	"io"
	"os"
	_ "time"
)

const (
	filename = "scraper.log"
	rotate   = 10
)

type logBot struct {
	file *os.File
	debug_level string
}

func newLogBot(dl string) *logBot {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//log.SetOutput(os.Stdout)

	lb := new(logBot)
	logger := &lum.Logger{
		Filename:   "log/" + filename,
		MaxSize:    100, // megabytes
		MaxBackups: 100,
		MaxAge:     28, //days
	}
	log.SetOutput(io.MultiWriter(logger, os.Stdout))

	// Only log the warning severity or above.
	//fmt.Printf(*debug_level)
	var l log.Level
	switch dl{
	case "debug":
		l = log.DebugLevel
	case "info":
		l = log.InfoLevel
	case "warn":
		l = log.WarnLevel
	case "error":
		l = log.ErrorLevel
	default:
		l= log.DebugLevel
	}

	log.SetLevel(l)

	go func() {
		//	<-closeQue
		//	logger.Rotate()
	}()

	return lb
}

func (lb *logBot) close() {
	lb.file.Close()
}

func (lb *logBot) reset() {
	lb.file.Close()
	var err error
	lb.file, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err == nil {
		//	log.SetOutput(file)
		log.SetOutput(io.MultiWriter(lb.file, os.Stdout))

	} else {
		log.Fatal("Failed to log to file, using default stderr")
	}
}

func (lb *logBot) GoID() int {
	return GoID()
}


func GoID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
func (lb *logBot) Debug(args ...interface{}) {
	log.WithFields(log.Fields{"gid": GoID()}).Debug(args...)
}

func (lb *logBot) Info( args ...interface{}) {
	log.WithFields(log.Fields{"gid": GoID()}).Debug(args...)
}
