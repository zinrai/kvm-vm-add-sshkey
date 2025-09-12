# kvm-install-sshkey

`kvm-install-sshkey` is a command-line tool designed to install SSH public keys into KVM (Kernel-based Virtual Machine) virtual machines.

## Features

- Installs SSH public keys into specified KVM virtual machines
- Verifies that the target virtual machine is stopped before proceeding
- Adds the public key to the specified user's `.ssh/authorized_keys` file
- Sets appropriate ownership and permissions for SSH files

## Notes

- This tool requires permissions to run `virsh` and `virt-customize` commands with `sudo`.
- Ensure that the target virtual machine is stopped before running this tool.
- If an `authorized_keys` file already exists, the new public key will be appended to it. Existing keys will be preserved.

## Installation

Build the tool:

```
$ go build
```

## Usage

Basic usage:

```
$ kvm-install-sshkey [options] VM_NAME
```

Example:

```
$ kvm-install-sshkey -user myuser -key /path/to/my/key.pub my-vm
```

## License

This project is licensed under the [MIT License](./LICENSE).
