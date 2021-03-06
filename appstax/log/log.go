package log

import (
	"github.com/mitchellh/go-homedir"
	"fmt"
	"os"
)

var stdoutEnabled = false
var path = "appstax.log"

func Path() string {
	return path
}

func SetStdoutEnabled(e bool) {
	stdoutEnabled = e
	if stdoutEnabled {
		Infof("Log enabled")
	}
}

func Infof(format string, a ...interface{}) {
	write(fmt.Sprintf("[INFO] "+format, a...))
}

func Debugf(format string, a ...interface{}) {
	write(fmt.Sprintf("[DEBUG] "+format, a...))
}

func Warnf(format string, a ...interface{}) {
	write(fmt.Sprintf("[WARN] "+format, a...))
}

func Panicf(format string, a ...interface{}) {
	write(fmt.Sprintf("[PANIC] "+format, a...))
}

func write(line string) {
	writeToSdtout(line)
	writeToFile(line)
}

func writeToSdtout(line string) {
	if stdoutEnabled {
		fmt.Println(line)
	}
}

func writeToFile(line string) {
	f, err := openFile()
	defer f.Close()
	if err == nil {
		f.WriteString(line + "\n")
	}
}

func openFile() (*os.File, error) {
	if path == "appstax.log" {
		home, err := homedir.Dir()
		if err == nil {
			path = home + "/.appstax/appstax.log"			
		}
	}
	return os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
}
