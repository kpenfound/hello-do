package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"dagger.io/dagger"
)

// go run terraform/main.go plan|apply|destroy
func main() {
	// determine tf operation
	if len(os.Args) < 2 {
		fmt.Printf("Must supply terraform command to run [plan | apply | destroy]")
		os.Exit(1)
	}
	subtask := os.Args[1]

	var tfcommand []string
	switch subtask {
	case "plan":
		tfcommand = []string{"plan"}
	case "apply":
		tfcommand = []string{"apply", "-auto-approve"}
	case "destroy":
		tfcommand = []string{"apply", "-destroy", "-auto-approve"}
	}

	// create client
	ctx := context.Background()
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		fmt.Printf("Unable to connect to dagger: %v", err)
		os.Exit(1)
	}
	defer client.Close()

	// load terraform directory
	tfdirectory := client.Host().Directory("terraform")

	// load terraform image
	tf := client.Container().From("hashicorp/terraform:latest")

	// mount terraform to container workdir
	tf = tf.WithMountedDirectory("/terraform", tfdirectory)
	tf = tf.WithWorkdir("/terraform")

	// pass through TF_TOKEN for Terraform Cloud
	tkn := os.Getenv("TF_TOKEN")
	tf = tf.WithEnvVariable("TF_TOKEN_app_terraform_io", tkn)

	// terraform init
	tf = tf.Exec(dagger.ContainerExecOpts{
		Args: []string{"init"},
	})

	// ensure terraform operation is never cached
	epoch := fmt.Sprintf("%d", time.Now().Unix())
	tf = tf.WithEnvVariable("CACHEBUSTER", epoch)

	// terraform command
	tf = tf.Exec(dagger.ContainerExecOpts{
		Args: tfcommand,
	})

	// execute against dagger engine
	_, err = tf.ExitCode(ctx)
	if err != nil {
		fmt.Printf("Unable to execute terraform command: %v", err)
		os.Exit(1)
	}
}
