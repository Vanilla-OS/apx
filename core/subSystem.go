package core

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Vanilla OS Contributors <https://github.com/vanilla-os/>
	Copyright: 2024
	Description:
		Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

type SubSystem struct {
	InternalName         string
	Name                 string
	Stack                *Stack
	Home                 string
	Status               string
	ExportedPrograms     map[string]map[string]string
	HasInit              bool
	IsManaged            bool
	IsRootfull           bool
	IsUnshared           bool
	HasNvidiaIntegration bool
	Hostname             string
}

func NewSubSystem(name string, stack *Stack, home string, hasInit bool, isManaged bool, isRootfull bool, isUnshared bool, hasNvidiaIntegration bool, hostname string) (*SubSystem, error) {
	internalName := genInternalName(name)
	return &SubSystem{
		InternalName:         internalName,
		Name:                 name,
		Stack:                stack,
		Home:                 home,
		HasInit:              hasInit,
		IsManaged:            isManaged,
		IsRootfull:           isRootfull,
		IsUnshared:           isUnshared,
		HasNvidiaIntegration: hasNvidiaIntegration,
		Hostname:             hostname,
	}, nil
}

func genInternalName(name string) string {
	return fmt.Sprintf("apx-%s", strings.ReplaceAll(strings.ToLower(name), " ", "-"))
}

func findExportedBinaries(internalName string) map[string]map[string]string {
	home, err := os.UserHomeDir()
	if err != nil {
		return map[string]map[string]string{}
	}
	binPath := filepath.Join(home, ".local", "bin")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		return map[string]map[string]string{}
	}
	binDir := os.DirFS(binPath)

	binaries := map[string]map[string]string{}
	err = fs.WalkDir(binDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		file, err := os.Open(filepath.Join(binPath, path))
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		const maxTokenSize = 1024 * 2048
		buf := make([]byte, maxTokenSize)
		scanner.Buffer(buf, maxTokenSize)
		for scanner.Scan() {
			if scanner.Text() == "# distrobox_binary" {
				scanner.Scan()
				if strings.HasSuffix(scanner.Text(), internalName) {
					name := filepath.Base(path)
					binaries[name] = map[string]string{
						"Exec": path,
						// "Icon":        pIcon,
						"Name": name,
						// "GenericName": pGenericName,
					}
				}
				break
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Printf("error reading binaries: %s\n", err)
		return map[string]map[string]string{}
	}

	return binaries
}

func findExportedPrograms(internalName string, name string) map[string]map[string]string {
	home, err := os.UserHomeDir()
	if err != nil {
		return map[string]map[string]string{}
	}

	files, err := filepath.Glob(fmt.Sprintf("%s/.local/share/applications/%s-*.desktop", home, internalName))
	if err != nil {
		return map[string]map[string]string{}
	}

	programs := map[string]map[string]string{}
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return map[string]map[string]string{}
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			continue
		}

		pName := ""
		pExec := ""
		pIcon := ""
		pGenericName := ""
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "Name=") {
				pName = strings.TrimPrefix(line, "Name=")
				pName = strings.ReplaceAll(pName, fmt.Sprintf(" on %s", name), "")
			}

			if strings.HasPrefix(line, "Exec=") {
				pExec = strings.TrimPrefix(line, "Exec=")
			}

			if strings.HasPrefix(line, "Icon=") {
				pIcon = strings.TrimPrefix(line, "Icon=")
			}

			if strings.HasPrefix(line, "GenericName=") {
				pGenericName = strings.TrimPrefix(line, "GenericName=")
			}
		}

		if pName != "" && pExec != "" {
			programs[pName] = map[string]string{
				"Exec":        pExec,
				"Icon":        pIcon,
				"Name":        pName,
				"GenericName": pGenericName,
			}
		}
	}

	return programs
}

func findExported(internalName string, name string) map[string]map[string]string {
	bins := findExportedBinaries(internalName)
	progs := findExportedPrograms(internalName, name)

	// If duplicate is found, give priority to application
	for k, v := range progs {
		bins[k] = v
	}

	return bins
}

