package winrmhelper

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/masterzen/winrm"
)

//WinRMResult holds the stdout, stderr and exit code of a powershell command
type WinRMResult struct {
	Stdout   string
	StdErr   string
	ExitCode int
}

// RunWinRMCommand will run a powershell command and return the stdout and stderr
// The output is converted to JSON if the json patameter is set to true.
func RunWinRMCommand(conn *winrm.Client, cmds []string, json bool) (*WinRMResult, error) {
	if json {
		cmds = append(cmds, "| convertto-json")
	}
	cmd := strings.Join(cmds, " ")
	encodedCmd := winrm.Powershell(cmd)
	log.Printf("[DEBUG] Running command %s via powershell", cmd)
	log.Printf("[DEBUG] Encoded command: %s", encodedCmd)
	stdout, stderr, res, err := conn.RunWithString(encodedCmd, "")
	log.Printf("[DEBUG] Powershell command exited with code %d", res)
	if res != 0 {
		log.Printf("[DEBUG] Stdout: %s, Stderr: %s", stdout, stderr)
	}
	if err != nil {
		log.Printf("[DEBUG] run error : %s", err)
		return nil, fmt.Errorf("powershell command failed with exit code %d\nstdout: %s\nstderr: %s", res, stdout, stderr)
	}

	result := &WinRMResult{
		Stdout:   stdout,
		StdErr:   stderr,
		ExitCode: res,
	}

	return result, nil
}

// SanitiseTFInput returns the value of a resource field after some basic sanitisation checks
// to protect ourselves from command injection
func SanitiseTFInput(d *schema.ResourceData, key string) string {
	// placeholder for now.
	return d.Get(key).(string)
}
