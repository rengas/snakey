package testutils

import (
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"time"
)

func Cleanup(code int, pool *dockertest.Pool, network *dockertest.Network, resources ...*dockertest.Resource) {
	for _, resource := range resources {
		if resource != nil {
			if err := pool.Purge(resource); err != nil {
				log.Fatalf("failed to purge resource: %s", err.Error())
			}
		}
	}
	if network != nil {
		if err := network.Close(); err != nil {
			log.Fatalf("failed to close network: %s", err.Error())
		}
	}
	os.Exit(code)
}

type OperationFunc func() error

func Retry(attempts int, delay time.Duration, factor float64, op OperationFunc) error {
	var err error
	f := 1.
	for attempt := 1; attempt <= attempts; attempt++ {
		err = op()
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(float64(delay) * f))
		f *= factor
	}
	return err
}
