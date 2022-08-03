package main

import (
	"flag"

	"github.com/sirupsen/logrus"
)

var (
	debug = flag.Bool("debug", false, "Enables debug logging")
)

func main() {
	flag.Parse()
	log := logrus.NewEntry(logrus.StandardLogger())
	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	log.Infoln("Hello world")
}
