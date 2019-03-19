package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

var cowStr = `	
     /                       \
    /X/                       \X\
   |XX\         _____         /XX|
   |XXX\     _/       \_     /XXX|___________
   \XXXXXXX             XXXXXXX/            \\\
      \XXXX    /     \    XXXXX/                \\\
           |   0     0   |                         \
            |           |                           \
             \         /                            |______//
              \       /                             |
               | O_O | \                            |
                \ _ /   \________________           |
                           | |  | |      \         /
     No Bullshit,          / |  / |       \______/
      Please...            \ |  \ |        \ |  \ |
                         __| |__| |      __| |__| |
                         |___||___|      |___||___|
`

func simpleHandler(name string, args ...string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		cmd := exec.Command(name, args...)
		//fmt.Println("cmd:", name, strings.Join(args, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Error: %s\n", err)
		}
		return nil
	}
}

func daemondHandler(name string, args ...string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		cmd := exec.Command(name, args...)
		fmt.Println("cmd:", name, strings.Join(args, " "))
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatalf("Error: %s\n", err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatalf("Error: %s\n", err)
		}
		return nil
	}
}

func shellHandler(line string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		cmd := exec.Command("sh", "-c", line)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			log.Fatalf("Error: %s\n", err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatalf("Error: %s\n", err)
		}
		return nil
	}
}

func main() {
	// cli config
	app := cli.NewApp()
	app.Name = "lulu-cli"
	app.Usage = "Make development easier"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:   "hello",
			Usage:  "say hello",
			Action: simpleHandler("echo", "Have fun!"),
		},
		{
			Name:   "ganache",
			Usage:  "Start private blockchain on port 8545 with network id 999",
			Action: daemondHandler("ganache-cli", "--port", "8545", "-i", "999"),
		},
		{
			Name:   "server",
			Usage:  "Connect to server",
			Action: daemondHandler("ssh", "fisher@192.168.1.5"),
		},
		{
			Name:   "webapp",
			Usage:  "Web app framework",
			Action: shellHandler("git clone https://github.com/jpsiyu/node-framework.git . && rm -rf .git"),
		},
		{
			Name:   "webapp-go",
			Usage:  "Web app framework go version",
			Action: shellHandler("git clone https://github.com/jpsiyu/webapp-go.git . && rm -rf .git"),
		},
		{
			Name:   "dockercc",
			Usage:  "Clear all docker containers",
			Action: shellHandler("docker rm -f $(docker ps -aq)"),
		},
		{
			Name:   "dockerci",
			Usage:  "Clear docker dangling images",
			Action: daemondHandler("docker", "image", "prune"),
		},
		{
			Name:   "dockerls",
			Usage:  "Show name and tag of docker image",
			Action: shellHandler("docker image ls | awk '{print $1,$2}'"),
		},
		{
			Name:   "space",
			Usage:  "Give me space",
			Action: simpleHandler("echo", cowStr),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
