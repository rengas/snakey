//go:build e2e

package test

import (
	"errors"
	"fmt"
	"github.com/ory/dockertest/v3"
	"log"
	"net/http"
	"os"
	"path"
	"snakey/pkg/testutils"
	"testing"
	"time"
)

var httpPort string
var connStr string

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("failed to connect to docker: %s", err.Error())
	}

	network, err := pool.CreateNetwork("snakey")
	if err != nil {
		log.Fatalf("failed to create docker network: %s", err.Error())
	}

	// run migrations
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %s", err.Error())
	}

	projectRoot := path.Join(pwd, "..")

	// run test
	time.Sleep(5 * time.Second)
	// setup consumer container
	apiContainer, err := createAPI(pool, network, projectRoot)
	if err != nil {
		log.Printf("failed to create api: %s", err.Error())
		testutils.Cleanup(1, pool, network, apiContainer)
	}

	testutils.Cleanup(m.Run(), pool, network, apiContainer)
}

func createAPI(pool *dockertest.Pool, network *dockertest.Network, projectRoot string) (*dockertest.Resource, error) {
	buildOpts := &dockertest.BuildOptions{
		Dockerfile: "cmd/api/Dockerfile",
		ContextDir: projectRoot,
	}
	runOpts := &dockertest.RunOptions{
		Name:     "api",
		Networks: []*dockertest.Network{network},
	}

	resource, err := pool.BuildAndRunWithBuildOptions(buildOpts, runOpts)
	if err != nil {
		return resource, fmt.Errorf("failed to build and run: %w", err)
	}

	httpPort = resource.GetPort("8080/tcp")

	if err = pool.Retry(func() error {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/health", httpPort))
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return errors.New("got http error code")
		}

		return nil
	}); err != nil {
		return resource, fmt.Errorf("failed to connect to container: %w", err)
	}

	return resource, nil
}
