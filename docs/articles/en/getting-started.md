---
Title: Getting Started
Description: How to start using Apx on your system.
PublicationDate: 2024-10-18
Listed: true
Authors: 
  - jardon
Tags:
  - getting-started
---

Apx is designed to be a versatile, distro-agnostic tool that can be installed on any Linux distribution. Follow the steps below to get Apx up and running on your system.

## Prerequisites

Before installing Apx, ensure you have the following software installed:

1. **Go**: This is required to compile Apx. Install it from your distribution's package manager.
2. **Git**: Needed to clone the Apx repository.
3. **Podman or Docker**: Either container runtime is suitable, but Podman is recommended.
4. **Make**: This utility is used for building and installing Apx.

## Installation Procedure

### Clone the Apx Repository

Open your terminal and run the following commands to clone the repository and navigate into it:

```bash
git clone --recursive https://github.com/Vanilla-OS/apx.git
cd apx
```

### Build Apx

Compile Apx by executing:

```bash
make build
```

### Install Apx

To install Apx system-wide, run:

```bash
sudo make install
```

### Install Apx Manpages

For the manual pages, execute:

```bash
sudo make install-manpages
```

## Custom Installation Destination

You can change the installation prefix or destination using `PREFIX` and `DESTDIR`. Here are examples for custom installations:

- **Install Apx to `~/.local`**:

  ```bash
  make install PREFIX=$HOME/.local
  make install-manpages PREFIX=$HOME/.local
  ```

- **Install Apx to a separate root**:
  ```bash
  make install DESTDIR=$HOME/altroot
  make install-manpages DESTDIR=$HOME/altroot
  ```

# Getting Started with Apx-GUI

Apx-GUI provides a graphical interface for managing your Apx installations. Follow the steps below to install Apx-GUI on your system.

## Dependencies

Before you start, ensure you have the following dependencies installed:

1. **build-essential**: A package containing essential compilation tools.
2. **meson**: A build system to configure the project.
3. **libadwaita-1-dev**: A library for building modern GNOME applications.
4. **gettext**: A utility for internationalization and localization.
5. **desktop-file-utils**: Tools for handling desktop entry files.
6. **apx (2.0+)**: Ensure you have Apx version 2.0 or higher installed.

You can install these dependencies using your distribution's package manager. For example, on Debian-based systems, you can run:

```bash
sudo apt update
sudo apt install build-essential meson libadwaita-1-dev gettext desktop-file-utils apx
```

## Installation Procedure

### Clone the Apx-GUI Repository

Open your terminal and run the following commands to clone the repository and navigate into it:

```bash
git clone https://github.com/Vanilla-OS/apx-gui.git
cd apx-gui
```

### Build Apx-GUI

Once you have cloned the repository, build Apx-GUI by running:

```bash
meson setup build
ninja -C build
```

> **NOTE:** you can set a custom installation destination by passing `--prefix=/path/to/dir` to `meson`

### Install Apx-GUI

After successfully building the application, install it with the following command:

```bash
sudo ninja -C build install
```

### Run Apx-GUI

You can launch Apx-GUI using the following command:

```bash
apx-gui
```

Follow these steps to successfully install and run Apx-GUI on your system. Enjoy managing your Apx installations with the graphical interface!
