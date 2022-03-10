package log

import (
	"log"

	"github.com/wqvoon/cbox/pkg/flags"
)

func init() {
	flags.ParseAll()

	if flags.IsDebugMode() {
		log.SetFlags(log.Llongfile)
	} else {
		log.SetFlags(0)
	}
}

func TODO(msg ...string) {
	if len(msg) == 0 {
		panic("TODO")
	} else {
		panic("TODO: " + msg[0])
	}
}

func Errorln(msg ...interface{}) {
	if flags.IsDebugMode() {
		log.Panicln(msg...)
	}

	log.Fatalln(msg...)
}

func Errorf(format string, v ...interface{}) {
	if flags.IsDebugMode() {
		log.Panicf(format, v...)
	}

	log.Fatalf(format, v...)
}

var (
	Println = log.Println
	Printf  = log.Printf
)