func (s *SubSystem) Create() error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	labels := map[string]string{
		"stack": strings.ReplaceAll(s.Stack.Name, " ", "\\ "),
		"name":  strings.ReplaceAll(s.Name, " ", "\\ "),
	}

	if s.IsManaged {
		labels["managed"] = "true"
	}

	if s.HasInit {
		labels["hasInit"] = "true"
	}

	if s.IsUnshared {
		labels["unshared"] = "true"
	}

	if s.HasNvidiaIntegration {
		labels["nvidia"] = "true"
	}

	err = dbox.CreateContainer(
		s.InternalName,
		s.Stack.Base,
		s.Stack.Packages,
		s.Home,
		labels,
		s.HasInit,
		s.IsRootfull,
		s.IsUnshared,
		s.HasNvidiaIntegration,
		s.Hostname,
	)
	if err != nil {
		return err
	}

	return nil
}

func LoadSubSystem(name string, isRootFull bool) (*SubSystem, error) {
	dbox, err := NewDbox()
	if err != nil {
		return nil, err
	}

	internalName := genInternalName(name)
	container, err := dbox.GetContainer(internalName, isRootFull)
	if err != nil {
		return nil, err
	}

	stack, err := LoadStack(container.Labels["stack"])
	if err != nil {
		return nil, err
	}
	return &SubSystem{
		InternalName: internalName,
		Name:         container.Labels["name"],
		Stack:        stack,
		Status:       container.Status,
		HasInit:      container.Labels["hasInit"] == "true",
		IsManaged:    container.Labels["managed"] == "true",
		IsRootfull:   isRootFull,
		IsUnshared:   container.Labels["unshared"] == "true",
	}, nil
}

func ListSubSystems(includeManaged bool, includeRootFull bool) ([]*SubSystem, error) {
	dbox, err := NewDbox()
	if err != nil {
		return nil, err
	}

	containers, err := dbox.ListContainers(includeRootFull)
	if err != nil {
		return nil, err
	}

	subsystems := []*SubSystem{}
	for _, container := range containers {
		if _, ok := container.Labels["name"]; !ok {
			// log.Printf("Container %s has no name label, skipping", container.Name)
			continue
		}

		if !includeManaged {
			if _, ok := container.Labels["managed"]; ok {
				continue
			}
		}

		stack, err := LoadStack(container.Labels["stack"])
		if err != nil {
			log.Printf("Error loading stack %s: %s", container.Labels["stack"], err)
			continue
		}

		internalName := genInternalName(container.Labels["name"])
		subsystem := &SubSystem{
			InternalName:     internalName,
			Name:             container.Labels["name"],
			Stack:            stack,
			Status:           container.Status,
			ExportedPrograms: findExported(internalName, container.Labels["name"]),
		}

		subsystems = append(subsystems, subsystem)
	}

	return subsystems, nil
}

// ListSubsystemForStack returns a list of subsystems for the specified stack.
func ListSubsystemForStack(stackName string) ([]*SubSystem, error) {
	dbox, err := NewDbox()
	if err != nil {
		return nil, err
	}

	containers, err := dbox.ListContainers(true)
	if err != nil {
		return nil, err
	}

	subsystems := []*SubSystem{}
	for _, container := range containers {
		if _, ok := container.Labels["name"]; !ok {
			continue
		}

		stack, err := LoadStack(stackName)
		if err != nil {
			log.Printf("Error loading stack %s: %s", stackName, err)
			continue
		}

		internalName := genInternalName(container.Labels["name"])
		subsystem := &SubSystem{
			InternalName:     internalName,
			Name:             container.Labels["name"],
			Stack:            stack,
			Status:           container.Status,
			ExportedPrograms: findExported(internalName, container.Labels["name"]),
		}

		if subsystem.Stack.Name == stack.Name {
			subsystems = append(subsystems, subsystem)
		}
	}

	return subsystems, nil
}

