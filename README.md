<div align="center">
  <img src="apx-logo.png" height="120">
  <h1 align="center">apx</h1>
  <p align="center">Apx (/à·peks/) is the default package manager in Vanilla OS. It is a wrapper around multiple package managers to install packages and run commands inside a managed container.</p>
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
  --apk         Install packages from the Alpine repository

Commands:
    autoremove  Remove all unused packages
    clean       Clean the apx package manager cache
    enter       Enter the container's shell
    export      Export/Recreate a program's desktop entry from a managed container
    help        Show this help message and exit
    init        Initialize a managed container
    install     Install packages inside the container
    list        List installed packages
    log         Show logs (This command is yet to be implemented)
    purge       Purge packages from the container
    run         Run a command inside the container
    remove      Remove packages from the container
    search      Search for packages
    show        Show details about a package
    unexport    Unexport/Remove a program's desktop entry
    update      Update the list of available packages
    upgrade     Upgrade the system by installing/upgrading available packages
```

## Docs

The official **documentation and manpage** for `apx` are available at <https://documentation.vanillaos.org/docs/apx/>.

## Other Distributions and Package Managers

> Please consider to keep the project name as `apx` to avoid confusion for users.

To use with another distro and/or package manager, re-compile editing the `config.json` file or simply edit `/etc/apx/config.json` on a pre-compiled installation
to point to the desired:
* Image:

by adding the docker image path for example `docker.io/library/ubuntu` to `"image": "",` under `"container": {`



where for a **Simple Ubuntu image setup** it would look something like:
`"image": "docker.io/library/ubuntu",`
* Package Manager:

by replacing `"bin": "/usr/bin/apt",` with `"bin": "ABSOLUTE PATH TO YOUR PACKAGE MANAGER",`
and replacing all `"cmd*":"",` with the apt alternative fot your package manager.


where for a **Simple Fedora dnf setup** it would look something like:
```
    "pkgmanager": { 
         "bin": "/usr/bin/dnf", 
         "sudo": true, 
         "cmdAutoremove": "autoremove", 
         "cmdClean": "clean", 
         "cmdInstall": "install", 
         "cmdList": "list", 
         "cmdPurge": "remove", 
         "cmdRemove": "remove", 
         "cmdSearch": "search", 
         "cmdShow": "repoquery -i", 
         "cmdUpdate": "update --refresh", 
         "cmdUpgrade": "distrosync" 
     }
```
