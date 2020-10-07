package winrmhelper

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mainzen/winrm"
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
		return nil, fmt.Errorf("powershell command failed with exit code %d\nstdout: %s\nstderr: %s\nerror: %s", res, stdout, stderr, err)
	}

	result := &WinRMResult{
		Stdout:   strings.TrimSpace(stdout),
		StdErr:   stderr,
		ExitCode: res,
	}

	return result, nil
}

// SanitiseTFInput returns the value of a resource field after some basic sanitisation checks
// to protect ourselves from command injection
func SanitiseTFInput(d *schema.ResourceData, key string) string {
	cleanupReplacer := strings.NewReplacer(
		"`", "``",
		`"`, "`\"",
		"$", "`$",
		"\x00", "`0",
		"\x07", "`a",
		"\x08", "`b",
		"\x1f", "`e",
		"\x0c", "`f",
		"\n", "`n",
		"\r", "`r",
		"\t", "`t",
		"\v", "`v",
	)

	out := cleanupReplacer.Replace(d.Get(key).(string))
	log.Printf("[DEBUG] sanitising key %q to: %s", key, out)
	return out
}

// SetMachineExtensionName will add the necessary GUIDs to the GPO's gPCMachineExtensionNames attribute.
// These are required for the security settings part of a GPO to work.
func SetMachineExtensionNames(client *winrm.Client, gpoDN, value string) error {
	cmd := fmt.Sprintf(`Set-ADObject -Identity "%s" -Replace @{gPCMachineExtensionNames="%s"}`, gpoDN, value)
	result, err := RunWinRMCommand(client, []string{cmd}, false)
	if err != nil {
		return fmt.Errorf("error while setting machine extension names for GPO %q: %s", gpoDN, err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("command to set machine extension names for GPO %q failed, stderr: %s, stdout: %s", gpoDN, result.StdErr, result.Stdout)
	}
	return nil
}
