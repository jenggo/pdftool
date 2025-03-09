package cron

import (
	"context"
	"math"
	"pdftool/types"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

func New() *cron.Cron {
	c := cron.New()

	s3, err := minio.New(types.Config.S3.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(types.Config.S3.Key.Access, types.Config.S3.Key.Secret, ""),
		Secure: true,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create minio client")
		return nil
	}
	s3.SetAppInfo(types.AppName, types.AppVersion)

	found, err := s3.BucketExists(context.Background(), types.Config.S3.Bucket)
	if err != nil {
		log.Error().Err(err).Msg("failed to check bucket existence")
		return nil
	}

	if !found {
		log.Error().Msgf("bucket %s not found", types.Config.S3.Bucket)
		return nil
	}

	_, _ = c.AddFunc("@every hour", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		objectList := s3.ListObjects(ctx, types.Config.S3.Bucket, minio.ListObjectsOptions{Recursive: true})
		for object := range objectList {
			if object.Err != nil {
				continue
			}

			lastmodified := object.LastModified.Local()
			umur := math.Round(time.Since(lastmodified).Hours())

			if umur > 1 {
				if err := s3.RemoveObject(ctx, types.Config.S3.Bucket, object.Key, minio.RemoveObjectOptions{ForceDelete: true}); err != nil {
					log.Error().Err(err).Msgf("failed to remove %s", object.Key)
				}
			}
		}
	})

	c.Start()
	return c
}
