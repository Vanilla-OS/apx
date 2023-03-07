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
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/vanilla-os/apx/settings"
	"github.com/vanilla-os/orchid/cmdr"
)

type ContainerType int

const (
	APT    ContainerType = iota // 0
	AUR    ContainerType = iota // 1
	DNF    ContainerType = iota // 2
	APK    ContainerType = iota // 3
	ZYPPER ContainerType = iota // 4
	XBPS   ContainerType = iota // 5
	SWUPD  ContainerType = iota // 6
)

// How many container types we offer. Must be always the same
// as the number of options above!
const CONTAINER_TYPES = 7

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
	case ZYPPER:
		return "registry.opensuse.org/opensuse/tumbleweed:latest", nil
	case XBPS:
		return "ghcr.io/void-linux/void-linux:latest-full-x86_64", nil
	case SWUPD:
		return "docker.io/library/clearlinux", nil
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
	case ZYPPER:
		cn.WriteString("apx_managed_zypper")
	case XBPS:
		cn.WriteString("apx_managed_xbps")
	case SWUPD:
		cn.WriteString("apx_managed_swupd")
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
	cmd.Env = append(cmd.Env, "STORAGE_DRIVER=vfs")

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
	cmd.Env = append(cmd.Env, "STORAGE_DRIVER=vfs")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		// last command failed, not apx
		if err.Error() == "exit status 1" {
			return nil
		}
		// avoid panic on ctrl-c
		if err.Error() != "exit status 130" {
			return err
		}
	}

	return nil
}

func (c *Container) Create() error {
	ExitIfOverlayTypeFS()

	if !CheckConnection("cloudflare.com", "443") {
		log.Default().Println("No internet connection. Please connect to the internet and try again.")
		return errors.New("failed to create container")
	}

	container_image, err := c.GetContainerImage()
	if err != nil {
		return err
	}

	container_name := c.GetContainerName()
	spinner, err := cmdr.Spinner.Start("Creating container...")
	if err != nil {
		return err
	}
	defer spinner.Stop()

	cmd := exec.Command(settings.Cnf.DistroboxPath, "create",
		"--name", container_name,
		"--image", container_image,
		"--yes",
		"--no-entry",
		"--additional-flags",
		"--label=manager=apx",
		"--yes")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "STORAGE_DRIVER=vfs")
	// mute command output
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//cmd.Stdin = os.Stdin
	//err = cmd.Run()
	_, err = cmd.Output()
	if err != nil {
		log.Fatalf("error creating container: %v", err)
	}

	if c.containerType == AUR {
		// Setup locales
		spinner.UpdateText("Setting up locales...")
		locales, err := GetArchLocales(c)
		if err != nil {
			log.Fatalf("error looking up locales: %v", err)
		}
		if err := InstallArchLocales(c, locales); err != nil {
			log.Fatalf("error generating locales: %v", err)
		}

		// Try to remove older yay versions downloaded by a previous Arch container
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("error detecting user home directory: %v", err)
		}
		yay_dir := fmt.Sprintf("%v/.local/src/yay", home)
		if err := os.RemoveAll(yay_dir); err != nil {
			log.Fatalf("error removing older yay version: %v", err)
		}

		// Download and install yay
		spinner.UpdateText("Downloading and installing Yay...")
		DownloadYay()
		c.Output(GetAurPkgCommand("install-yay-deps")...)
		c.Run(GetAurPkgCommand("install-yay")...)
	}

	if c.containerType == SWUPD {
		c.Run(GetSwupdPkgCommand("install-os-core-search")...)
	}

	spinner.Success("Container created.")

	return err
}

func (c *Container) Stop() error {
	ExitIfOverlayTypeFS()

	container_name := c.GetContainerName()
	spinner, err := cmdr.Spinner.Start("Stopping container...")
	if err != nil {
		return err
	}
	defer spinner.Stop()

	cmd := exec.Command(settings.Cnf.DistroboxPath, "stop", container_name, "--yes")
	_, err = cmd.Output()

	spinner.Success("Container stopped.")

	if err != nil {
		log.Fatalf("error stopping container: %v", err)
	}

	return err
}

