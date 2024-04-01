package main_test

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

func TestMyProjectInDocker(t *testing.T) {

	// Enable detailed logging for Testcontainers
	os.Setenv("DEBUG", "true")

	ctx := context.Background()

	// Obtain the current directory's absolute path
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %s", err)
	}

	// Define the Docker container request
	req := testcontainers.ContainerRequest{
		Image:      "ubuntu:latest",
		Cmd:        []string{"/bin/bash", "-c", "while true; do sleep 5; done"},                                   // Keeps the container running
		WaitingFor: wait.ForLog("sleep 5").WithStartupTimeout(10 * time.Second).WithPollInterval(3 * time.Second), // Simple wait strategy
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			// Using HostConfigModifier to specify mounts
			hostConfig.Mounts = []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: currentDir, // Host path to bind mount
					Target: "/lda",     // Container path where the host path is mounted
				},
			}
			// Enable auto-removal of the container on stop
			hostConfig.AutoRemove = true
		},
	}

	// Start the container
	log.Println("Starting the container...")
	cntr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %s", err)
	}
	defer cntr.Terminate(ctx)

	// Execute a command inside the running container
	fmt.Printf("Executing command in container %s\n", cntr.GetContainerID())
	i, execResult, err := cntr.Exec(ctx, []string{"/bin/bash", "-c", "myproject --version"})
	if err != nil {
		t.Fatalf("Failed to execute command inside container: %s", err)
	}

	fmt.Printf("Command executed in container %d\n", i)
	stdout, err := io.ReadAll(execResult)
	if err != nil {
		t.Fatalf("Failed to read command output: %s", err)
	}

	fmt.Printf("Command output: %s\n", string(stdout))

	// Here you could add assertions based on the output of your project's execution
}
