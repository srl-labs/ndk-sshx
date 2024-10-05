package app

import (
	"bufio"
	"context"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	SSHXBinPath = "/opt/sshx/sshx-bin"
)

var shells = map[string]string{
	"cli":  "sr_cli",
	"bash": "/bin/bash",
}

func (a *App) startSSHX(ctx context.Context) {
	a.logger.Info().Msg("Starting sshx")

	a.KillSSHX(ctx, a.sshxPid)

	cmd := exec.CommandContext(ctx, "ip", "netns", "exec", "srbase-mgmt", SSHXBinPath, "--shell", shells[a.configState.Shell])

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

			// uncomment to see the captured output
			// a.logger.Info().Msg(line)

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

	a.sshxPid = cmd.Process.Pid
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

// KillSSHX kills the sshx process.
func (a *App) KillSSHX(ctx context.Context, pid int) {
	// kill previous sshx process if it exists
	if a.sshxPid != 0 {
		a.logger.Info().
			Int("PID", a.sshxPid).
			Msg("Killing previous sshx process")

		err := exec.CommandContext(ctx, "kill", "-9", strconv.Itoa(a.sshxPid)).Run()
		if err != nil {
			a.logger.Err(err).
				Int("PID", a.sshxPid).
				Msg("Failed to kill previous sshx process")
		}
	}
}
