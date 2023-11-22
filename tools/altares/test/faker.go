package test

import (
	"math/rand"
	"time"

	"github.com/jaswdr/faker"
)

var fake faker.Faker

func init() {
	fake = faker.NewWithSeed(rand.NewSource(time.Now().UnixNano()))
}

func FakeBucketName() string {
	bucketName := ""
	for len(bucketName) <= 3 {
		bucketName = fake.Lorem().Word()
	}
	return bucketName
}
