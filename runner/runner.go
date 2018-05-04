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
	"strconv"

	"github.com/docker/go-connections/nat"
	"github.com/phayes/freeport"

	client "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
)

var dockerClient *client.Client

// Runner represents runner docker information
type Runner struct {
	ID      string `json:"id" bson:"id"`
	Port    string `json:"port" bson:"port"`
	RedisID string `json:"rid" bson:"redisid"`
}

func init() {
	// NewEnvClient initializes a new API client based on environment variables.
	// Use DOCKER_HOST to set the url to the docker server.
	// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	// TODO corrects network creation
	/*
		if _, err := cli.NetworkCreate(context.Background(), "isrc", types.NetworkCreate{
			CheckDuplicate: false,
		}); err != nil {
			panic(err)
		}
	*/

	dockerClient = cli
}

// New creates runner docker with given user name
// mgu represents mongo url that is used in runners
// for collecting errors and access to thing data
func New(name string, envs []Env) (Runner, error) {
	ctx := context.Background()

	rid, err := createRedis(name)

	if err != nil {
		return Runner{}, err
	}

	gid, eport, err := createRunner(name, envs)

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
		ID:      gid,
		Port:    eport,
		RedisID: rid,
	}, nil
}

func createRedis(name string) (string, error) {
	ctx := context.Background()

	imageName := "redis:alpine"

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

func createRunner(name string, envs []Env) (string, string, error) {
	ctx := context.Background()

	imageName := "aiotrc/gorunner"

	lport, _ := nat.NewPort("tcp", "8080")
	eport, err := freeport.GetFreePort()
	if err != nil {
		return "", "", err
	}

	dockerEnvs := []string{
		fmt.Sprintf("REDIS_HOST=rd_%s", name),
		fmt.Sprintf("NAME=%s", name),
	}

	// There is at least one user defined environment variables
	if envs != nil {
		for _, e := range envs {
			dockerEnvs = append(dockerEnvs, fmt.Sprintf("%s=%s", e.Name, e.Value))
		}
	}

	resp, err := dockerClient.ContainerCreate(ctx,
		&container.Config{
			Image: imageName,
			ExposedPorts: nat.PortSet{
				lport: struct{}{},
			},
			Env: dockerEnvs,
		},
		&container.HostConfig{
			NetworkMode: "isrc",
			PortBindings: nat.PortMap{
				lport: []nat.PortBinding{
					nat.PortBinding{
						HostIP:   "0.0.0.0",
						HostPort: strconv.Itoa(eport),
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

	return resp.ID, strconv.Itoa(eport), nil
}

// Remove removes runner docker
func (r Runner) Remove() error {
	ctx := context.Background()

	if err := dockerClient.ContainerRemove(ctx, r.RedisID, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return err
	}

	err := dockerClient.ContainerRemove(ctx, r.ID, types.ContainerRemoveOptions{
		Force: true,
	})
	return err
}
