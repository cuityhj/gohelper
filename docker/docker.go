package docker

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

func Run(containerConfig *container.Config, containerHostConfig *container.HostConfig, containerName string, timeout time.Duration) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	defer cli.Close()
	resp, err := cli.ContainerCreate(context.Background(), containerConfig,
		containerHostConfig, nil, nil, containerName)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	runCh := make(chan struct{})
	errCh := make(chan error)
	go func() {
		for {
			if stats, err := cli.ContainerInspect(context.Background(), resp.ID); err != nil {
				errCh <- err
				return
			} else {
				if stats.ContainerJSONBase != nil && stats.ContainerJSONBase.State != nil &&
					stats.ContainerJSONBase.State.Running {
					runCh <- struct{}{}
					return
				} else {
					time.Sleep(time.Second)
				}
			}
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-runCh:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("docker run %s timeout", containerName)
	}
}

func Remove(containerName string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	defer cli.Close()
	return cli.ContainerRemove(context.Background(), containerName, types.ContainerRemoveOptions{Force: true})
}

func Restart(containerName string, timeout *int) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	defer cli.Close()
	return cli.ContainerRestart(context.Background(), containerName, container.StopOptions{Timeout: timeout})
}

func Stop(containerName string, timeout *int) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	defer cli.Close()
	return cli.ContainerStop(context.Background(), containerName, container.StopOptions{Timeout: timeout})
}

type DockerContext struct {
	ContainerName string
	Command       string
}

func Exec(ctx *DockerContext) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	defer cli.Close()
	execResp, err := cli.ContainerExecCreate(context.Background(), ctx.ContainerName, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		Cmd:          strslice.StrSlice([]string{ctx.Command}),
	})
	if err != nil {
		return "", err
	}

	resp, err := cli.ContainerExecAttach(context.Background(), execResp.ID, types.ExecStartCheck{
		Detach: false,
		Tty:    false,
	})
	if err != nil {
		return "", err
	}

	defer resp.Close()
	out, err := ioutil.ReadAll(resp.Reader)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
