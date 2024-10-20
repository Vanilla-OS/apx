---
Title: Working with Stacks
Description: Learn how to create, manage, and customize stacks in Apx to simplify the management of applications and dependencies.
PublicationDate: 2024-10-18
Listed: true
Authors:
  - jardon
Tags:
  - working
  - stacks
---

Stacks are preconfigured operating system images in APX that can be reused across different subsystems. This allows for consistent environments and simplifies the management of applications and dependencies.

To see the available commands for managing stacks, use the following command:

```bash
apx stacks --help
```

```
Work with the stacks that are available in apx.

Usage:
  apx stacks [command]

Available Commands:
  export      Export the specified stack.
  import      Import the specified stack.
  list        List all available stacks.
  new         Create a new stack.
  rm          Remove the specified stack.
  show        Show information about the specified stack.
  update      Update the specified stack.

Flags:
  -h, --help   help for stacks

Use "apx stacks [command] --help" for more information about a command.
```

## Creating a Stack

Creating a stack allows you to define a customized operating environment. APX ships with several default stacks, but you can create user-defined stacks tailored to your needs.

To see the options available for creating a new stack, run:

```bash
apx stacks new --help
```

```
Create a new stack.

Usage:
  apx stacks new [flags]

Flags:
  -b, --base string          The base distribution image to use. (For a list of compatible images view: https://distrobox.it/compatibility/#containers-distros)
  -h, --help                 help for new
  -n, --name string          The name of the stack.
  -y, --no-prompt            Assume defaults to all prompts.
  -p, --packages string      The packages to install.
  -k, --pkg-manager string   The package manager to use.
```

> **NOTE:** Depending on the stack that needs to be added, it may be necessary to create a package manager first.

We want to use the latest Ubuntu LTS which at the time of writing is 24.04 (Noble Numbat) in a subsystem. Let's begin by listing off the currently available stacks.

```bash
apx stacks list
```

```
 INFO  Found 7 stacks:
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ NAME        ┊ BASE                                                     ┊ BUILT-IN ┊ PKGS ┊ PKG MANAGER ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ alpine      ┊ alpine:latest                                            ┊ yes      ┊ 0    ┊ apk         ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ arch        ┊ ghcr.io/archlinux/archlinux:multilib-devel               ┊ yes      ┊ 0    ┊ pacman      ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ fedora      ┊ fedora:latest                                            ┊ yes      ┊ 0    ┊ dnf         ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ opensuse    ┊ registry.opensuse.org/opensuse/tumbleweed:latest         ┊ yes      ┊ 0    ┊ zypper      ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ ubuntu      ┊ ubuntu:latest                                            ┊ yes      ┊ 0    ┊ apt         ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ vanilla-dev ┊ ghcr.io/vanilla-os/dev:main                              ┊ yes      ┊ 0    ┊ apt         ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ vanilla     ┊ ghcr.io/vanilla-os/pico:main                             ┊ yes      ┊ 0    ┊ apt         ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
```

Looking at the above output we see that there is no stack for Ubuntu 24.04. This means that we need to create one! To do this, we tell `apx` to create a new stack and pass it some info.

```bash
apx stacks new
```

> **NOTE:** These options can be passed as CLI args as shown above in the help output.

```
 INFO  Choose a name:
noble
 INFO  Choose a base (e.g. 'vanillaos/pico'):
ubuntu:noble
 INFO  Choose a package manager:
1. yum
2. apk
3. apt
4. dnf
5. pacman
6. zypper
 INFO  Select a package manager [1-6]:
3
 INFO  You have not provided any packages to install in the stack. Do you want to add some now?[y/N]
y
 INFO  Please type the packages you want to install in the stack, separated by a space:
neofetch vim
 SUCCESS  Created stack 'noble'.
```

We have received a "SUCCESS" message from `apx`! To confirm that the new `noble` stack has been created, just list off the available stacks.

```bash
apx stacks list
```

