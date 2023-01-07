package core

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around apt to make it works inside a container
	from outside, directly on the host.
*/

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/vanilla-os/apx/settings"
)

type ContainerType int

const (
	APT ContainerType = iota // 0
	AUR ContainerType = iota // 1
	DNF ContainerType = iota // 2
	APK ContainerType = iota // 3
)

type Container struct {
	containerType ContainerType
	customName    string
}

func NewContainer(kind ContainerType) *Container {
	return &Container{
		containerType: kind,
	}
}
func NewNamedContainer(kind ContainerType, name string) *Container {
	return &Container{
		containerType: kind,
		customName:    name,
	}
}
func (c *Container) GetContainerImage() (image string, err error) {
	switch c.containerType {
	case APT:
		return GetHostImage()
	case AUR:
		return "docker.io/library/archlinux", nil
	case DNF:
		return "docker.io/library/fedora", nil
	case APK:
		return "docker.io/library/alpine", nil
	default:
		image = ""
		err = errors.New("can't retrieve image for unknown container")
	}
	return image, err
}

func (c *Container) GetContainerName() (name string) {
	var cn strings.Builder
	switch c.containerType {
	case APT:
		cn.WriteString("apx_managed")
	case AUR:
		cn.WriteString("apx_managed_aur")
	case DNF:
		cn.WriteString("apx_managed_dnf")
	case APK:
		cn.WriteString("apx_managed_apk")
	default:
		log.Fatal(fmt.Errorf("unspecified container type"))
	}
	if len(c.customName) > 0 {
		cn.WriteString("_")
		cn.WriteString(strings.Replace(c.customName, " ", "", -1))
	}
	return cn.String()
}

func ContainerManager() string {
	docker := exec.Command("sh", "-c", "command -v docker")
	podman := exec.Command("sh", "-c", "command -v podman")

	// prefer podman over docker if both are present
	if err := podman.Run(); err == nil {
		return "podman"
	} else if err := docker.Run(); err == nil {
		return "docker"
	}

	log.Fatal("no container engine found. Please install Podman or Docker.")
	return ""
}

func GetHostImage() (img string, err error) {
	if settings.Cnf.Image != "" {
		return settings.Cnf.Image, nil
	}

	distro_raw, err := exec.Command("lsb_release", "-is").Output()
	if err != nil {
		return "", err
	}
	distro := strings.ToLower(strings.Trim(string(distro_raw), "\r\n"))

	release_raw, err := exec.Command("lsb_release", "-rs").Output()
	if err != nil {
		return "", err
	}
	release := strings.ToLower(strings.Trim(string(release_raw), "\r\n"))

	return fmt.Sprintf("%v:%v", distro, release), nil
}

func GetDistroboxVersion() (version string, err error) {
	output, err := exec.Command(settings.Cnf.DistroboxPath, "version").Output()
	if err != nil {
		return "", err
	}

	splitted := strings.Split(string(output), "distrobox: ")
	if len(splitted) != 2 {
		return "", errors.New("can't retrieve distrobox version")
	}

	return splitted[1], nil
}

func (c *Container) Exec(capture_output bool, args ...string) (string, error) {
	ExitIfOverlayTypeFS()

	if !c.Exists() {
		err := c.Create()
		if err != nil {
			log.Default().Println("Failed to initialize the container. Try manually with `apx init`.")
			return "", err
		}
	}

	container_name := c.GetContainerName()

	cmd := exec.Command(settings.Cnf.DistroboxPath, "enter", container_name, "--")
	cmd.Args = append(cmd.Args, args...)
	cmd.Env = os.Environ()

	if capture_output {
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}

		return string(out), nil
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		return "", cmd.Run()
	}
}

// Run executes a command with args inside the container, piping stdout, stderr,
// and stdin to the shell.
func (c *Container) Run(args ...string) error {
	_, err := c.Exec(false, args...)
	return err
}

// Output executes a command with args insinde the container, capturing and
// returning the output
func (c *Container) Output(args ...string) (string, error) {
	return c.Exec(true, args...)
}

func (c *Container) Enter() error {
	ExitIfOverlayTypeFS()

	if !c.Exists() {
		log.Default().Printf("Managed container does not exist.\nTry: apx init")
		return errors.New("managed container does not exist")
	}

	container_name := c.GetContainerName()

	cmd := exec.Command(settings.Cnf.DistroboxPath, "enter", container_name)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		if err.Error() != "exit status 130" {
			return err
		}
	}

	return nil
}