func (c *Container) Remove() error {
	ExitIfOverlayTypeFS()

	container_name := c.GetContainerName()
	spinner, err := cmdr.Spinner.Start("Removing container...")
	if err != nil {
		return err
	}
	defer spinner.Stop()

	if !c.Exists() {
		return nil
	}

	err = c.Stop()
	if err != nil {
		return err
	}

	cmd := exec.Command(settings.Cnf.DistroboxPath, "rm", container_name, "--yes")
	_, err = cmd.Output()

	spinner.Success("Container removed.")

	return err
}

func (c *Container) ExportDesktopEntry(program string) {
	c.Run("sh", "-c", "distrobox-export --app "+program+" >/dev/null 2>/dev/null || true")
}

func (c *Container) ExportBinary(bin string) error {
	// Get host's $PATH
	out, err := c.Output("sh", "-c", "distrobox-host-exec $(readlink -fn $(getent passwd $USER | cut -f 7 -d :)) -l -i -c printenv | grep -E ^PATH=")
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to execute printenv: %s", err))
	}

	spinner, err := cmdr.Spinner.Start(fmt.Sprintf("Exporting binary: %v.", bin))
	if err != nil {
		return err
	}
	defer spinner.Stop()

	// If bin name not in $PATH, export to .local/bin
	// Otherwise, export with suffix based on container name
	if !strings.HasPrefix(out, "PATH=") {
		return errors.New("Failed to read host's $PATH")
	}
	_, host_path, _ := strings.Cut(out, "=")

	// Ensure `~/.local/bin` exists
	local_bin_path := fmt.Sprintf("/home/%s/.local/bin", os.Getenv("USER"))
	if _, err := os.Stat(local_bin_path); os.IsNotExist(err) {
		os.Mkdir(local_bin_path, 0775)
	}

	paths := strings.Split(host_path, ":")
	bin_rename := ""
	for _, path := range paths {
		// Skip directory if it doesn't exist
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("Could not read directory %s: %s", path, err)
		}

		duplicate_found := false
		for _, entry := range entries {
			// If duplicate is located in ~/.local/bin, we'll handle it later
			if entry.Name() == bin && !strings.Contains(path, "/.local/bin") {
				switch c.containerType {
				case APT:
					bin_rename = fmt.Sprintf("apt_%s", bin)
				case AUR:
					bin_rename = fmt.Sprintf("aur_%s", bin)
				case DNF:
					bin_rename = fmt.Sprintf("dnf_%s", bin)
				case APK:
					bin_rename = fmt.Sprintf("apk_%s", bin)
				case ZYPPER:
					bin_rename = fmt.Sprintf("zypper_%s", bin)
				case XBPS:
					bin_rename = fmt.Sprintf("xbps_%s", bin)
				case SWUPD:
					bin_rename = fmt.Sprintf("swupd_%s", bin)
				default:
					return errors.New("can't export binary from unknown container")
				}

				cmdr.Warning.Printf("Another program with name `%s` already exists on host, exporting as `%s`.\n", bin, bin_rename)
				duplicate_found = true
				break
			}
		}

		// No need to keep searching if we alrady found a duplicate name
		if duplicate_found {
			break
		}
	}

	// If returns error, binary could not be found
	bin_path, err := c.Output("sh", "-c", "command -v "+bin)
	if err != nil {
		return fmt.Errorf("Error: Could not find a binary with name `%s` in $PATH. Nothing to export.", bin)
	}

	// Binaries in ~/.local/bin are already accessible by the host
	if strings.Contains(bin_path, "/.local/bin") {
		msg := fmt.Sprintf("`%s` is already shared with host system. There's no need to export it.\n", bin)
		spinner.Info(msg)
		return nil
	}

	c.Run("sh", "-c", "distrobox-export --bin "+string(bin_path)+" --export-path ~/.local/bin >/dev/null 2>/dev/null || true")
	if bin_rename != "" {
		if err := c.Run("sh", "-c", "mv ~/.local/bin/"+bin+" ~/.local/bin/"+bin_rename); err != nil {
			return err
		}

		spinner.Success(fmt.Sprintf("Binary exported to `~/.local/bin/%s`.", bin_rename))
		return nil
	}

	spinner.Success(fmt.Sprintf("Binary exported to `~/.local/bin/%s`.", bin))
	return nil
}

