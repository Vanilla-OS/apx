<div align="center">
  <img src="apx-logo.png" height="120">
  <h1 align="center">Apx</h1>

[![Translation Status][weblate-image]][weblate-url]

[weblate-url]: https://hosted.weblate.org/engage/vanilla-os/
[weblate-image]: https://hosted.weblate.org/widgets/vanilla-os/-/apx/svg-badge.svg
[repology-url]: https://repology.org/project/apx-package-manager/versions
[repology-image]: https://repology.org/badge/vertical-allrepos/apx-package-manager.svg

  <p align="center">Apx (/à·peks/) is the default package manager in Vanilla OS. It is a wrapper around multiple package managers to install packages and run commands inside a managed container.</p>
  <small>Special thanks to <a href="https://github.com/89luca89/distrobox">distrobox</a> for making this possible.</small>
</div>

<br/>

## Help

```
Apx is a package manager with support for multiple sources, allowing you to install packages in subsystems.

Usage:
  apx [command]

Available Commands:
  [subsystem] Work with the specified subsystem, accessing the package manager and environment.
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  pkgmanagers Work with the package managers that are available in apx.
  stacks      Work with the stacks that are available in apx.
  subsystems  Work with the subsystems that are available in apx.

Flags:
  -h, --help      help for apx
  -v, --version   version for apx

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

### Generating man pages for translations

Once the translation is complete in Weblate and the changes committed, clone the repository using `git` and perform `go build`, create a directory using the `mkdir man/<language_code> && mkdir man/<language_code>/man1` command, and execute this command `LANG=<language_code> ./apx man > man/<language_code>/man1/apx.1`. Open a PR for the generated manpage here.

## Instructions for using Apx in other distributions

Apx has been designed as a distro-agnostic tool, allowing it to work with any distribution.

### Prerequisites

- You must have `go` installed from your distribution's native repositories to compile `apx`.
- You must have `git` installed to clone the repository.
- You must have either `podman` (suggested) or `docker` installed.
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
