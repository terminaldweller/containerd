package main

import (
	"context"
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

func main() {
	if err := redisExample(); err != nil {
		log.Fatal(err)
	}
}

func redisExample() error {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "exampple")
	image, err := client.Pull(ctx, "docker.io/library/redis:alpine", containerd.WithPullUnpack)
	if err != nil {
		return err
	}
	log.Printf("Successfully pulled image %s.\n", image.Name())

	container, err := client.NewContainer(
		ctx,
		"redis-server",
		containerd.WithNewSnapshot("redis-server-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return err
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)
	log.Printf("Successfully created container with ID %s and snapshot with ID redis-server-snapshot", container.ID())

	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		fmt.Println("fuck 1")
		return err
	}
	defer task.Delete(ctx)

	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		fmt.Println("fuck 2")
		return err
	}

	if err := task.Start(ctx); err != nil {
		fmt.Println("fuck 3")
		return err
	}

	time.Sleep(10 * time.Second)

	if err := task.Kill(ctx, syscall.SIGTERM); err != nil {
		fmt.Println("fuck 4")
		return err
	}

	status := <-exitStatusC
	code, _, err := status.Result()
	if err != nil {
		fmt.Println("fuck 5")
		return err
	}
	log.Printf("redis-server exited with status: %d\n", code)

	return nil
}
