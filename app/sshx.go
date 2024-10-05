package app

import (
	"bufio"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	SSHXBinPath = "/opt/sshx/sshx-bin"
)

func (a *App) runSSHX() {
	a.logger.Info().Msg("Starting sshx")
	cmd := exec.Command("ip", "netns", "exec", "srbase-mgmt", SSHXBinPath)

	// Set the process group ID to create a new process group
	// cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		a.logger.Err(err).Msgf("Error creating stdout pipe: %v\n", err)
		return
	}

	urlChan := make(chan string, 1)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			a.logger.Info().Msg(line)

			if strings.Contains(line, "https://sshx.io") {
				url := extractURL(line)
				urlChan <- url

				close(urlChan)

				break
			}
		}
	}()

	// Capture the output
	err = cmd.Start()
	if err != nil {
		a.logger.Err(err).Msgf("Error starting sshx: %v\n", err)
		return
	}

	// Print PID
	a.logger.Info().
		Int("PID", cmd.Process.Pid).
		Msg("sshx started")

	select {
	case url := <-urlChan:
		a.logger.Info().Msgf("Parsed SSHX URL: %s", url)
		a.configState.URL = url
	case <-time.After(30 * time.Second):
		a.logger.Warn().Msg("Timeout waiting for SSHX URL")
	}

	// Let the command continue running
	go func() {
		err := cmd.Wait()
		if err != nil {
			a.logger.Err(err).Msg("sshx process ended with error")
		} else {
			a.logger.Info().Msg("sshx process ended")
		}
	}()
}

// extractURL extracts the sshx URL from the given line.
func extractURL(line string) string {
	// Simple regex to extract URL
	re := regexp.MustCompile(`https://sshx\.io/s/[A-Za-z0-9#]+`)
	return re.FindString(line)
}
