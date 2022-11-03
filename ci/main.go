package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

const (
	publishAddress = "registry.digitalocean.com/hello-do/hello:latest"
)

func main() {
	// Create dagger client
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Build our app
	project := client.Host().Workdir()

	builder := client.Container().
		From("golang:latest").
		WithMountedDirectory("/src", project).
		WithWorkdir("/src").
		WithEnvVariable("GOOS", "linux").
		WithEnvVariable("GOARCH", "amd64").
		WithEnvVariable("CGO_ENABLED", "0").
		Exec(dagger.ContainerExecOpts{
			Args: []string{"go", "build", "-o", "hello"},
		})

	// Get built binary
	build := builder.File("/src/hello")

	// Publish binary on Alpine base
	addr, err := client.Container().
		From("alpine").
		WithMountedFile("/tmp/hello", build).
		Exec(dagger.ContainerExecOpts{
			Args: []string{"cp", "/tmp/hello", "/bin/hello"},
		}).
		WithEntrypoint([]string{"/bin/hello"}).
		Publish(ctx, publishAddress)
	if err != nil {
		panic(err)
	}

	fmt.Println(addr)
}
