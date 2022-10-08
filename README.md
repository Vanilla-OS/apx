<div align="center">
  <img src="apx-logo.png" height="120">
  <h1 align="center">apx</h1>
  <p align="center">Apx is the default package manager in Vanilla OS. It installs packages inside a managed container by default but can be used to install packages in the host as well.</p>
  <small>Special thanks to <a href="https://github.com/89luca89/distrobox">distrobox</a> for making this possible.</small>
</div>

<br/>

## Help
```
Special thanks to distrobox <https://github.com/89luca89/distrobox> for
  making this possible.

Usage: apx [options] [command] [arguments]

Options:
  -h, --help    Show this help message and exit
  -v, --version Show version and exit
  --sys         Perform operations on the system instead of the managed container

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

## Other distros
To use with another package manager, re-compile editing the `config.json` file
to point to the desired package manger and image.
