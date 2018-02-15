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
	"time"

	"github.com/docker/go-connections/nat"

	client "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
)

var dockerClient *client.Client

// Runner represents runner docker information
type Runner struct {
	ID   string `json:"id"`
	Port string `json:"port"`
	rID  string
}

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
// mge represents mongo url that is used in runners
// for collecting errors
func New(name string, mgu string) (Runner, error) {
	ctx := context.Background()

	rid, err := createRedis(name)

	if err != nil {
		return Runner{}, err
	}

	gid, eport, err := createRunner(name, mgu)

	if err != nil {
		// Removes redis container
		if err := dockerClient.ContainerRemove(ctx, rid, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			return Runner{}, err
		}

		return Runner{}, err
	}

	return Runner{
		ID:   gid,
		Port: eport,
		rID:  rid,
	}, nil
}

func createRedis(name string) (string, error) {
	ctx := context.Background()

	imageName := "redis:alpine"

	_, err := dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}

	lport, _ := nat.NewPort("tcp", "6379")

	resp, err := dockerClient.ContainerCreate(ctx,
		&container.Config{
			Image: imageName,
			ExposedPorts: nat.PortSet{
				lport: struct{}{},
			},
		},
		&container.HostConfig{
			NetworkMode: "isrc",
		}, nil, fmt.Sprintf("rd_%s", name))
	if err != nil {
		return "", err
	}

	if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}

	return resp.ID, nil
}

func createRunner(name string, mgu string) (string, string, error) {
	rand := rand.New(rand.NewSource(time.Now().Unix()))
	ctx := context.Background()

	imageName := "aiotrc/gorunner"

	_, err := dockerClient.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return "", "", err
	}

	lport, _ := nat.NewPort("tcp", "8080")
	eport := fmt.Sprintf("%d", 8080+rand.Intn(100))

	resp, err := dockerClient.ContainerCreate(ctx,
		&container.Config{
			Image: imageName,
			ExposedPorts: nat.PortSet{
				lport: struct{}{},
			},
			Env: []string{
				fmt.Sprintf("REDIS_HOST=rd_%s", name),
				fmt.Sprintf("MONGO_URL=%s", mgu),
			},
		},
		&container.HostConfig{
			NetworkMode: "isrc",
			PortBindings: nat.PortMap{
				lport: []nat.PortBinding{
					nat.PortBinding{
						HostIP:   "0.0.0.0",
						HostPort: eport,
					},
				},
			},
		}, nil, fmt.Sprintf("el_%s", name))
	if err != nil {
		return "", "", err
	}

	if err := dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", "", err
	}

	return resp.ID, eport, nil
}

// Remove removes runner docker
func (r Runner) Remove() error {
	ctx := context.Background()

	if err := dockerClient.ContainerRemove(ctx, r.rID, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return err
	}

	err := dockerClient.ContainerRemove(ctx, r.ID, types.ContainerRemoveOptions{
		Force: true,
	})
	return err
}
