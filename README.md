<div align="center">
  <img src="apx-logo.png" height="120">
  <h1 align="center">apx</h1>
  <p align="center">Apx (/à·peks/) is the default package manager in Vanilla OS. It is a wrapper around multiple package managers to install packages and run commands inside a managed container.</p>
  <small>Special thanks to <a href="https://github.com/89luca89/distrobox">distrobox</a> for making this possible.</small>
</div>

<br/>

## Help

```
Apx is a package manager with support for multiple sources allowing you to install packages in a managed container.

Usage:
  apx [command]

Available Commands:
  autoremove  Remove all unused packages automatically
  clean       Clean the apx package manager cache
  completion  Generate the autocompletion script for the specified shell
  enter       Enter in the container shell
  export      Export/Recreate a program's desktop entry from a managed container
  help        Help about any command
  init        Initialize the managed container
  install     Install packages inside a managed container
  list        List installed packages.
  purge       Purge packages inside a managed container
  remove      Remove packages inside a managed container.
  run         Run a program inside a managed container.
  search      Search for packages in a managed container.
  show        Show details about a package
  unexport    Unexport/Remove a program's desktop entry from a managed container
  update      Update the list of available packages
  upgrade     Upgrade the system by installing/upgrading available packages.

Global Flags:
      --apk           Install packages from the Alpine repository.
      --apt           Install packages from the Ubuntu repository.
      --aur           Install packages from the AUR (Arch User Repository).
      --dnf           Install packages from the Fedora's DNF (Dandified YUM) repository.
      --zypper        Install packages from the openSUSE repository.
      --xbps          Install packages from the Void (Linux) repository.
  -h, --help          help for apx
  -n, --name string   Create or use custom container with this name.
  -v, --version       version for apx

Use "apx [command] --help" for more information about a command.
```

## Docs

The official **documentation and manpage** for `apx` are available at <https://documentation.vanillaos.org/docs/apx/>.

## Dependencies

To add new dependencies, use `go get` as usual, then run `go mod tidy` and finally `go mod vendor` before
committing code.

## Testing Translations locally

To test translations made in the `.yml` file locally, perform `go build` first in the correct directory then execute this command `LANG=<language_code> ./apx man > man/<language_code>/apx.1` (i.e `LANG=sv ./apx man > man/sv/apx.1`).

## Instructions for using Apx in other distributions

Apx has been designed in a distro-agnostic manner allowing it to work with any distribution. (Note: Nix command in Apx requires systemd)

### Prerequisites

- You must install `go` from your distribution's native repositories to compile `apx`.
- You can install either `git` or `gh` to clone the repository.

### Procedure

- Navigate to your directory to clone the repository using `cd`.
- Clone apx's repository using `git` or `gh`:-

```bash
git clone https://github.com/Vanilla-OS/apx.git
gh repo clone Vanilla-OS/apx
```

-  For the Apx binary to work with the terminal, you need to add the cloned location to your PATH using the following command:-

```bash
sudo cp apx /usr/bin
```

> In the above command, you can replace the path with: `/usr/bin/` or `~/.local/bin` if required.

- Create a directory to store Distrobox and configure Apx using the following steps:-

```bash
sudo mkdir /etc/apx
sudo cp config/config.json /etc/apx/
sudo mkdir /usr/lib/apx
```

- Now, we need to install the distrobox binary using the following steps:-

```bash
curl -s https://raw.githubusercontent.com/89luca89/distrobox/main/install | sh -s -- --prefix ~/.local
sudo mv ~/.local/share/distrobox* /usr/lib/apx
```

> Note:- Apx uses a fork of distrobox called `micro-distrobox`, but it is currently unavailable for other distributions, affecting the export of desktop entries.

- To fix exporting desktop entries, you need to perform the following step:-

```bash
sudo chown <username> ~/.local/share/icons -R
```
