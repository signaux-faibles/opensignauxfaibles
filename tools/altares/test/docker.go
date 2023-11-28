package test

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var s3AccessKey = "test"
var s3SecretKey = "testtest"

var s3Credentials = credentials.NewStaticV4(s3AccessKey, s3SecretKey, "")

const s3ContainerName = "s3_by_minio_"

func NewS3ForTest(t *testing.T) *minio.Client {
	s3Test := startMinio(t)
	apiHostAndPort := s3Test.GetHostPort("9000/tcp")
	slog.Debug(
		"l'api S3 est disponible",
		slog.String("endpoint", apiHostAndPort),
	)
	client, err := minio.New(apiHostAndPort, &minio.Options{Creds: s3Credentials})
	require.NoError(t, err)
	return client
}

func startMinio(t *testing.T) *dockertest.Resource {
	s3ContainerName := s3ContainerName + Fake.Lorem().Word() + "_" + Fake.Lorem().Word()
	slog.Info(
		"démarre le container minio s3",
		slog.String("name", s3ContainerName),
	)
	pool, err := dockertest.NewPool("")
	require.NoError(t, err)
	s3, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       s3ContainerName,
		Repository: "quay.io/minio/minio",
		//Tag:        "15-alpine",
		Cmd: []string{"server", "/data"},
		Env: []string{
			"MINIO_ROOT_USER=" + s3AccessKey,
			"MINIO_ROOT_PASSWORD=" + s3SecretKey,
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
		require.NoError(t, err)
	}
	// container stops after 20"
	setContainerExpiration(s3, 20)
	hostPort := s3.GetHostPort("9000/tcp")
	start := time.Now()
	waitS3IsReady(pool, hostPort)
	end := time.Now()
	slog.Info("temps d'attente de démarrage de s3", slog.Any("duration", end.Sub(start)))
	return s3
}

func setContainerExpiration(s3 *dockertest.Resource, seconds uint) {
	if err := s3.Expire(seconds); err != nil {
		killContainer(s3)
		slog.Error(
			"erreur pendant la configuration de l'expiration du container",
			slog.String("container", s3ContainerName),
			slog.Any("error", err),
		)
		panic(err)
	}
}

func killContainer(resource *dockertest.Resource) {
	if resource == nil {
		return
	}
	if err := resource.Close(); err != nil {
		log.Panicf("Erreur lors de la purge des resources : %v", err)
	}
}

func waitS3IsReady(pool *dockertest.Pool, s3Endpoint string) {
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	// the minio client does not do service discovery for you (i.e. it does not check if connection can be established), so we have to use the health check
	if err := pool.Retry(func() error {
		url := fmt.Sprintf("http://%s/minio/health/live", s3Endpoint)
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			slog.Info("s3 pas encore prêt", slog.Int("statusCode", resp.StatusCode))
			return fmt.Errorf("status code not OK")
		}
		slog.Info("s3 prêt")
		return nil
	}); err != nil {
		panic(errors.Wrap(err, "erreur de connexion à Docker"))
	}
	slog.Info("le stockage objet est prêt", slog.String("url", s3Endpoint))
}
