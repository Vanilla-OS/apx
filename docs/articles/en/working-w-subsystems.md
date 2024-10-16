## Working with Subsystems
You can interact with subsystems using the `apx` CLI as shown below.  Managing subsystems is as easy as simple command in a shell.

To see the available commands for managing subsystems, use:
```bash
apx subsystems --help
```
```
Work with the subsystems that are available in apx.

Usage:
  apx subsystems [command]

Available Commands:
  list        List all available subsystems.
  new         Create a new subsystem.
  reset       Reset the specified subsystem.
  rm          Remove the specified subsystem.

Flags:
  -h, --help   help for subsystems

Use "apx subsystems [command] --help" for more information about a command.
```

## Creating a Subsystem
Subsystems are individual operating system containers built on the concept of stacks in APX. These subsystems can be created and deleted independently of stacks. You can have multiple isolated subsystems built on the same stack!

To see how to create a new subsystem, run:
```bash
apx subsystems new --help
```
```
Create a new subsystem.

Usage:
  apx subsystems new [flags]

Flags:
  -h, --help           help for new
  -H, --home string    The custom home directory of the subsystem.
  -i, --init           Use systemd inside the subsystem.
  -n, --name string    The name of the subsystem.
  -s, --stack string   The stack to use.
```

Now that we have created a stack for Ubuntu 24.04, we want to create a new subsystem that is built from that stack.  To do this we need to supply a couple parameters to `apx`.
```bash
apx subsystems new
```
You will be prompted to choose a name and select a stack:
```
 INFO  Choose a name:
noble-test
 INFO  Available stacks:
1. noble
2. alpine
3. arch
4. fedora
5. opensuse
6. ubuntu
7. vanilla-dev
8. vanilla
 INFO  Select a stack [1-8]:
1
▀  Creating subsystem 'noble-test' with stack 'noble'… (0s)Resolved "ubuntu" as an alias (/etc/containers/registries.conf.d/shortnames.conf)
Trying to pull docker.io/library/ubuntu:noble...
 ▀ Creating subsystem 'noble-test' with stack 'noble'… (1s)Getting image source signatures
 ▄ Creating subsystem 'noble-test' with stack 'noble'… (1s)Copying blob ff65ddf9395b [--------------------------------------] 0.0b / 28.4MiB (skipped: 0.0b = 0.00%)                           Copying blob ff65ddf9395b [--------------------------------------] 0.0b / 28.4MiB (skipped: 0.0b = 0.00%)                           Copying blob ff65ddf9395b [=>------------------------------------] 1.4MiBCopying blob ff65ddf9395b [=>------------------------------Copying blob ff65ddf9395b [=>------------------------------------] 1.7MiB / 28.4MiB | 2.2 MiB/s                                     Copying blob ff65ddf9395b [===>----------------------------------] 3.2MiB / 28.4MiB | 16.5 MiB/s                                    Copying blob ff65ddf9395b [======>-------------------------------] 4.9MiBCopying blob ff65ddf9395b [========>-----------------------Copying blob ff65ddf9395b [===========>--------------------------] 8.8MiB / 28.4MiB | 19.9 MiB/s                                    Copying blob ff65ddf9395b [==============>-----------------------] 10.9MiB / 28.4MiB | 17.6 MiB/s                                   Copying blob ff65ddf9395b [===============>----------------------] 11.9MiCopying blob ff65ddf9395b [=================>--------------Copying blob ff65ddf9395b [====================>-----------------] 16.0MiB / 28.4MiB | 31.8 MiB/s                                   Copying blob ff65ddf9395b [========================>-------------] 18.7MiB / 28.4MiB | 23.7 MiB/s                                   Copying blob ff65ddf9395b [============================>---------] 21.5MiCopying blob ff65ddf9395b [================================Copying blob ff65ddf9395b [====================================>-] 27.5MiCopying blob ff65ddf9395b done   |                         Copying blob ff65ddf9395b done   | 
Copying config 59ab366372 done   | 
Writing manifest to image destination
59ab366372d56772eb54e426183435e6b0642152cb449ec7ab52473af8ca6e3f
 ▄ Creating subsystem 'noble-test' with stack 'noble'… (4s) [ OK ]
Distrobox 'apx-noble-test' successfully created.
To enter, run:

apx noble-test enter
 SUCCESS  Created subsystem 'noble-test'.
```
In the above example, we created a new subsystem named `noble-test` from the `noble` stack.  We can confirm this by listing the available subsystems in `apx`.
```bash
apx subsystems list
```
```
 INFO  Found 1 subsystems:
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼
┊ NAME       ┊ STACK ┊ STATUS       ┊ PKGS ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼
┊ noble-test ┊ noble ┊ Up 7 minutes ┊ 2    ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┼
```
Now that we see that `noble-test` exists, let's try to use the subsystem.  We want to `enter` the subsystem which should give us a shell inside the container.
```bash
apx noble-test enter
```
```
Starting container...                   	 [ OK ]
Installing basic packages...            	 [ OK ]
Setting up devpts mounts...             	 [ OK ]
Setting up read-only mounts...          	 [ OK ]
Setting up read-write mounts...         	 [ OK ]
Setting up host's sockets integration...	 [ OK ]
Setting up host's nvidia integration... 	 [ OK ]
Integrating host's themes, icons, fonts...	 [ OK ]
Setting up package manager exceptions...	 [ OK ]
Setting up package manager hooks...     	 [ OK ]
Setting up dpkg exceptions...           	 [ OK ]
Setting up apt hooks...                 	 [ OK ]
Setting up distrobox profile...         	 [ OK ]
Setting up sudo...                      	 [ OK ]
Setting up user groups...               	 [ OK ]
Setting up kerberos integration...      	 [ OK ]
Setting up user's group list...         	 [ OK ]
Setting up existing user...             	 [ OK ]
Setting up user home...                 	 [ OK ]
Ensuring user's access...               	 [ OK ]

Container Setup Complete!
jardon@apx-noble-test:~$
```
Now that we are inside the container, we can confirm the stack and that it installed `neofetch` like we asked it to when creating the stack.
```bash
neofetch --off --color_blocks off --stdout
```
```
jardon@lagann 
------------- 
OS: Ubuntu 24.04.1 LTS x86_64 
Host: Precision 5560 
Kernel: 6.9.8-amd64 
Uptime: 1 day, 20 hours, 14 mins 
Packages: 438 (dpkg) 
Shell: bash 5.2.21 
Resolution: 3840x2400 
DE: GNOME 
WM: Mutter 
Theme: Yaru [GTK3] 
Icons: Yaru [GTK3] 
Terminal: conmon 
CPU: 11th Gen Intel i7-11850H (16) @ 4.800GHz 
GPU: Intel TigerLake-H GT1 [UHD Graphics] 
GPU: NVIDIA T1200 Laptop GPU 
Memory: 12599MiB / 31827MiB 
```
With that, we have successfully created and used a subsystem built on a previously user-defined stack!

## Deleting a Subsystem
Removing a subsystem with `apx` is easy.  Just pass the name of the subsystem to the `apx` command and confirm the deletion.
```bash
apx subsystem rm --help
```
```
Remove the specified subsystem.

Usage:
  apx subsystems rm [flags]

Flags:
  -f, --force         Force removal of the subsystem.
  -h, --help          help for rm
  -n, --name string   The name of the subsystem to remove.
```

```bash
apx subsystem rm -n noble-test
```
```
 INFO  Are you sure you want to remove 'noble-test'? [y/N]
y
 SUCCESS  Removed subsystem 'noble-test'.
```

And with that you have successfully deleted the subsystem!