```
 INFO  Found 8 stacks:
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ NAME        ┊ BASE                                                     ┊ BUILT-IN ┊ PKGS ┊ PKG MANAGER ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ noble       ┊ ubuntu:noble                                             ┊ no       ┊ 2    ┊ apt         ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ alpine      ┊ alpine:latest                                            ┊ yes      ┊ 0    ┊ apk         ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ arch        ┊ ghcr.io/archlinux/archlinux:multilib-devel               ┊ yes      ┊ 0    ┊ pacman      ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ fedora      ┊ fedora:latest                                            ┊ yes      ┊ 0    ┊ dnf         ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ opensuse    ┊ registry.opensuse.org/opensuse/tumbleweed:latest         ┊ yes      ┊ 0    ┊ zypper      ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ ubuntu      ┊ ubuntu:latest                                            ┊ yes      ┊ 0    ┊ apt         ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ vanilla-dev ┊ ghcr.io/vanilla-os/dev:main                              ┊ yes      ┊ 0    ┊ apt         ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ vanilla     ┊ ghcr.io/vanilla-os/pico:main                             ┊ yes      ┊ 0    ┊ apt         ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
```

> **NOTE:** The new `noble` stack is user-defined and therefore not considered a "built-in".

## Updating Stacks

Occasionally, you may need to update an existing stack. This could involve changing the base image to a newer version or modifying the list of installed packages.

To see how to update a stack, you can check the help command:

```bash
apx stacks update --help
```

```
Update the specified stack.

Usage:
  apx stacks update [flags]

Flags:
  -b, --base string          The base subsystem to use.
  -h, --help                 help for update
  -n, --name string          The name of the stack.
  -y, --no-prompt            Assume defaults to all prompts.
  -p, --packages string      The packages to install.
  -k, --pkg-manager string   The package manager to use.
```

In this example, we are going to update the list of installed packages to include `git`.

```bash
apx stacks update noble
```

```
 INFO  Type a new base or confirm the current one (ubuntu:noble):

 INFO  Choose a new package manager or confirm the current one (apt):

 INFO  You have not provided any packages to install in the stack. Do you want to add some now?[y/N]
y
 INFO  Type the packages you want to install in the stack, separated by a space:
neofetch vim git
 INFO  Updated stack 'noble'.
```

Let's check the stack to see if the packages were updated correct.

```bash
apx stacks show noble
```

```
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Name            ┊ noble              ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Base            ┊ ubuntu:noble       ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Packages        ┊ neofetch, vim, git ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Package manager ┊ apt                ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
```

## Exporting Stacks

Exporting a stack allows you to save its configuration and installed packages to a file, which can be shared or stored as a backup.

To see how to export a stack, you can run:

```bash
apx stacks export --help
```

```
Export the specified stack.

Usage:
  apx stacks export [flags]

Flags:
  -h, --help            help for export
  -n, --name string     The name of the stack to export.
  -o, --output string   The path to export the stack to.
```

To export the noble stack to a file named noble-stack.tar.gz, you can use the following command:

```bash
apx stacks export -n noble -o .
cat noble.yml
```

```
name: noble
base: ubuntu:noble
packages:
- git
- neofetch
- vim
pkgmanager: apt
builtin: false
```

## Importing Stacks

Importing a stack allows you to bring a previously exported stack back into your APX environment. This can be useful for sharing stacks between different systems or restoring a stack from a backup.

To see how to import a stack, run the following command:

```bash
apx stacks import --help
```

```
Import the specified stack.

Usage:
  apx stacks import [flags]

Flags:
  -h, --help           help for import
  -i, --input string   The path to import the stack from.
```

Assuming you have a stack file named `noble.yml`, you can import it like this:

```bash
apx stacks import -i noble.yml
```

```
 INFO  Imported stack from 'noble'.
```

```bash
apx stacks show noble
```

```
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Name            ┊ noble              ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Base            ┊ ubuntu:noble       ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Packages        ┊ git, neofetch, vim ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Package manager ┊ apt                ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
```

## Deleting Stacks

Removing a stack allows you to delete a stack that you no longer need. This can help keep your environment organized and free from unused resources.

To see how to remove a stack, run:

```bash
apx stacks rm --help
```

```
Remove the specified stack.

Usage:
  apx stacks rm [flags]

Flags:
  -f, --force         Force removal of the stack.
  -h, --help          help for rm
  -n, --name string   The name of the stack to remove.
```

Stacks can be removed easily with `apx`. We just need to pass the name argument to the `rm` command.

```bash
apx stacks rm -n noble
```

```
 INFO  Are you sure you want to remove 'noble'? [y/N]
y
 INFO  Removed stack 'noble'.
```
