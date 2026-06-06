# kvm-vm-add-sshkey

`kvm-vm-add-sshkey` is a command-line tool that adds an SSH public key to a KVM (Kernel-based Virtual Machine) virtual machine by writing directly to the guest disk while the VM is stopped.

## Features

- Appends an SSH public key to the specified user's `.ssh/authorized_keys` file inside the guest
- Verifies that the target virtual machine is stopped before proceeding
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
$ kvm-vm-add-sshkey [options] VM_NAME
```

Example:

```
$ kvm-vm-add-sshkey -user myuser -key /path/to/my/key.pub my-vm
```

## License

This project is licensed under the [MIT License](./LICENSE).
