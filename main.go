package main

import (
	"os"

	"github.com/ory-am/hydra/cmd"
	"github.com/pkg/profile"
	"github.com/Sirupsen/logrus"
)

func main() {
	if os.Getenv("HYDRA_PROFILING") == "1" {
		defer profile.Start().Stop()
	}

	switch (os.Getenv("LOG_LEVEL")) {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		break;
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
		break;
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
		break;
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
		break;
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
		break;
	}

	cmd.Execute()
}
