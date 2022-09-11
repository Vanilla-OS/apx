<div align="center">
  <img src="apx-logo.svg" height="64">
  <h1 align="center">apx</h1>
  <p align="center">Apx is a wrapper around apt to make it works inside a container from outside, directly on the host.</p>
  <small>Special thanks to <a href="https://github.com/89luca89/distrobox">distrobox</a> for making this possible.</small>
</div>

<br/>


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
