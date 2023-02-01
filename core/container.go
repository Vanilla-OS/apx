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
	"time"

	"github.com/briandowns/spinner"
	"github.com/vanilla-os/apx/settings"
)

type Container struct {
	containerType settings.DistroInfo
	customName    string
}

func NewContainer(kind settings.DistroInfo) *Container {
	return &Container{
		containerType: kind,
	}
}

func NewNamedContainer(kind settings.DistroInfo, name string) *Container {
	return &Container{
		containerType: kind,
		customName:    name,
	}
}

func (c *Container) GetContainerName() (name string) {
	var cn strings.Builder
	cn.WriteString(c.containerType.ContainerName)

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

	container_name := c.GetContainerName()
	spinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	spinner.Suffix = " Creating container..."

	spinner.Start()

	cmd := exec.Command(settings.Cnf.DistroboxPath, "create",
		"--name", container_name,
		"--image", c.containerType.Image,
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
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("error creating container: %v", err)
	}

	if c.containerType.Pkgmanager == settings.Yay {
		// Setup locales
		spinner.Suffix = " Setting up locales..."
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
		spinner.Suffix = " Downloading and installing Yay..."
		DownloadYay()
		c.Output(GetAurPkgCommand("install-yay-deps")...)
		c.Run(GetAurPkgCommand("install-yay")...)
	}

	spinner.Stop()

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
	out, err := c.Output("sh", "-c", "distrobox-host-exec $(readlink -fn $(getent passwd $USER | cut -f 7 -d :)) -l -i -c printenv | grep -E ^PATH=")
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to execute printenv: %s", err))
	}

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
			return errors.New(fmt.Sprintf("Could not read directory %s: %s", path, err))
		}

		duplicate_found := false
		for _, entry := range entries {
			// If duplicate is located in ~/.local/bin, we'll handle it later
			if entry.Name() == bin && !strings.Contains(path, "/.local/bin") {
				switch c.containerType.Pkgmanager {
				case settings.Apt:
					bin_rename = fmt.Sprintf("apt_%s", bin)
				case settings.Yay:
					bin_rename = fmt.Sprintf("aur_%s", bin)
				case settings.Dnf:
					bin_rename = fmt.Sprintf("dnf_%s", bin)
				case settings.Apk:
					bin_rename = fmt.Sprintf("apk_%s", bin)
				case settings.Zypper:
					bin_rename = fmt.Sprintf("zypper_%s", bin)
				case settings.Xbps:
					bin_rename = fmt.Sprintf("xbps_%s", bin)
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

	// If returns error, binary could not be found
	bin_path, err := c.Output("sh", "-c", "command -v "+bin)
	if err != nil {
		return errors.New(fmt.Sprintf("Error: Could not find a binary with name `%s` in $PATH. Nothing to export.", bin))
	}

	// Binaries in ~/.local/bin are already accessible by the host
	if strings.Contains(bin_path, "/.local/bin") {
		fmt.Printf("`%s` is already shared with host system. There's no need to export it.\n", bin)
		return nil
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

// RemoveBinary unexports/removes an exported binary application.
// fail_silently will not return an error if the file was not found or is invalid,
// this is useful for when removing packages, where the binary may not have been exported
// fail_silently will still return an error on unexpected issues
func (c *Container) RemoveBinary(bin string, fail_silently bool) error {
	// Check file exists in ~/.local/bin
	local_bin_file := fmt.Sprintf("/home/%s/.local/bin/%s", os.Getenv("USER"), bin)
	if _, err := os.Stat(local_bin_file); os.IsNotExist(err) {
		// Try to look for a prefixed file
		var prefix string
		switch c.containerType.Pkgmanager {
		case settings.Apt:
			prefix = fmt.Sprintf("apt_%s", bin)
		case settings.Yay:
			prefix = fmt.Sprintf("aur_%s", bin)
		case settings.Dnf:
			prefix = fmt.Sprintf("dnf_%s", bin)
		case settings.Apk:
			prefix = fmt.Sprintf("apk_%s", bin)
		case settings.Zypper:
			prefix = fmt.Sprintf("zypper_%s", bin)
		case settings.Xbps:
			prefix = fmt.Sprintf("xbps_%s", bin)
		default:
			return errors.New("can't unexport binary from unknown container")
		}

		prefixed_bin_file := fmt.Sprintf("/home/%s/.local/bin/%s", os.Getenv("USER"), prefix)
		if _, prefix_err := os.Stat(prefixed_bin_file); os.IsNotExist(prefix_err) {
			if !fail_silently {
				return errors.New(fmt.Sprintf("Binary `%s` is not exported.", bin))
			} else {
				return nil
			}
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
					return errors.New(fmt.Sprintf("`~/.local/bin/%s` is not an apx export, refusing to remove.", bin))
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
	for i := 0; i < 6; i++ {
		var distro settings.DistroInfo
		switch i {
		case 0:
			distro = settings.DistroUbuntu
		case 1:
			distro = settings.DistroArch
		case 2:
			distro = settings.DistroFedora
		case 3:
			distro = settings.DistroAlpine
		case 4:
			distro = settings.DistroOpensuse
		case 5:
			distro = settings.DistroVoid
		}

		container := NewContainer(distro)
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
