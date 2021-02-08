package winrmhelper

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os/exec"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/masterzen/winrm"
)

// SID is a common structure by all "security principals". This means domains, users, computers, and groups.
// The structure we get from powershell contains more fields, but we're only interested in the Value.
type SID struct {
	Value string `json:"Value"`
}

//WinRMResult holds the stdout, stderr and exit code of a powershell command
type WinRMResult struct {
	Stdout   string
	StdErr   string
	ExitCode int
}

type psString string

func (s *psString) UnmarshalText(text []byte) error {
	str := string(text[:])
	str = strings.TrimSpace(str)
	if str[0] == '+' && len(str) > 2 {
		*s = psString(fmt.Sprintf("\n%s", str[2:]))
	} else {
		*s = psString(str)
	}

	return nil
}

// PSOutput is used to unmarshall CLIXML output
// Right now we are only using this to extract error messages, but it can be extended
// to unpack more elements if required.
type PSOutput struct {
	PSStrings []psString `xml:"S"`
}

func (s *PSOutput) stringSlice() []string {
	out := make([]string, len(s.PSStrings))
	for idx, v := range s.PSStrings {
		out[idx] = string(v)
	}
	return out
}

// String() return a string containing the error message that was serialised in a CLIXML message
func (p *PSOutput) String() string {
	str := strings.Join(p.stringSlice(), "")
	replacer := strings.NewReplacer("_x000D_", "", "_x000A_", "")
	str = replacer.Replace(str)
	return str
}

func decodeXMLCli(xmlDoc string) (string, error) {
	// If stderr is formatted in CLIXML try to extract the error message
	if strings.Contains(xmlDoc, "#< CLIXML") {
		xmlDoc = strings.Replace(xmlDoc, "#< CLIXML", "", -1)
		var v PSOutput
		err := xml.Unmarshal([]byte(xmlDoc), &v)
		if err != nil {
			return "", fmt.Errorf("while unmarshalling CLIXML document: %s", err)
		}
		xmlDoc = strings.TrimSpace(v.String())
	}
	return xmlDoc, nil
}

// PowerShell struct
type PowerShell struct {
	powerShell string
}

// NewPS create new local session
func NewPS() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

const defaultFailedCode = 1

// ExecutePScmd will execute the powershell command using exec
func (p *PowerShell) ExecutePScmd(args ...string) (stdout string, stderr string, exitCode int, err error) {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(p.powerShell, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err = cmd.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			exitCode = defaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	return
}

// RunWinRMCommand will run a powershell command and return the stdout and stderr
// The output is converted to JSON if the json patameter is set to true.
func RunWinRMCommand(conn *winrm.Client, cmds []string, json bool, forceArray bool, execLocally bool) (*WinRMResult, error) {
	if json {
		cmds = append(cmds, "| ConvertTo-Json")
	}

	cmd := strings.Join(cmds, " ")
	encodedCmd := winrm.Powershell(cmd)
	log.Printf("[DEBUG] Running command %s via powershell", cmd)
	log.Printf("[DEBUG] Encoded command: %s", encodedCmd)

	var (
		stdout string
		stderr string
		res    int
		err    error
	)

	if execLocally == false && conn != nil {
		log.Printf("[DEBUG] Executing command on remote host")
		stdout, stderr, res, err = conn.RunWithString(encodedCmd, "")
		log.Printf("[DEBUG] Powershell command exited with code %d", res)
	} else {
		log.Printf("[DEBUG] Creating local shell")
		localShell := NewPS()
		log.Printf("[DEBUG] Executing command on local host")
		stdout, stderr, res, err = localShell.ExecutePScmd(encodedCmd)
	}

	log.Printf("[DEBUG] Powershell command exited with code %d", res)
	if res != 0 {
		log.Printf("[DEBUG] Stdout: %s, Stderr: %s", stdout, stderr)
	}

	// Decode stderr here for the error to be human readable if we need to return early
	stderr, xmlErr := decodeXMLCli(stderr)
	if xmlErr != nil {
		log.Printf("[DEBUG] stderr was not serialised as CLIXML, passing back as is")
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
func SetMachineExtensionNames(client *winrm.Client, gpoDN, value string, execLocally bool) error {
	cmd := fmt.Sprintf(`Set-ADObject -Identity "%s" -Replace @{gPCMachineExtensionNames="%s"}`, gpoDN, value)
	result, err := RunWinRMCommand(client, []string{cmd}, false, false, execLocally)
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
