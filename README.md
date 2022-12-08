<div align="center">
  <img src="apx-logo.png" height="120">
  <h1 align="center">apx</h1>
  <p align="center">Apx is the default package manager in Vanilla OS. It is a wrapper around multiple package managers to install packages and run commands inside a managed container.</p>
  <small>Special thanks to <a href="https://github.com/89luca89/distrobox">distrobox</a> for making this possible.</small>
</div>

<br/>

## Help

```
Usage: apx [options] [command] [arguments]

Options:
  -h, --help    Show this help message and exit
  -v, --version Show version and exit
  --aur         Install packages from the AUR repository
  --dnf         Install packages from the Fedora repository

Commands:
    autoremove  Remove all unused packages
    clean       Clean the apx package manager cache
    enter       Enter the container's shell
    export      Export/Recreate a program's desktop entry from a managed container
    help        Show this help message and exit
    init        Initialize a managed container
    install     Install packages inside the container
    list        List installed packages
    log         Show logs
    purge       Purge packages from the container
    run         Run a command inside the container
    remove      Remove packages from the container
    search      Search for packages
    show        Show details about a package
    unexport    Unexport/Remove a program's desktop entry
    update      Update the list of available packages
    upgrade     Upgrade the system by installing/upgrading available packages
    version     Show version and exit
```

## Other distros

> Please consider to keep the project name as `apx` to avoid confusion for users.

To use with another package manager, re-compile editing the `config.json` file
to point to the desired package manger and image.
