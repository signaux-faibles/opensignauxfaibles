package test

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"

	"opensignauxfaibles/tools/altares/pkg/utils"
)

var s3AccessKey = "test"
var s3SecretKey = "testtest"

var s3Credentials = credentials.NewStaticV4(s3AccessKey, s3SecretKey, "")

const s3ContainerName = "s3_by_minio"

func NewS3ForTest(t *testing.T) *minio.Client {
	s3Test := startMinio(t)
	apiHostAndPort := s3Test.GetHostPort("9000/tcp")
	slog.Debug(
		"l'api S3 est disponible",
		slog.String("endpoint", apiHostAndPort),
	)
	time.Sleep(time.Second)
	client, err := minio.New(apiHostAndPort, &minio.Options{Creds: s3Credentials})
	require.NoError(t, err)
	return client
}

func startMinio(t *testing.T) *dockertest.Resource {
	dir, err := os.MkdirTemp(os.TempDir(), "s3_volume_*")
	require.NoError(t, err)
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	s3, found := pool.ContainerByName(s3ContainerName)
	if found {
		slog.Info("le containers s3 existe déjà", slog.String("volume", sourcesFrom(s3.Container.Mounts)))
		return s3
	}
	slog.Info(
		"démarre le container minio s3",
		slog.String("name", s3ContainerName),
		slog.String("path", dir),
	)
	s3, err = pool.RunWithOptions(&dockertest.RunOptions{
		Name:       s3ContainerName,
		Repository: "quay.io/minio/minio",
		//Tag:        "15-alpine",
		Cmd: []string{"server", "/data"},
		Env: []string{
			"MINIO_ROOT_USER=" + s3AccessKey,
			"MINIO_ROOT_PASSWORD=" + s3SecretKey,
		},
		Mounts: []string{
			dir + ":/data",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		killContainer(s3)
		slog.Error(
			"erreur pendant le démarrage du container",
			slog.String("container", s3ContainerName),
			slog.Any("error", err),
		)
		require.NoError(t, err)
	}
	// container stops after 20'
	if err = s3.Expire(10); err != nil {
		killContainer(s3)
		slog.Error(
			"erreur pendant la configuration de l'expiration du container",
			slog.String("container", s3ContainerName),
			slog.Any("error", err),
		)
		panic(err)
	}
	wait := time.Second
	slog.Debug("attends que S3 soit prêt", slog.Any("wait", wait))
	time.Sleep(wait)
	return s3
}

func sourcesFrom(mounts []docker.Mount) string {
	r := ""
	for _, m := range mounts {
		r = fmt.Sprint(m.Source + ":" + m.Destination)
	}
	return r
}

func killContainer(resource *dockertest.Resource) {
	if resource == nil {
		return
	}
	if err := resource.Close(); err != nil {
		log.Panicf("Erreur lors de la purge des resources : %v", err)
	}
}

func wait4S3IsReady(s3Endpoint string) {
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool, err := dockertest.NewPool("")
	utils.ManageError(err, "erreur à création du pool docker par défaut")
	pool.MaxWait = 60 * time.Second
	if err := pool.Retry(func() error {
		client, err := minio.New(s3Endpoint, &minio.Options{
			Creds:  credentials.NewEnvMinio(),
			Secure: true,
		})
		if err != nil {
			return err
		}
		_, err = client.BucketExists(context.Background(), "random")
		return err
	}); err != nil {
		slog.Error(
			"erreur lors de la connexion au conteneur S3",
			slog.Any("error", err),
			slog.String("url", s3Endpoint),
		)
		panic(err)
	}
	slog.Debug("le stockage objet est prêt", slog.String("url", s3Endpoint))
}
