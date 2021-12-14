package winrmhelper

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"
	"github.com/packer-community/winrmcp/winrmcp"
)

// SID is a common structure by all "security principals". This means domains, users, computers, and groups.
// The structure we get from powershell contains more fields, but we're only interested in the Value.
type SID struct {
	Value string `json:"Value"`
}

// LocalPSSession struct
type LocalPSSession struct {
	powerShell string
}

// NewLocalPSSession create new local session
func NewLocalPSSession() *LocalPSSession {
	ps, _ := exec.LookPath("powershell.exe")
	return &LocalPSSession{
		powerShell: ps,
	}
}

const defaultFailedCode = 1

// ExecutePScmd will execute the powershell command using exec
func (l *LocalPSSession) ExecutePScmd(args ...string) (stdout string, stderr string, exitCode int, err error) {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(l.powerShell, args...)
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

// SetMachineExtensionNames will add the necessary GUIDs to the GPO's gPCMachineExtensionNames attribute.
// These are required for the security settings part of a GPO to work.
func SetMachineExtensionNames(conf *config.ProviderConf, gpoDN, value string) error {
	cmd := fmt.Sprintf(`Set-ADObject -Identity "%s" -Replace @{gPCMachineExtensionNames="%s"}`, gpoDN, value)
	psOpts := CreatePSCommandOpts{
		JSONOutput:      false,
		ForceArray:      false,
		ExecLocally:     conf.IsConnectionTypeLocal(),
		PassCredentials: conf.IsPassCredentialsEnabled(),
		Username:        conf.Settings.WinRMUsername,
		Password:        conf.Settings.WinRMPassword,
		Server:          conf.Settings.DomainName,
	}
	psCmd := NewPSCommand([]string{cmd}, psOpts)
	result, err := psCmd.Run(conf)
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

// SortInnerSlice is used to sort multivalued custom attributes.
// Custom attributes can be single valued or multi valued. Multi-value attribute values are represented by a json
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

func UploadFiletoSYSVOL(conf *config.ProviderConf, cpClient *winrmcp.Winrmcp, buf io.Reader, destPath string) error {
	tmpPathCmd := NewPSCommand([]string{"$randompath=[System.IO.Path]::GetRandomFileName(); echo $env:TMP\\$randompath"}, CreatePSCommandOpts{
		ForceArray:      false,
		JSONOutput:      false,
		ExecLocally:     conf.IsConnectionTypeLocal(),
		PassCredentials: false,
		SkipCredPrefix:  true,
		SkipCredSuffix:  true,
	})
	tmpPathResult, err := tmpPathCmd.Run(conf)
	if err != nil {
		return fmt.Errorf("while renaming GPO: %s", err)
	} else if tmpPathResult != nil && tmpPathResult.ExitCode != 0 {
		return fmt.Errorf("while renaming GPO stderr: %s", tmpPathResult.StdErr)
	}
	tmpPath := tmpPathResult.Stdout

	err = cpClient.Write(tmpPath, buf)
	if err != nil {
		return fmt.Errorf("error while writing ini file to %q: %s", destPath, err)
	}

	toks := strings.Split(destPath, `\`)
	x := toks[:len(toks)-1]
	destDir := strings.Join(x, `\`)
	mdCmd := fmt.Sprintf(`$check=Test-Path "%s"; if (!$check)  {md "%s"}`, destDir, destDir)
	domainName := conf.Settings.DomainName
	if conf.Settings.KrbRealm == domainName {
		domainName = "$env:computername"
	}
	mdPSComamnd := NewPSCommand([]string{mdCmd}, CreatePSCommandOpts{
		ExecLocally:     conf.IsConnectionTypeLocal(),
		JSONOutput:      false,
		ForceArray:      false,
		PassCredentials: conf.IsPassCredentialsEnabled(),
		InvokeCommand:   conf.IsPassCredentialsEnabled(),
		Username:        conf.Settings.WinRMUsername,
		Password:        conf.Settings.WinRMPassword,
		Server:          domainName,
	})
	mdOutput, err := mdPSComamnd.Run(conf)
	if err != nil {
		return fmt.Errorf("while renaming GPO: %s", err)
	} else if mdOutput != nil && mdOutput.ExitCode != 0 {
		return fmt.Errorf("while renaming GPO stderr: %s", mdOutput.StdErr)
	}

	cpCmd := fmt.Sprintf(`Copy-Item "%s" "%s"; Remove-Item "%s"`, tmpPath, destPath, tmpPath)
	cpPSComamnd := NewPSCommand([]string{cpCmd}, CreatePSCommandOpts{
		ExecLocally:     conf.IsConnectionTypeLocal(),
		JSONOutput:      false,
		ForceArray:      false,
		PassCredentials: conf.IsPassCredentialsEnabled(),
		InvokeCommand:   conf.IsPassCredentialsEnabled(),
		Username:        conf.Settings.WinRMUsername,
		Password:        conf.Settings.WinRMPassword,
		Server:          domainName,
	})
	cpOutput, err := cpPSComamnd.Run(conf)
	if err != nil {
		return fmt.Errorf("while renaming GPO: %s", err)
	} else if cpOutput != nil && cpOutput.ExitCode != 0 {
		return fmt.Errorf("while renaming GPO stderr: %s", cpOutput.StdErr)
	}

	return nil
}

func GetOtherAttributes(customAttributes map[string]interface{}) (string, error) {
	out := []string{}
	for k, v := range customAttributes {
		cleanKey := SanitiseString(k)
		var cleanValue string
		if reflect.ValueOf(v).Kind() == reflect.Slice {
			quotedStrings := make([]string, len(v.([]interface{})))
			for idx, s := range v.([]interface{}) {
				// Using %q here will cause double quotes inside the string to be escaped with \"
				// which is not desirable in Powershell
				quotedStrings[idx] = GetString(s.(string))
			}
			cleanValue = strings.Join(quotedStrings, ",")
		} else {
			cleanValue = GetString(v.(string))
		}
		out = append(out, fmt.Sprintf(`'%s'=%s`, cleanKey, cleanValue))
	}
	finalAttrString := strings.Join(out, ";")
	return fmt.Sprintf("@{%s}", finalAttrString), nil
}

func GetChangesForCustomAttributes(oldValue, newValue interface{}) ([]string, error) {
	cmds := []string{}
	newMap, err := structure.ExpandJsonFromString(newValue.(string))
	if err != nil {
		return nil, err
	}
	newSortedMap := SortInnerSlice(newMap)
	toClear := []string{}
	toReplace := []string{}
	toAdd := []string{}

	var oldSortedMap map[string]interface{}
	if oldValue.(string) != "" {
		oldMap, err := structure.ExpandJsonFromString(oldValue.(string))
		if err != nil {
			return nil, fmt.Errorf("while expanding CA json string %s: %s", oldValue.(string), err)
		}
		oldSortedMap = SortInnerSlice(oldMap)
	}

	for k, v := range oldSortedMap {
		if newVal, ok := newSortedMap[k]; ok {
			if !reflect.DeepEqual(v, newVal) {
				var out string
				if reflect.ValueOf(newVal).Kind() == reflect.Slice {
					quotedStrings := make([]string, len(newVal.([]string)))
					for idx, s := range newVal.([]string) {
						quotedStrings[idx] = s
					}
					out = strings.Join(quotedStrings, ",")
				} else {
					out = newVal.(string)
				}
				toReplace = append(toReplace, fmt.Sprintf("%s=%s", SanitiseString(k), out))
			}
		} else {
			toClear = append(toClear, SanitiseString(k))
		}
	}

	for k, newVal := range newSortedMap {
		if _, ok := oldSortedMap[k]; !ok {
			var out string
			if reflect.ValueOf(newVal).Kind() == reflect.Slice {
				quotedStrings := make([]string, len(newVal.([]string)))
				for idx, s := range newVal.([]string) {
					// Using %q here will cause double quotes inside the string to be escaped with \"
					// which is not desirable in Powershell
					quotedStrings[idx] = s
				}
				out = strings.Join(quotedStrings, ",")
			} else {
				out = newVal.(string)
			}
			toAdd = append(toAdd, fmt.Sprintf("%s=%s", SanitiseString(k), out))
		}
	}

	if len(toClear) > 0 {
		cmds = append(cmds, fmt.Sprintf(`-Clear %s`, strings.Join(toClear, ",")))
	}

	if len(toReplace) > 0 {
		cmds = append(cmds, fmt.Sprintf(`-Replace @{%s}`, strings.Join(toReplace, ";")))
	}

	if len(toAdd) > 0 {
		cmds = append(cmds, fmt.Sprintf(`-Add @{%s}`, strings.Join(toAdd, ";")))
	}
	return cmds, nil
}