func (c *Container) RemoveDesktopEntry(program string) error {
	container_name := c.GetContainerName()
	spinner, err := cmdr.Spinner.Start(fmt.Sprintf("Removing desktop entry: %v", program))
	if err != nil {
		return err
	}
	defer spinner.Stop()

	home_dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(home_dir + "/.local/share/applications")
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(strings.ToLower(file.Name()),
			strings.ToLower(container_name+"-"+program)) {
			err := os.Remove(home_dir + "/.local/share/applications/" + file.Name())
			if err != nil {
				return err
			}
			spinner.Success(fmt.Sprintf("Desktop entry for %v removed.", program))
			return nil
		}
	}

	spinner.Info(fmt.Sprintf("Desktop entry %v not found.", program))
	return nil
}

// RemoveBinary unexports/removes an exported binary application.
// fail_silently will not return an error if the file was not found or is invalid,
// this is useful for when removing packages, where the binary may not have been exported
// fail_silently will still return an error on unexpected issues
func (c *Container) RemoveBinary(bin string, fail_silently bool) error {
	// Check file exists in ~/.local/bin
	spinner, err := cmdr.Spinner.Start(fmt.Sprintf("Removing binary export: %v.", bin))
	if err != nil {
		return err
	}
	defer spinner.Stop()

	local_bin_file := fmt.Sprintf("/home/%s/.local/bin/%s", os.Getenv("USER"), bin)
	if _, err := os.Stat(local_bin_file); os.IsNotExist(err) {
		// Try to look for a prefixed file
		var prefix string
		switch c.containerType {
		case APT:
			prefix = fmt.Sprintf("apt_%s", bin)
		case AUR:
			prefix = fmt.Sprintf("aur_%s", bin)
		case DNF:
			prefix = fmt.Sprintf("dnf_%s", bin)
		case APK:
			prefix = fmt.Sprintf("apk_%s", bin)
		case ZYPPER:
			prefix = fmt.Sprintf("zypper_%s", bin)
		case XBPS:
			prefix = fmt.Sprintf("xbps_%s", bin)
		case SWUPD:
			prefix = fmt.Sprintf("swupd_%s", bin)
		default:
			return errors.New("Can't unexport binary from unknown container")
		}

		prefixed_bin_file := fmt.Sprintf("/home/%s/.local/bin/%s", os.Getenv("USER"), prefix)
		if _, prefix_err := os.Stat(prefixed_bin_file); os.IsNotExist(prefix_err) {
			if !fail_silently {
				spinner.Info(fmt.Sprintf("Binary `%s` is not exported.", bin))
			}

			return nil
		} else {
			local_bin_file = prefixed_bin_file
		}
	}

	// Ensure it's a distrobox export by reading the file's second line
	file, err := os.Open(local_bin_file)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for linenr := 1; scanner.Scan(); linenr++ {
		if linenr == 2 {
			text := scanner.Text()
			if text != "# distrobox_binary" {
				if !fail_silently {
					return fmt.Errorf("`~/.local/bin/%s` is not an apx export, refusing to remove.", bin)
				} else {
					return nil
				}
			} else {
				break
			}
		}
	}

	// Delete it
	err = os.Remove(local_bin_file)
	if err != nil {
		return err
	}

	spinner.Success(fmt.Sprintf("Binary %v unexported.", bin))

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

func ApplyForAll(command string, flags []string) error {
	for i := 0; i < CONTAINER_TYPES; i++ {
		container := NewContainer(ContainerType(i))
		if !container.Exists() {
			continue
		}

		name := container.GetContainerName()

		fmt.Println()
		log.Default().Println(fmt.Sprintf("Running %s in %s...", command, name))

		command := append([]string{}, container.GetPkgCommand(command)...)
		for _, flag := range flags {
			command = append(command, flag)
		}

		if err := container.Run(command...); err != nil {
			return err
		}
	}

	return nil
}
