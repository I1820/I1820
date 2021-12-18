package runner

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/phayes/freeport"
	"github.com/sirupsen/logrus"
)

const (
	runnerImage string = "i1820/elrunner"
	redisImage  string = "redis:alpine"
	network     string = "i1820_projects"
)

// Manager manages runners and their containers.
type Manager interface {
	New(context.Context, string, []Env) (Runner, error)
	Restart(context.Context, Runner) error
	Remove(context.Context, Runner) error
	Pull(context.Context) ([2]string, error)
}

const (
	NanoCPUs = 1000 * 1000 * 1000
	Memory   = 2 * 1000 * 1000 * 100
)

// DockerManager is a docker based manager.
type DockerManager struct {
	Client *client.Client
}

// New creates a new manager.
func New() (Manager, error) {
	// NewEnvClient initializes a new API client based on environment variables.
	// Use DOCKER_HOST to set the url to the docker server.
	// Use DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest.
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	if _, err := cli.NetworkCreate(context.Background(), network, types.NetworkCreate{
		Driver: "bridge",
		Options: map[string]string{
			"subnet":  "192.168.72.0/24",
			"gateway": "192.168.72.1",
		},
	}); err != nil {
		return nil, err
	}

	return &DockerManager{
		Client: cli,
	}, nil
}

// New creates runner docker with given user name
// mgu represents mongo url that is used in runners
// for collecting errors and access to thing data
func (m *DockerManager) New(ctx context.Context, name string, envs []Env) (Runner, error) {
	rid, err := m.createRedis(ctx, name)

	if err != nil {
		return Runner{}, err
	}

	gid, port, err := m.createRunner(ctx, name, envs)

	if err != nil {
		// Removes redis container
		if err := m.Client.ContainerRemove(ctx, rid, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			return Runner{}, err
		}

		return Runner{}, err
	}

	return Runner{
		ID:      gid,
		Port:    port,
		RedisID: rid,
	}, nil
}

// createRedis creates a redis container by using rd_{name} as its name.
func (m *DockerManager) createRedis(ctx context.Context, name string) (string, error) {
	port, _ := nat.NewPort("tcp", "6379")

	resp, err := m.Client.ContainerCreate(ctx,
		&container.Config{
			Image: redisImage,
			ExposedPorts: nat.PortSet{
				port: struct{}{},
			},
		},
		&container.HostConfig{
			NetworkMode: container.NetworkMode(network),
		}, nil, nil, fmt.Sprintf("rd_%s", name))
	if err != nil {
		return "", err
	}

	if err := m.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		if err := m.Client.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			logrus.Errorf("container deletion %s", err)
		}

		return "", err
	}

	return resp.ID, nil
}

// createRunner creates a runner container by using el_{name} as its name.
func (m *DockerManager) createRunner(ctx context.Context, name string, envs []Env) (string, string, error) {
	lport, _ := nat.NewPort("tcp", "8080") // local port

	eport, err := freeport.GetFreePort() // exposed port
	if err != nil {
		return "", "", err
	}

	dockerEnvs := []string{
		fmt.Sprintf("REDIS_HOST=rd_%s", name),
		fmt.Sprintf("NAME=%s", name),
		"PORT=8080",
		"ADDR=0.0.0.0",
	}

	// There is at least one user defined environment variables
	for _, e := range envs {
		dockerEnvs = append(dockerEnvs, fmt.Sprintf("%s=%s", e.Name, e.Value))
	}

	resp, err := m.Client.ContainerCreate(ctx,
		&container.Config{
			Image: runnerImage,
			ExposedPorts: nat.PortSet{
				lport: struct{}{},
			},
			Env: dockerEnvs,
		},
		&container.HostConfig{
			Resources: container.Resources{
				Memory:   Memory,
				NanoCPUs: NanoCPUs,
			},
			NetworkMode: container.NetworkMode(network),
			PortBindings: nat.PortMap{
				lport: []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: strconv.Itoa(eport),
					},
				},
			},
		}, nil, nil, fmt.Sprintf("dstn_%s", name))
	if err != nil {
		return "", "", err
	}

	if err := m.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		if err := m.Client.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			logrus.Errorf("container deletion %s", err)
		}

		return "", "", err
	}

	return resp.ID, strconv.Itoa(eport), nil
}

// Restart restarts runner docker (not redis)
func (m *DockerManager) Restart(ctx context.Context, r Runner) error {
	td := 1 * time.Second

	return m.Client.ContainerRestart(ctx, r.ID, &td)
}

// Remove removes runner and redis dockers
func (m *DockerManager) Remove(ctx context.Context, r Runner) error {
	if err := m.Client.ContainerRemove(ctx, r.RedisID, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		logrus.Errorf("container deletion %s", err)
	}

	if err := m.Client.ContainerRemove(ctx, r.ID, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		logrus.Errorf("container deletion %s", err)
	}

	return nil
}

// Pull pulls latest images of i1820/elrunner and redis:alpine.
// Please consider that image names are defined globally.
func (m *DockerManager) Pull(ctx context.Context) ([2]string, error) {
	var results [2]string

	re, err := m.Client.ImagePull(ctx, runnerImage, types.ImagePullOptions{})
	if err != nil {
		return results, err
	}

	be, err := ioutil.ReadAll(re)
	if err != nil {
		return results, err
	}

	rr, err := m.Client.ImagePull(ctx, redisImage, types.ImagePullOptions{})
	if err != nil {
		return results, err
	}

	br, err := ioutil.ReadAll(rr)
	if err != nil {
		return results, err
	}

	results[0] = string(be)
	results[1] = string(br)

	return results, nil
}
