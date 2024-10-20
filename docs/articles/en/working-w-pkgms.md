---
Title: Working with Package Managers
Description: Learn how to work with package managers in Apx, including creating, editing, deleting, importing, and exporting configurations.
PublicationDate: 2024-10-18
Listed: true
Authors:
  - jardon
Tags:
  - working
  - package-managers
---

Package managers can be manipulated in various ways with `apx`. This includes creation, editing, deleting, importing configurations, and exporting configurations.

```bash
apx pkgmanagers --help
```

```
Work with the package managers that are available in apx.

Usage:
  apx pkgmanagers [command]

Available Commands:
  export      Export the specified package manager.
  import      Import the specified package manager.
  list        List all available package managers.
  new         Create a new package manager.
  rm          Remove the specified package manager.
  show        Show information about the specified package manager.
  update      Update the specified package manager.

Flags:
  -h, --help   help for pkgmanagers

Use "apx pkgmanagers [command] --help" for more information about a command.
```

## Creating a Package Manager

APX ships with several package managers out-of-the-box. These cover a wide range of different Linux distributions. In the event that you need to add a package manager for a new operating system, this can be done using the `apx` CLI. Here are the available options for adding a package manager:

```bash
apx pkgmanagers new --help
```

````
Create a new package manager.

Usage:
  apx pkgmanagers new [flags]

Flags:
  -a, --autoremove string   The command to run to autoremove packages.
  -c, --clean string        The command to run to clean the package manager's cache.
  -h, --help                help for new
  -i, --install string      The command to run to install packages.
  -l, --list string         The command to run to list installed packages.
  -n, --name string         The name of the package manager.
  -S, --need-sudo           Whether the package manager needs sudo to run.
  -y, --no-prompt           Assume defaults to all prompts.
  -p, --purge string        The command to run to purge packages.
  -r, --remove string       The command to run to remove packages.
  -s, --search string       The command to run to search for packages.
  -w, --show string         The command to run to show information about packages.
  -u, --update string       The command to run to update the list of available packages.
  -U, --upgrade string      The command to run to upgrade packages.
```bash
apx pkgmanagers list
````

```
 INFO  Found 5 package managers
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
┊ NAME   ┊ BUILT-IN ┊
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
┊ apk    ┊ yes      ┊
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
┊ apt    ┊ yes      ┊
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
┊ dnf    ┊ yes      ┊
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
┊ pacman ┊ yes      ┊
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
┊ zypper ┊ yes      ┊
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
```

Looking at the list, we see that `yum` is not preconfigured for us. This means that we need to add a package manager to `apx` and give it all the required parameters.

```bash
apx pkgmanagers new
```

> **NOTE:** These options can be passed as CLI args as shown above in the help output.

```
 INFO  Choose a name:
yum
 INFO  Does the package manager need sudo to run? [y/N]
y
 INFO  Enter the command for 'list':
yum list
 INFO  Enter the command for 'show':
yum info
 INFO  Enter the command for 'autoRemove':
yum autoremove
 INFO  Enter the command for 'clean':
yum clean
 INFO  Enter the command for 'remove':
yum remove
 INFO  Enter the command for 'search':
yum search
 INFO  Enter the command for 'update':
yum upgrade --refresh
 INFO  Enter the command for 'upgrade':
yum upgrade
 INFO  Enter the command for 'install':
yum install
 INFO  Enter the command for 'purge':
yum remove
 SUCCESS  Package manager yum created successfully.
```

Now we should be able to see the new yum package manager in the list of available managers.

```bash
apx pkgmanagers list
```

```
 INFO  Found 6 package managers
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
┊ NAME   ┊ BUILT-IN ┊
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
┊ yum    ┊ no       ┊
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
┊ apk    ┊ yes      ┊
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
┊ apt    ┊ yes      ┊
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
┊ dnf    ┊ yes      ┊
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
┊ pacman ┊ yes      ┊
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
┊ zypper ┊ yes      ┊
┼┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┼
```

> **NOTE:** The new `yum` package manager is user-defined and therefore not considered a "built-in".

## Updating a Package Manager

Updates to package managers can be done similarly to other operations in `apx`.

```bash
apx pkgmanagers update --help
```

```
Update the specified package manager.

Usage:
  apx pkgmanagers update [flags]

Flags:
  -a, --autoremove string   The command to run to autoremove packages.
  -c, --clean string        The command to run to clean the package manager's cache.
  -h, --help                help for update
  -i, --install string      The command to run to install packages.
  -l, --list string         The command to run to list installed packages.
  -n, --name string         The name of the package manager.
  -S, --need-sudo           Whether the package manager needs sudo to run.
  -y, --no-prompt           Assume defaults to all prompts.
  -p, --purge string        The command to run to purge packages.
  -r, --remove string       The command to run to remove packages.
  -s, --search string       The command to run to search for packages.
  -w, --show string         The command to run to show information about packages.
  -u, --update string       The command to run to update the list of available packages.
  -U, --upgrade string      The command to run to upgrade packages.
