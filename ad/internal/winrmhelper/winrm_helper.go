package winrmhelper

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
func RunWinRMCommand(conn *winrm.Client, cmds []string, json bool, forceArray bool) (*WinRMResult, error) {
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

	if json && forceArray && result.Stdout != "" && string(result.Stdout[0]) != "[" {
		result.Stdout = fmt.Sprintf("[%s]", result.Stdout)
	}

	return result, nil
}

// SanitiseTFInput returns the value of a resource field after passing it through SanitiseString
func SanitiseTFInput(d *schema.ResourceData, key string) string {
	return SanitiseString(d.Get(key).(string))

}

// SanitiseString returns the value of a string after some basic sanitisation checks
// to protect ourselves from command injection
func SanitiseString(key string) string {
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
	out := cleanupReplacer.Replace(key)
	log.Printf("[DEBUG] sanitising key %q to: %s", key, out)
	return out
}

// SetMachineExtensionName will add the necessary GUIDs to the GPO's gPCMachineExtensionNames attribute.
// These are required for the security settings part of a GPO to work.
func SetMachineExtensionNames(client *winrm.Client, gpoDN, value string) error {
	cmd := fmt.Sprintf(`Set-ADObject -Identity "%s" -Replace @{gPCMachineExtensionNames="%s"}`, gpoDN, value)
	result, err := RunWinRMCommand(client, []string{cmd}, false, false)
	if err != nil {
		return fmt.Errorf("error while setting machine extension names for GPO %q: %s", gpoDN, err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("command to set machine extension names for GPO %q failed, stderr: %s, stdout: %s", gpoDN, result.StdErr, result.Stdout)
	}
	return nil
}

func GetString(v interface{}) string {
	var out string
	kind := reflect.ValueOf(v).Kind()
	switch kind {
	case reflect.String:
		out = SanitiseString(v.(string))
	case reflect.Float64:
		out = strconv.FormatFloat(v.(float64), 'E', -1, 64)
	case reflect.Int64:
		out = strconv.FormatInt(v.(int64), 10)
	case reflect.Bool:
		out = strconv.FormatBool(v.(bool))
	}
	return fmt.Sprintf(`"%s"`, out)
}

// custom attributes can be single valued or multi valued. Multi-value attribute values are represented by a json
// array that gets converted to a list. It's not guaranteed that the order of the values returned by windows
// will match the order set by the user in the config, so we just check the members of the custom attributes map
// and if a slice is found then it's sorted before we compare it.
func SortInnerSlice(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		if reflect.ValueOf(v).Kind() == reflect.Slice {
			newVal := make([]string, len(v.([]interface{})))
			for idx, attr := range v.([]interface{}) {
				newVal[idx] = GetString(attr)
			}
			sort.Strings(newVal)
			m[k] = newVal
		} else {
			m[k] = GetString(v)
		}
	}
	return m
}
