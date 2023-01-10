package settings

import "errors"

type PackageManager string

const (
	Apt = "apt"
	Yay = "yay"
	Dnf = "dnf"
	Apk = "apk"
)

type DistroInfo struct {
	Id            string
	Image         string
	Pkgmanager    PackageManager
	ContainerName string
}

func FromPackageManger(pkgmanager string) (DistroInfo, error) {
	switch pkgmanager {
	case Apt:
		return DistroUbuntu, nil
	case Yay:
		return DistroArch, nil
	case Dnf:
		return DistroFedora, nil
	case Apk:
		return DistroAlpine, nil
	default:
		return DistroInfo{}, errors.New("Invalid package manager")
	}
}

var DistroUbuntu = DistroInfo{Id: "ubuntu", Image: "docker.io/library/ubuntu", Pkgmanager: Apt, ContainerName: "apx_managed_ubuntu"}
var DistroArch = DistroInfo{Id: "arch", Image: "docker.io/library/archlinux", Pkgmanager: Yay, ContainerName: "apx_managed_archlinux"}
var DistroFedora = DistroInfo{Id: "fedora", Image: "docker.io/library/fedora", Pkgmanager: Dnf, ContainerName: "apx_managed_fedora"}
var DistroAlpine = DistroInfo{Id: "alpine", Image: "docker.io/library/alpine", Pkgmanager: Apk, ContainerName: "apx_managed_alpine"}
