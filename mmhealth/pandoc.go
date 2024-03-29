package mmhealth

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/pkg/errors"
)

func ReportToPDF(inputFileName, outputFileName string) error {

	cmdArgs := []string{
		"--template=/app/template.tex",
		filepath.Join("/files", inputFileName),
		"-o",
		filepath.Join("/files", outputFileName),
	}

	err := runDockerCommand(cmdArgs)
	if err != nil {
		return errors.Wrap(err, "failed to generate the pdf report")
	}

	fmt.Println("PDF generated successfully.")
	return nil
}

var DockerImageProd = "ghcr.io/coltoneshaw/mmhealth"
var DockerImageDev = "mmhealth"

var DockerImage = func() string {
	if BuildVersion == "(devel)" {
		return DockerImageDev
	}
	return DockerImageProd + ":" + "v1.2.1"
}()

// Responsible for passing any docker commands to the mmhealth container.
func runDockerCommand(cmdArgs []string) error {
	currentUser, err := user.Current()
	if err != nil {
		return errors.Wrap(err, "failed to get the current user")
	}

	pwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "failed to get current directory")
	}

	dockerArgs := []string{
		"run",
		"--rm",
		"--platform=linux/amd64",
		"--volume", pwd + ":/files",
		"--user", currentUser.Uid + ":" + currentUser.Gid,
		DockerImage,
	}

	mergedArgs := append(dockerArgs, cmdArgs...)

	return runCommand("docker", mergedArgs)
}

func runCommand(name string, args []string) error {
	cmd := exec.Command(name, args...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start the command")
	}

	go copyOutput(stdout)
	go copyOutput(stderr)

	err = cmd.Wait()
	if err != nil {
		return errors.Wrap(err, "failed to wait for the command to finish")
	}

	return nil
}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