func (s *SubSystem) Exec(captureOutput, detachedMode bool, args ...string) (string, error) {
	dbox, err := NewDbox()
	if err != nil {
		return "", err
	}

	out, err := dbox.ContainerExec(s.InternalName, captureOutput, false, s.IsRootfull, detachedMode, args...)
	if err != nil {
		return "", err
	}

	if captureOutput {
		return out, nil
	}

	return "", nil
}

func (s *SubSystem) Enter() error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}
	return dbox.ContainerEnter(s.InternalName, s.IsRootfull)
}

func (s *SubSystem) Start() error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}
	return dbox.ContainerStart(s.InternalName, s.IsRootfull)
}

func (s *SubSystem) Stop() error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}
	return dbox.ContainerStop(s.InternalName, s.IsRootfull)
}

func (s *SubSystem) Remove() error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	return dbox.ContainerDelete(s.InternalName, s.IsRootfull)
}

func (s *SubSystem) Reset() error {
	err := s.Remove()
	if err != nil {
		return err
	}

	return s.Create()
}

func (s *SubSystem) ExportDesktopEntry(appName string) error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	return dbox.ContainerExportDesktopEntry(s.InternalName, appName, fmt.Sprintf("on %s", s.Name), s.IsRootfull)
}

func (s *SubSystem) ExportDesktopEntries(args ...string) (int, error) {
	exportedN := 0

	for _, appName := range args {
		err := s.ExportDesktopEntry(appName)
		if err != nil {
			return exportedN, err
		}

		exportedN++
	}

	return exportedN, nil
}

func (s *SubSystem) UnexportDesktopEntries(args ...string) (int, error) {
	exportedN := 0

	for _, appName := range args {
		err := s.UnexportDesktopEntry(appName)
		if err != nil {
			return exportedN, err
		}

		exportedN++
	}

	return exportedN, nil
}

func (s *SubSystem) ExportBin(binary string, exportPath string) error {
	if !strings.HasPrefix(binary, "/") {
		binaryPath, err := s.Exec(true, false, "which", binary)
		if err != nil {
			return err
		}

		binary = strings.TrimSpace(binaryPath)
	}

	binaryName := filepath.Base(binary)

	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if exportPath == "" {
		exportPath = filepath.Join(homeDir, ".local", "bin")
	}

	joinedPath := filepath.Join(exportPath, binaryName)
	if _, err := os.Stat(joinedPath); err == nil {
		tmpExportPath := fmt.Sprintf("/tmp/%s", uuid.New().String())
		err = os.MkdirAll(tmpExportPath, 0o755)
		if err != nil {
			return err
		}

		err = dbox.ContainerExportBin(s.InternalName, binary, tmpExportPath, s.IsRootfull)
		if err != nil {
			return err
		}

		err = CopyFile(filepath.Join(tmpExportPath, binaryName), filepath.Join(exportPath, fmt.Sprintf("%s-%s", binaryName, s.InternalName)))
		if err != nil {
			return err
		}

		err = os.RemoveAll(tmpExportPath)
		if err != nil {
			return err
		}

		err = os.Chmod(filepath.Join(exportPath, fmt.Sprintf("%s-%s", binaryName, s.InternalName)), 0o755)
		if err != nil {
			return err
		}

		return nil
	}

	err = os.MkdirAll(exportPath, 0o755)
	if err != nil {
		return err
	}

	err = dbox.ContainerExportBin(s.InternalName, binary, exportPath, s.IsRootfull)
	if err != nil {
		return err
	}

	return nil
}

func (s *SubSystem) UnexportDesktopEntry(appName string) error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	return dbox.ContainerUnexportDesktopEntry(s.InternalName, appName, s.IsRootfull)
}

func (s *SubSystem) UnexportBin(binary string, exportPath string) error {
	if !strings.HasPrefix(binary, "/") {
		binaryPath, err := s.Exec(true, false, "which", binary)
		if err != nil {
			return err
		}

		binary = strings.TrimSpace(binaryPath)
	}

	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	return dbox.ContainerUnexportBin(s.InternalName, binary, s.IsRootfull)
}
