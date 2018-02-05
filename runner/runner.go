/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 18-11-2017
 * |
 * | File Name:     runner/runner.go
 * +===============================================
 */

package runner

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/docker/go-connections/nat"

	client "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
)

var dockerClient *client.Client

func init() {
	// NewEnvClient initializes a new API client based on environment variables.
	// Use DOCKER_HOST to set the url to the docker server.
	// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	dockerClient = cli
}

// New creates runner docker with given user name
func New(name string) string {
	ctx := context.Background()

	imageName := "aiotrc/gorunner"

	_, err := dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	resp, err := dockerClient.ContainerCreate(ctx,
		&container.Config{
			Image: imageName,
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"8080": []nat.PortBinding{
					nat.PortBinding{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprintf("%d", 8080+rand.Intn(100)),
					},
				},
			},
		}, nil, fmt.Sprintf("el-%s", name))
	if err != nil {
		panic(err)
	}

	if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID
}
