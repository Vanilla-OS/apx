<div align="center">
  <img src="apx-logo.png" height="120">
  <h1 align="center">Apx</h1>

[![Translation Status][weblate-image]][weblate-url]
[![build result][build-image]][build-url]

[weblate-url]: https://hosted.weblate.org/engage/vanilla-os/
[weblate-image]: https://hosted.weblate.org/widgets/vanilla-os/-/apx/svg-badge.svg
[weblate-status-image]: https://hosted.weblate.org/widgets/vanilla-os/-/apx/multi-auto.svg
[repology-url]: https://repology.org/project/apx-package-manager/versions
[repology-image]: https://repology.org/badge/vertical-allrepos/apx-package-manager.svg
[build-image]: https://build.opensuse.org/projects/home:fabricators:orchid/packages/apx/badge.svg?type=default
[build-url]: https://build.opensuse.org/package/show/home:fabricators:orchid/apx

  <p align="center">Apx (/à·peks/) is the default package manager in Vanilla OS. It is a wrapper around multiple package managers to install packages and run commands inside a managed container.</p>
  <small>Special thanks to <a href="https://github.com/89luca89/distrobox">distrobox</a> for making this possible.</small>
</div>

<br/>

## Help

```
Apx is a package manager with support for multiple sources,
allowing you to install packages in a managed container.

Usage:
  apx [command]

Managed Container Commands
  autoremove  Remove all unused packages automatically
  clean       Clean the apx package manager cache
  enter       Enter a shell in the managed container
  export      Export/Recreate a program's desktop entry from a managed container
  init        Initialize a managed container
  install     Install packages inside a managed container.
  list        List installed packages.
  purge       Purge packages inside a managed container
  remove      Remove packages inside a managed container.
  run         Run a program inside a managed container.
  search      Search for packages in a managed container.
  show        Show details about a package
  unexport    Unexport/Remove a program's desktop entry from a managed container
  update      Update the list of available packages
  upgrade     Upgrade the system by installing/upgrading available packages.

Additional Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -v, --verbose       show more detailed output
      --apt           Install packages from the Ubuntu repository.
      --aur           Install packages from the AUR (Arch User Repository).
      --dnf           Install packages from the Fedora's DNF (Dandified YUM) repository.
      --apk           Install packages from the Alpine repository.
      --zypper        Install packages from the OpenSUSE repository.
      --xbps          Install packages from the Void (Linux) repository.
      --nix           Install packages from the Nixpkgs (Nix packages) repository.
  -n, --name string   Apply to custom container with this name.
  -h, --help          help for apx
      --version       version for apx

Use "apx [command] --help" for more information about a command.
```

## Documentation and Guides

### Documentation

The official **documentation and manpage** for `apx` are available at <https://documentation.vanillaos.org/docs/apx/>.

### Guides

A guide for Installing applications in `apx` is available at <https://handbook.vanillaos.org/2023/01/11/install-and-manage-applications.html>.

## Dependencies

To add new dependencies, use `go get` as usual, then run `go mod tidy` and `go mod vendor` before committing the code.

## Translations

Contribute translations for the manpage and help page in [Weblate](https://hosted.weblate.org/projects/vanilla-os/apx).

[![Translation Status][weblate-status-image]][weblate-url]

### Generating man pages for translations

Once the translation is complete in Weblate and the changes committed, clone the repository using `git` and perform `go build`, create a directory using the `mkdir man/<language_code> && mkdir man/<language_code>/man1` command, and execute this command `LANG=<language_code> ./apx man > man/<language_code>/man1/apx.1`. Open a PR for the generated manpage here.

## Instructions for using Apx in other distributions

Apx has been designed as a distro-agnostic tool, allowing it to work with any distribution. (Note: The Nix integration in Apx requires SystemD)

### Prerequisites

- You must have `go` installed from your distribution's native repositories to compile `apx`.
- You must have `git` installed to clone the repository.
- You must have either `podman` or `docker` installed.
- You must have `make` installed for the installation.

### Procedure

Clone apx's repository using `git` and enter it using `cd`:-

``` bash
git clone --recursive https://github.com/Vanilla-OS/apx.git
cd apx
```

Build apx using `make`:-

``` bash
make build
```

Install apx using `make`:-

``` bash
sudo make install
```

Install the apx manpages using `make`:-

``` bash
sudo make install-manpages
```

Uninstall apx using `make`:-

```bash
make uninstall
```

Uninstall apx man pages using `make`:-

```bash
make uninstall-manpages
```

### Installing apx in a custom destination

The prefix or installation destination can be changed using `PREFIX` and `DESTDIR`.

To install apx into `~/.local` for example:-

``` bash
make install PREFIX=$HOME/.local
make install-manpages PREFIX=$HOME/.local
```

Or into a separate root:-

``` bash
make install DESTDIR=$HOME/altroot
make install-manpages DESTDIR=$HOME/altroot
```

## Community Packages

Apx is packaged in various sources by our community, if you aren't comfortable with building `apx` manually you can install a package listed below.

(**Note:** These packages are unofficial if there are any issues with packaging report them to the respective maintainers, if there are any issues with apx functionalities report them here.)

[![Packaging status][repology-image]][repology-url]