func (c *Container) Create() error {
	ExitIfOverlayTypeFS()

	if !CheckConnection() {
		log.Default().Println("No internet connection. Please connect to the internet and try again.")
		return errors.New("failed to create container")
	}

	container_image, err := c.GetContainerImage()
	if err != nil {
		return err
	}

	container_name := c.GetContainerName()
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = " Creating container..."

	spinner.Start()

	cmd := exec.Command(settings.Cnf.DistroboxPath, "create",
		"--name", container_name,
		"--image", container_image,
		"--yes",
		"--no-entry",
		"--additional-flags",
		"--label=manager=apx",
		"--yes")
	cmd.Env = os.Environ()
	// mute command output
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//cmd.Stdin = os.Stdin
	//err = cmd.Run()
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("error creating container: %v", err)
	}

	spinner.Stop()

	if c.containerType == AUR {
		DownloadYay()
		c.Run(GetAurPkgCommand("install-yay")...)
	}

	log.Default().Println("Container created")

	return err
}

func (c *Container) Stop() error {
	ExitIfOverlayTypeFS()

	container_name := c.GetContainerName()
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = " Stopping container..."

	spinner.Start()

	cmd := exec.Command(settings.Cnf.DistroboxPath, "stop", container_name, "--yes")
	_, err := cmd.Output()

	spinner.Stop()

	if err != nil {
		log.Fatalf("error stopping container: %v", err)
	}

	log.Default().Println("Container stopped")

	return err
}

func (c *Container) Remove() error {
	ExitIfOverlayTypeFS()

	container_name := c.GetContainerName()
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = " Removing container..."

	if !c.Exists() {
		return nil
	}

	err := c.Stop()
	if err != nil {
		return err
	}

	spinner.Start()

	cmd := exec.Command(settings.Cnf.DistroboxPath, "rm", container_name, "--yes")
	_, err = cmd.Output()

	spinner.Stop()

	log.Default().Println("Container removed")

	return err
}

func (c *Container) ExportDesktopEntry(program string) {
	c.Run("sh", "-c", "distrobox-export --app "+program+" 2>/dev/null || true")
}

func (c *Container) ExportBinary(bin string) error {
	// Get host's $PATH
	out, err := c.Output("sh", "-c", "distrobox-host-exec printenv | grep -E ^PATH=")
	if err != nil {
		return err
	}

	// If bin name not in $PATH, export to .local/bin
	// Otherwise, export with suffix based on container name
	if !strings.HasPrefix(out, "PATH=") {
		return errors.New("Failed to read host's $PATH")
	}
	_, host_path, _ := strings.Cut(out, "=")

	paths := strings.Split(host_path, ":")
	bin_rename := ""
	for _, path := range paths {
		// Skip directory if it doesn't exist
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		duplicate_found := false
		for _, entry := range entries {
			if entry.Name() == bin {
				switch c.containerType {
				case APT:
					bin_rename = fmt.Sprintf("apt_%s", bin)
				case AUR:
					bin_rename = fmt.Sprintf("aur_%s", bin)
				case DNF:
					bin_rename = fmt.Sprintf("dnf_%s", bin)
				case APK:
					bin_rename = fmt.Sprintf("apk_%s", bin)
				default:
					return errors.New("can't export binary from unknown container")
				}

				fmt.Printf("Warning: another program with name `%s` already exists on host, exporting as `%s`.\n", bin, bin_rename)
				duplicate_found = true
				break
			}
		}

		// No need to keep searching if we alrady found a duplicate name
		if duplicate_found {
			break
		}
	}

	bin_path, err := c.Output("sh", "-c", "command -v "+bin)
	if err != nil {
		return err
	}

	c.Run("sh", "-c", "distrobox-export --bin "+string(bin_path)+" --export-path ~/.local/bin >/dev/null 2>/dev/null || true")
	if bin_rename != "" {
		if err := c.Run("sh", "-c", "mv ~/.local/bin/"+bin+" ~/.local/bin/"+bin_rename); err != nil {
			return err
		}

		fmt.Printf("Binary exported to `~/.local/bin/%s`.\n", bin_rename)
		return nil
	}

	fmt.Printf("Binary exported to `~/.local/bin/%s`.\n", bin)
	return nil
}

func (c *Container) RemoveDesktopEntry(program string) error {
	container_name := c.GetContainerName()
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = fmt.Sprintf("Removing desktop entry: %v\n", program)

	home_dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	spinner.Start()

	files, err := ioutil.ReadDir(home_dir + "/.local/share/applications")
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(strings.ToLower(file.Name()),
			strings.ToLower(container_name+"-"+program)) {
			spinner.Stop()
			err := os.Remove(home_dir + "/.local/share/applications/" + file.Name())
			if err != nil {
				return err
			}
		}
	}

	spinner.Stop()

	log.Default().Printf("Desktop entry %v not found.\n", program)
	return nil
}

func (c *Container) Exists() bool {
	container_name := c.GetContainerName()
	manager := ContainerManager()

	cmd := exec.Command(manager, "ps", "-a", "-q", "-f", "name="+container_name+"$")
	output, _ := cmd.Output()

	// fmt.Println("container_name: ", container_name)
	// fmt.Println("command: ", cmd.String())
	// fmt.Println("output: ", string(output))

	return len(output) > 0
}
