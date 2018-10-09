package corsx

import (
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/ory/go-convenience/stringsx"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func ParseOptions() cors.Options {
	allowCredentials, _ := strconv.ParseBool(viper.GetString("CORS_ALLOWED_CREDENTIALS"))
	debug, _ := strconv.ParseBool(viper.GetString("CORS_DEBUG"))
	maxAge, _ := strconv.Atoi(viper.GetString("CORS_MAX_AGE"))
	return cors.Options{
		AllowedOrigins:   stringsx.Splitx(viper.GetString("CORS_ALLOWED_ORIGINS"), ","),
		AllowedMethods:   stringsx.Splitx(viper.GetString("CORS_ALLOWED_METHODS"), ","),
		AllowedHeaders:   stringsx.Splitx(viper.GetString("CORS_ALLOWED_HEADERS"), ","),
		ExposedHeaders:   stringsx.Splitx(viper.GetString("CORS_EXPOSED_HEADERS"), ","),
		AllowCredentials: allowCredentials,
		MaxAge:           maxAge,
		Debug:            debug,
	}
}

func Initialize(h http.Handler, l logrus.FieldLogger) http.Handler {
	if viper.GetString("CORS_ENABLED") == "true" {
		l.Info("CORS is enabled")
		return cors.New(ParseOptions()).Handler(h)
	}

	return h
}
