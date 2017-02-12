package persistence

import (
	"gopkg.in/ory-am/dockertest.v3"
	"testing"
	"gopkg.in/mgo.v2"
	"fmt"
)

func TestMongoRepository_Query(t *testing.T) {
	pool, dockerResource := connectToDockerMongo(t)
	defer pool.Purge(dockerResource)
}

func connectToDockerMongo(t *testing.T) (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	fatalIfError(t, err)

	dockerResource, err := pool.Run("mongo", "3.3", []string{""})
	fatalIfError(t, err)

	err = pool.Retry(func() error {
		var e error
		session, e := mgo.Dial(fmt.Sprintf("localhost:%s", dockerResource.GetPort("27017")))
		if e != nil {
			t.Error(e)
			return e
		}
		return session.Ping()
	})
	fatalIfError(t, err)

	return pool, dockerResource
}

func fatalIfError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
