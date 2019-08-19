package driver

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/ory/viper"

	"github.com/ory/hydra/driver/configuration"
)

type DefaultDriver struct {
	c configuration.Provider
	r Registry
}

func NewDefaultDriver(l logrus.FieldLogger, forcedHTTP bool, insecureRedirects []string, version, build, date string, validate bool) Driver {
	fmt.Printf("migrate dsn set viper 2: %s\n", viper.GetString("dsn"))
	c := configuration.NewViperProvider(l, forcedHTTP, insecureRedirects)
	fmt.Printf("migrate dsn set viper 3: %s\n", viper.GetString("dsn"))

	if validate {
		configuration.MustValidate(l, c)
	}

	r, err := NewRegistry(c)
	if err != nil {
		l.WithError(err).Fatal("Unable to instantiate service registry.")
	}

	r.
		WithConfig(c).
		WithLogger(l).
		WithBuildInfo(version, build, date)

	if err = r.Init(); err != nil {
		l.WithError(err).Fatal("Unable to initialize service registry.")
	}

	return &DefaultDriver{r: r, c: c}
}

func (r *DefaultDriver) Configuration() configuration.Provider {
	return r.c
}

func (r *DefaultDriver) Registry() Registry {
	return r.r
}

func (r *DefaultDriver) CallRegistry() Driver {
	CallRegistry(r.Registry())
	return r
}
