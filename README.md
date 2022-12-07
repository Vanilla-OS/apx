<div align="center">
  <img src="apx-logo.png" height="120">
  <h1 align="center">apx</h1>
  <p align="center">Apx is the default package manager in Vanilla OS. It installs packages inside a managed container by default but can be used to install packages in the host as well.</p>
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
    autoremove  Remove all unused packages automatically
    clean       Clean the apx package manager cache
    enter       Enter in the container shell
    export      Export/Recreate a program's desktop entry from a managed container
    help        Show this help message and exit
    init        Initialize the managed container
    install     Install packages inside a managed container
    list        List installed packages
    log         Show logs
    purge       Purge packages inside a managed container
    run         Run a command inside a managed container
    remove      Remove packages inside a managed container
    search      Search for packages in a managed container
    show        Show details about a package
    unexport    Unexport/Remove a program's desktop entry from a managed container
    update      Update the list of available packages
    upgrade     Upgrade the system by installing/upgrading available packages
    version     Show version and exit
```

## Other distros

> Please consider to keep the project name as `apx` to avoid confusion for users.

To use with another package manager, re-compile editing the `config.json` file
to point to the desired package manger and image.
