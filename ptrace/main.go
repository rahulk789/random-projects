package main

import (
	"fmt"
	"os/exec"
	"syscall"
)

func main() {
	// Replace with your actual arguments for `buildah commit`
	buildahArgs := []string{"commit", "nginx-working-container", "my-nginx-image:latest"}

	// Prepare the command
	cmd := exec.Command("buildah", buildahArgs...)

	// Set the command to be traced
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting the command: %v\n", err)
		return
	}

	// Attach to the process
	pid := cmd.Process.Pid
	fmt.Printf("Tracing PID: %d\n", pid)

	// Wait for the process to stop
	var status syscall.WaitStatus
	if _, err := syscall.Wait4(pid, &status, 0, nil); err != nil {
		fmt.Printf("Error waiting for the process: %v\n", err)
		return
	}

	// Continue the process and trace it
	for {
		// Check if the process has exited
		if status.Exited() {
			break
		}

		// Continue the traced process
		if err := syscall.PtraceSyscall(pid, 0); err != nil {
			fmt.Printf("Error continuing process: %v\n", err)
			return
		}

		// Wait for the next stop
		if _, err := syscall.Wait4(pid, &status, 0, nil); err != nil {
			fmt.Printf("Error waiting for process: %v\n", err)
			return
		}
	}

	// Process finished
	fmt.Printf("Process exited with status: %d\n", status.ExitStatus())
}
