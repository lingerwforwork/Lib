package container

import (
	"bytes"
	"testing"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
)

func TestDockerContainer(t *testing.T) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	dockerContainer := NewDockerContainer(cli, &container.Config{
		Image:        "python:3.12-slim",
		WorkingDir:   "/app",
		Cmd:          []string{"python", "-c", "import datetime,random;d=[random.randint(1,100) for _ in range(5)];s={'总数':len(d),'总和':sum(d),'平均值':round(sum(d)/len(d),2),'最大值':max(d),'最小值':min(d)};open('test.txt','w',encoding='utf-8').write(f'生成时间: {datetime.datetime.now()}\\\\n随机数: {d}\\\\n统计:\\\\n'+'\\\\n'.join(f'  {k}: {v}' for k,v in s.items()));print('✅ 已生成: test.txt')"},
		AttachStdout: true,
		AttachStderr: true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   "/Users/mac/Desktop/work/cn_project/backend/epass-v2/Lib/pkg/container",
				Target:   "/app",
				ReadOnly: false,
			},
		},
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	waitResult, err := dockerContainer.Run([]uint32{1, 2}, &stdout, &stderr)
	defer dockerContainer.Clear()
	if err != nil {
		t.Error(err)
	}
	select {
	case <-waitResult.Result:
		t.Log(stdout.String())
		t.Log(stderr.String())
	case e := <-waitResult.Error:
		t.Fatal("容器执行失败" + e.Error())
	}
}
