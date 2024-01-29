package test

import (
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/jaswdr/faker"

	"opensignauxfaibles/tools/altares/pkg/utils"
)

var Fake faker.Faker

func init() {
	Fake = faker.NewWithSeed(rand.NewSource(time.Now().UnixNano()))
}

func FakeBucketName() string {
	bucketName := ""
	for len(bucketName) <= 3 {
		bucketName = Fake.Lorem().Word()
	}
	return bucketName
}

func CreateRandomFile() *os.File {
	temp, err := os.CreateTemp(os.TempDir(), "fake_*")
	utils.ManageError(err, "erreur à la création du fichier temporaire")
	err = os.WriteFile(temp.Name(), Fake.Lorem().Bytes(10124), 666)
	utils.ManageError(err, "erreur à l'écriture du fichier temporaire", slog.Any("file", temp))
	return temp
}
