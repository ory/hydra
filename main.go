package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/ory-am/hydra/cmd"
	"github.com/pkg/profile"
)

func main() {
	if os.Getenv("PROFILING") == "cpu" {
		defer profile.Start(profile.CPUProfile).Stop()
	} else if os.Getenv("PROFILING") == "memory" {
		defer profile.Start(profile.MemProfile).Stop()
	}

	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		break
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
		break
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
		break
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
		break
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
		break
	}

	cmd.Execute()
}
