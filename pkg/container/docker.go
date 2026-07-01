package container

import (
	"context"
	"io"

	"github.com/lingerwforwork/Lib/pkg/errors"
	"github.com/moby/moby/api/pkg/stdcopy"
	libContainer "github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type DockerContainer interface {
	Run(gpuIds []uint32, stdout, stderr io.Writer) (*client.ContainerWaitResult, error)
	Terminate() error //终止
	Clear() error
}

type Container struct {
	id         string
	cli        *client.Client
	config     *libContainer.Config
	hostConfig *libContainer.HostConfig
}

func NewDockerContainer(cli *client.Client, config *libContainer.Config, hostConfig *libContainer.HostConfig) *Container {
	return &Container{
		cli:        cli,
		config:     config,
		hostConfig: hostConfig,
	}
}

func (container *Container) Run(gpuIds []uint32, stdout, stderr io.Writer) (*client.ContainerWaitResult, error) {
	if container.cli == nil || container.config == nil {
		_, _ = io.WriteString(stderr, "Docker client not initialized")
		return nil, errors.NilError
	}
	ctx := context.Background()
	container.config.AttachStdout = true
	container.config.AttachStderr = true
	//1.创建容器
	createRes, err := container.cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config:     container.config,
		HostConfig: container.hostConfig,
	})
	if err != nil {
		_, _ = io.WriteString(stderr, "failed to create container")
		return nil, err
	}
	container.id = createRes.ID
	//defer container.Terminate()
	//2.连接容器
	attachResp, err := container.cli.ContainerAttach(ctx, createRes.ID, client.ContainerAttachOptions{
		Stderr: true,
		Stdout: true,
		Stream: true,
	})
	if err != nil {
		_, _ = io.WriteString(stderr, "failed to attach container")
		return nil, err
	}
	//3.启动容器
	_, err = container.cli.ContainerStart(ctx, createRes.ID, client.ContainerStartOptions{})
	if err != nil {
		_, _ = io.WriteString(stderr, "failed to start container")
		return nil, err
	}
	_, err = stdcopy.StdCopy(stdout, stderr, attachResp.Reader)
	if err != nil {
		_, _ = io.WriteString(stderr, "failed to attach container")
		return nil, err
	}
	//4.等待容器结束
	waitResult := container.cli.ContainerWait(ctx, createRes.ID, client.ContainerWaitOptions{})
	return &waitResult, nil
}

func (container *Container) Terminate() error {
	if container.cli == nil || container.id == "" {
		return nil
	}
	ctx := context.Background()
	_, err := container.cli.ContainerRemove(ctx, container.id, client.ContainerRemoveOptions{
		Force: true,
	})
	return err
}

func (container *Container) Clear() error {
	return container.Terminate()
}
