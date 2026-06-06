package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	username := flag.String("user", "", "Username on the virtual machine (optional)")
	keyPath := flag.String("key", "", "Path to the SSH public key file (default: $HOME/.ssh/id_rsa.pub)")

	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Println("Usage: kvm-vm-add-sshkey [options] VM_NAME")
		flag.PrintDefaults()
		os.Exit(1)
	}

	vmName := flag.Arg(0)

	if !isVMStopped(vmName) {
		fmt.Printf("The virtual machine %s is not stopped. Please stop it before proceeding.\n", vmName)
		os.Exit(1)
	}

	if *keyPath == "" {
		currentUser := os.Getenv("USER")
		userHome, err := getUserHome(currentUser)
		if err != nil {
			fmt.Printf("Failed to get home directory for user %s: %v\n", currentUser, err)
			os.Exit(1)
		}

		*keyPath = filepath.Join(userHome, ".ssh", "id_rsa.pub")
	}

	pubKey, err := os.ReadFile(*keyPath)
	if err != nil {
		fmt.Printf("Failed to read public key: %v\n", err)
		os.Exit(1)
	}

	if *username == "" {
		*username = os.Getenv("USER")
		if *username == "" {
			fmt.Println("Failed to get current user. Please specify a username using the --user flag.")
			os.Exit(1)
		}
	}

	err = installSSHKey(vmName, *username, string(pubKey))
	if err != nil {
		fmt.Printf("Failed to install SSH key: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully installed SSH public key for user %s on the VM %s.\n", *username, vmName)
}

func installSSHKey(vmName, username, pubKey string) error {
	script := fmt.Sprintf(`
mkdir -p /home/%[1]s/.ssh
echo '%[2]s' >> /home/%[1]s/.ssh/authorized_keys
chown -R %[1]s:%[1]s /home/%[1]s/.ssh
chmod 700 /home/%[1]s/.ssh
chmod 600 /home/%[1]s/.ssh/authorized_keys
`, username, strings.TrimSpace(pubKey))

	cmd := exec.Command("sudo", "virt-customize", "-d", vmName, "--run-command", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("virt-customize failed: %v, output: %s", err, output)
	}

	return nil
}

func isVMStopped(vmName string) bool {
	cmd := exec.Command("sudo", "virsh", "list", "--name", "--state-running")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to get list of running VMs: %v\n", err)
		return false
	}

	runningVMs := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, vm := range runningVMs {
		if vm == vmName {
			return false
		}
	}
	return true
}

func getUserHome(username string) (string, error) {
	cmd := exec.Command("getent", "passwd", username)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %v", err)
	}

	fields := strings.Split(string(output), ":")
	if len(fields) < 6 {
		return "", fmt.Errorf("unexpected output format from getent")
	}

	return fields[5], nil
}
