package main

import (
	"io"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"pdftool/cron"
	"pdftool/docs"
	"pdftool/server"
	"pdftool/types"

	zlogsentry "github.com/archdx/zerolog-sentry"
	"github.com/gofiber/storage/minio"
	"github.com/ilyakaznacheev/cleanenv"
	mo "github.com/minio/minio-go/v7"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//	@title			pdfTool

//	@host		localhost:2804
//	@BasePath	/
//	@schemes	http https

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description					Enter the token with the `Bearer: ` prefix, e.g. "Bearer abcde12345".
func main() {
	if err := cleanenv.ReadEnv(&types.Config); err != nil {
		log.Fatal().Err(err).Send()
	}

	zerolog.SetGlobalLevel(zerolog.Level(types.Config.App.LogLevel))
	var writeLog io.Writer = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "[Mon] [2006-01-02] [15:04:05]",
	}

	if types.Config.App.Sentry != "" {
		host, _ := os.Hostname()
		w, err := zlogsentry.New(
			types.Config.App.Sentry,
			zlogsentry.WithRelease(types.AppVersion),
			zlogsentry.WithServerName(host),
			zlogsentry.WithSampleRate(1.0),
		)
		if err != nil {
			log.Fatal().Err(err).Msg("error initializing Sentry client")
		}

		writeLog = zerolog.MultiLevelWriter(w, writeLog)
	}
	log.Logger = zerolog.New(writeLog).With().Timestamp().Logger()

	// Set storage
	if types.Config.S3.Enable {
		types.Config.S3.Storage = minio.New(minio.Config{
			Endpoint: types.Config.S3.Endpoint,
			Bucket:   types.Config.S3.Bucket,
			Secure:   true,
			Credentials: minio.Credentials{
				AccessKeyID:     types.Config.S3.Key.Access,
				SecretAccessKey: types.Config.S3.Key.Secret,
			},
			PutObjectOptions: mo.PutObjectOptions{
				UserMetadata: map[string]string{"x-amz-acl": "public-read"}, // force to public access
			},
		})
	}

	// Set swagger if enable
	if types.Config.Swagger.Enable {
		log.Info().Msgf("✓ Swagger enabled, path: %s", types.Config.Swagger.Path)
		u, _ := url.Parse(types.Config.App.BaseURL)
		docs.SwaggerInfo.Host = u.Host
		docs.SwaggerInfo.Title = types.AppName
		docs.SwaggerInfo.Version = types.AppVersion
	}

	// Starting server
	server := server.New()
	defer func() {
		// Shutdown server
		if err := server.Shutdown(); err != nil {
			log.Error().Err(err).Send()
		}
	}()

	// Starting cron
	if types.Config.S3.Enable {
		c := cron.New()
		defer c.Stop()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
}
