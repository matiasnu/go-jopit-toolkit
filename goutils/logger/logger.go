/**
* @author mnunez
 */

package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var Log *logrus.Logger

func InitLog(loggingPath, loggingFile, loggingLevel string) {
	Log = logrus.New()
	// Creating log dir if not exists
	if _, err := os.Stat(loggingPath); os.IsNotExist(err) {
		if err = os.MkdirAll(loggingPath, 0777); err != nil {
			if os.IsPermission(err) {
				fmt.Println("Try fix the permission issue, by creating the dir structure and try again.")
				panic(err)
			}
		}
	}

	f := filepath.Join(loggingPath, loggingFile)
	fmt.Printf("Logging on : %s\n", f)
	file, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err == nil {
		Log.SetOutput(file)
	} else {
		fmt.Println("Failed to log to file, using default stderr : ", err)
		Log.SetOutput(os.Stderr)
	}
	Log.SetOutput(file)
	// Log as JSON instead of the default ASCII formatter.
	Log.Formatter = new(prefixed.TextFormatter)
	Log.Formatter.(*prefixed.TextFormatter).ForceFormatting = true
	Log.Formatter.(*prefixed.TextFormatter).FullTimestamp = true

	// Only log the warning severity or above.
	lvl, err := logrus.ParseLevel(loggingLevel)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	Log.SetLevel(lvl)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

type LogInterface interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

var _ LogInterface = (*logrus.Logger)(nil) // (*logrus.Logger) implements LogInterface
