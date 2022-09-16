<div align="center">
  <img src="apx-logo.svg" height="64">
  <h1 align="center">apx</h1>
  <p align="center">Apx is a wrapper around apt to make it works inside a container from outside, directly on the host.</p>
  <small>Special thanks to <a href="https://github.com/89luca89/distrobox">distrobox</a> for making this possible.</small>
</div>

<br/>

## Help
```
Apx is a wrapper around apt to make it works inside a container
  from outside, directly on the host.

Special thanks to distrobox <https://github.com/89luca89/distrobox> for
  making this possible.

Usage: apx [options] [command]

Options:
  -h, --help    Show this help message and exit
  -v, --version Show version and exit

Commands:
    autoremove  Remove automatically all unused packages
    clean       Clean the apt cache
    enter       Enter the container
    help        Show this help message and exit
    init        Initialize the container
    install     Install packages
    list        List packages based on package names
    log         Show logs
    purge       Purge packages
    run         Run a command inside the container
    remove      Remove packages
    search      Search in package descriptions
    show        Show package details
    update      Update list of available packages
    upgrade     Upgrade the system by installing/upgrading packages
    version     Show version and exit
```

## Notes for packagers
If you want to package apx for your distribution without using an Ubuntu image, you
need to configure your own package manager interface and define the container image
name. To do so, you need to create a file named `packager-rules.cfg` in `/etc/apx/`
with the following content:

```bash
PKG_CONTAINER_IMAGE="your-image-name" # e.g. "archlinux:latest"
PKG_MANAGER_NEED_SUDO="true|false" # If your package manager needs sudo
PKG_MANAGER_INTERFACE="your-package-manager-interface" # e.g. pacman
```

now to write your own package manager interface, you need to create a file named
`packager-manager-interface.cfg` in `/etc/apx/` and update it following the
[packager-manager-interface.cfg.example](packager-manager-interface.cfg.example)
file.
