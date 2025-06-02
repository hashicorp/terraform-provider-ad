// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package winrmhelper

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-ad/ad/internal/config"

	"github.com/masterzen/winrm"
)

type CreatePSCommandOpts struct {
	ExecLocally     bool
	ForceArray      bool
	InvokeCommand   bool
	JSONOutput      bool
	PassCredentials bool
	Password        string
	Server          string
	SkipCredPrefix  bool
	SkipCredSuffix  bool
	Username        string
}

type PSCommand struct {
	CreatePSCommandOpts
	cmd string
}

func NewPSCommand(cmds []string, opts CreatePSCommandOpts) *PSCommand {
	if opts.InvokeCommand && opts.PassCredentials {
		invokeCmds := []string{"Invoke-Command -Authentication Kerberos"}
		if opts.JSONOutput {
			cmds = append(cmds, "| ConvertTo-Json")
		}

		invokeCmds = append(invokeCmds, fmt.Sprintf("-ScriptBlock {%s}", strings.Join(cmds, " ")))
		cmds = invokeCmds
	}

	if opts.PassCredentials {
		if !opts.SkipCredPrefix {
			cmdUsername := fmt.Sprintf("$User = \"%s\"\n", opts.Username)
			cmdPassword := fmt.Sprintf("$Password = ConvertTo-SecureString -String \"%s\" -AsPlainText -Force\n", opts.Password)
			cmds = append([]string{"$Credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $User, $Password\n"}, cmds...)
			cmds = append([]string{cmdUsername}, cmds...)
			cmds = append([]string{cmdPassword}, cmds...)
		}
		if !opts.SkipCredSuffix {
			cmds = append(cmds, "-Credential $Credential")
		}
	}

	if opts.PassCredentials && opts.Server != "" {
		switch {
		case opts.InvokeCommand:
			cmds = append(cmds, fmt.Sprintf("-Computername %s", opts.Server))
		default:
			cmds = append(cmds, fmt.Sprintf("-Server %s", opts.Server))
		}
	}

	if !opts.InvokeCommand && opts.JSONOutput {
		cmds = append(cmds, "| ConvertTo-Json")
	}

	cmd := strings.Join(cmds, " ")

	logStr := cmd
	if opts.PassCredentials {
		logStr = strings.ReplaceAll(cmd, opts.Password, "<REDACTED>")
	}
	log.Printf("[DEBUG] Constructing powerrshell command: %s ", logStr)

	res := PSCommand{
		CreatePSCommandOpts: opts,
		cmd:                 cmd,
	}

	return &res
}

// Run will run a powershell command and return the stdout and stderr
// The output is converted to JSON if the json parameter is set to true.
func (p *PSCommand) Run(conf *config.ProviderConf) (*PSCommandResult, error) {
	var (
		stdout string
		stderr string
		res    int
		err    error
	)
	conn, err := conf.AcquireWinRMClient()
	if err != nil {
		return nil, fmt.Errorf("while acquiring winrm client: %s", err)
	}
	defer conf.ReleaseWinRMClient(conn)

	encodedCmd := winrm.Powershell(p.cmd)

	if !p.ExecLocally && conn != nil {
		log.Printf("[DEBUG] Executing command on remote host")
		stdout, stderr, res, err = conn.RunWithString(encodedCmd, "")
		log.Printf("[DEBUG] Powershell command exited with code %d", res)
	} else {
		log.Printf("[DEBUG] Creating local shell")
		localShell := NewLocalPSSession()
		log.Printf("[DEBUG] Executing command on local host")
		stdout, stderr, res, err = localShell.ExecutePScmd(encodedCmd)
	}

	if err != nil {
		log.Printf("[DEBUG] run error : %s", err)
		return nil, fmt.Errorf("powershell command failed with exit code %d\nstdout: %s\nstderr: %s\nerror: %s", res, stdout, stderr, err)
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

	result := &PSCommandResult{
		Stdout:   strings.TrimSpace(stdout),
		StdErr:   stderr,
		ExitCode: res,
	}

	if p.ForceArray && result.Stdout != "" && string(result.Stdout[0]) != "[" {
		result.Stdout = fmt.Sprintf("[%s]", result.Stdout)
	}

	return result, nil
}

func (p *PSCommand) String() string {
	return p.cmd
}

// PSCommandResult holds the stdout, stderr and exit code of a powershell command
type PSCommandResult struct {
	Stdout   string
	StdErr   string
	ExitCode int
}

type psString string

func (s *psString) UnmarshalText(text []byte) error {
	str := string(text)
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

func (p *PSOutput) stringSlice() []string {
	out := make([]string, len(p.PSStrings))
	for idx, v := range p.PSStrings {
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
