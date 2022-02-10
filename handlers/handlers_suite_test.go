package handlers_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"log"
	"time"
)

func TestHandlers(t *testing.T) {
	pool, err := dockertest.NewPool("")
	pool.MaxWait = time.Minute * 2

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// 8080 for http, 8443 for https
	options := &dockertest.RunOptions{
		Repository: "wiremock/wiremock",
		PortBindings: map[dc.Port][]dc.PortBinding{
			"8080/tcp": []dc.PortBinding{{HostPort: "8080"}},
			"8443/tcp": []dc.PortBinding{{HostPort: "8443"}},
		},
	}
	resource, rerr := pool.RunWithOptions(options)

	if rerr != nil {
		log.Fatalf("Could not start resource: %s", rerr)
	}

	if err = pool.Retry(func() error {
		wiremockRetryUrl := "http://localhost:8080/__admin/mappings"
		_, _, errs := fiber.Get(wiremockRetryUrl).String()
		if errs != nil {
			err = errs[0]
		} else {
			err = nil
		}
		return err
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")

	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	defer GinkgoRecover()
}
