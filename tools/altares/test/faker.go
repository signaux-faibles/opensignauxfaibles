package test

import (
	"math/rand"
	"time"

	"github.com/jaswdr/faker"
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