```

Now let's update the autoremove command to add verbosity. This requires adding the `-v` argument to the command.

```bash
apx pkgmanagers update -n yum -S
```

> **NOTE:** Currently, sudo access with `-S` needs to be used every time an update occurs to maintain it. Leaving off the `-S` will remove sudo access from the command.

```
 INFO  Enter new command for 'autoRemove' (leave empty to keep 'yum autoremove'):
yum autoremove -v
 INFO  Enter new command for 'clean' (leave empty to keep 'yum clean'):

 INFO  Enter new command for 'install' (leave empty to keep 'yum'):

 INFO  Enter new command for 'list' (leave empty to keep 'nstall'):

 INFO  Enter new command for 'purge' (leave empty to keep 'v'):

 INFO  Enter new command for 'remove' (leave empty to keep 'yum remove'):

 INFO  Enter new command for 'search' (leave empty to keep 'yum search'):

 INFO  Enter new command for 'show' (leave empty to keep 'yum info'):

 INFO  Enter new command for 'update' (leave empty to keep 'yum upgrade --refresh'):

 INFO  Enter new command for 'upgrade' (leave empty to keep 'yum upgrade'):

 INFO  Updated package manager 'yum'.
```

```bash
apx pkgmanagers show yum
```

```
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Name       ┊ yum                   ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ NeedSudo   ┊ true                  ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ AutoRemove ┊ yum autoremove -v     ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Clean      ┊ yum clean             ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Install    ┊ yum install           ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ List       ┊ yum list              ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Purge      ┊ yum remove            ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Remove     ┊ yum remove            ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Search     ┊ yum search            ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Show       ┊ yum info              ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Update     ┊ yum upgrade --refresh ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Upgrade    ┊ yum upgrade           ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┼
```

## Exporting a Package Manager

`apx` can export a package manager to a yaml file to be imported later.

```bash
apx pkgmanagers export --help
```

```
Export the specified package manager.

Usage:
  apx pkgmanagers export [flags]

Flags:
  -h, --help            help for export
  -n, --name string     The name of the package manager to export.
  -o, --output string   The path to export the stack to.
```

Now let's say we are wanting to back up our package manager configuration for usage with some automation tooling. We just need to specify the name of the package manager and a location to export the yaml file.

```bash
apx pkgmanagers export -n yum -o .
cat yum.yml
```

```
model: 2
name: yum
needsudo: true
cmdautoremove: yum remove
cmdclean: yum clean
cmdinstall: yum install
cmdlist: yum list
cmdpurge: yum remove
cmdremove: yum remove
cmdsearch: yum search
cmdshow: yum show
cmdupdate: yum upgrade
cmdupgrade: yum upgrade
builtin: false
```

The package manager has been successfully exported!

## Importing a Package Manager

Package managers can be imported easily with `apx`.

```bash
apx pkgmanagers import --help
```

```
Import the specified package manager.

Usage:
  apx pkgmanagers import [flags]

Flags:
  -h, --help           help for import
  -i, --input string   The path to import the package manager from.
```

For this example, we are going to import the yaml file that we exported in the previous section.

```bash
apx pkgmanagers import -i yum.yml
```

```
 INFO  Imported package manager from 'yum'.
```

```bash
apx pkgmanagers show yum
```

```
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Name       ┊ yum         ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ NeedSudo   ┊ true        ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ AutoRemove ┊ yum remove  ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Clean      ┊ yum clean   ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Install    ┊ yum install ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ List       ┊ yum list    ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Purge      ┊ yum remove  ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Remove     ┊ yum remove  ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Search     ┊ yum search  ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Show       ┊ yum show    ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Update     ┊ yum upgrade ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
┊ Upgrade    ┊ yum upgrade ┊
┼┄┄┄┄┄┄┄┄┄┄┄┄┼┄┄┄┄┄┄┄┄┄┄┄┄┄┼
```

The package manager has been successfully imported!

## Deleting a Package Manager

```bash
apx pkgmanagers rm --help
```

```
Remove the specified package manager.

Usage:
  apx pkgmanagers rm [flags]

Flags:
  -f, --force         Force removal of the package manager.
  -h, --help          help for rm
  -n, --name string   The name of the package manager to remove.
```

Since we are not using the `yum` package manager, we want to delete it. This can be done by passing the name of the package manager to the `rm` command.

```bash
apx pkgmanagers rm -n yum
```

```
 INFO  Are you sure you want to remove 'yum'? [y/N]
y
 INFO  Removed package manager 'yum'.
```

And with that we no longer have the yum package manager!
