package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type logFormatter struct {
	AppName string
}

func (l *logFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Format(time.RFC3339)
	hostname, _ := os.Hostname()
	return []byte(fmt.Sprintf("%s %s %s[%d]: %s %s\n", timestamp, hostname, l.AppName, os.Getpid(), strings.ToUpper(entry.Level.String()), entry.Message)), nil
}

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&logFormatter{AppName: ""})
}

func AppName(appName string) {
	log.SetFormatter(&logFormatter{AppName: appName})
}

func Level() string {
	return string(log.GetLevel())
}

func SetLevel(level string) error {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		return err
	}
	log.SetLevel(lvl)
	return nil
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Print(args ...interface{}) {
	log.Print(args...)
}

func Warning(args ...interface{}) {
	log.Warning(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Panic(args ...interface{}) {
	log.Panic(args...)
}

func Debugln(args ...interface{}) {
	log.Debugln(args...)
}

func Infoln(args ...interface{}) {
	log.Infoln(args...)
}

func Println(args ...interface{}) {
	log.Println(args...)
}

func Warningln(args ...interface{}) {
	log.Warningln(args...)
}

func Errorln(args ...interface{}) {
	log.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	log.Fatalln(args...)
}

func Panicln(args ...interface{}) {
	log.Panicln(args...)
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Environment(prefix string) {
  	log.Debug("environment variables: ")
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, prefix) {
			log.Debug(e)
		}
	}
}

