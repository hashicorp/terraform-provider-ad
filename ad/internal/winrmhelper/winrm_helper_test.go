// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package winrmhelper

import "testing"

const clixmlError = ` #< CLIXML
	<Objs Version="1.1.0.1" xmlns="http://schemas.microsoft.com/powershell/2004/04"><Obj S="progress" RefId="0">
    <TN RefId="0"><T>System.Management.Automation.PSCustomObject</T><T>System.Object</T></TN><MS>
    <I64 N="SourceId">1</I64><PR N="Record"><AV>Loading Active Directory module for Windows PowerShell with default drive 'AD:'</AV>
	<AI>0</AI><Nil /><PI>-1</PI><PC>0</PC><T>Processing</T><SR>-1</SR><SD> </SD></PR></MS></Obj><Obj S="progress" RefId="1">
	<TNRef RefId="0" /><MS><I64 N="SourceId">1</I64><PR N="Record"><AV>Loading Active Directory module for Windows PowerShell with default drive 'AD:'</AV>
	<AI>0</AI><Nil /><PI>-1</PI><PC>25</PC><T>Processing</T><SR>-1</SR><SD> </SD></PR></MS></Obj><Obj S="progress" RefId="2">
	<TNRef RefId="0" /><MS><I64 N="SourceId">1</I64><PR N="Record"><AV>Loading Active Directory module for Windows PowerShell with default drive 'AD:'</AV>
	<AI>0</AI><Nil /><PI>-1</PI><PC>50</PC><T>Processing</T><SR>-1</SR><SD> </SD></PR></MS></Obj><Obj S="progress" RefId="3">
	<TNRef RefId="0" /><MS><I64 N="SourceId">1</I64><PR N="Record"><AV>Loading Active Directory module for Windows PowerShell with default drive 'AD:'</AV>
	<AI>0</AI><Nil /><PI>-1</PI><PC>75</PC><T>Processing</T><SR>-1</SR><SD> </SD></PR></MS></Obj><Obj S="progress" RefId="4">
	<TNRef RefId="0" /><MS><I64 N="SourceId">1</I64><PR N="Record"><AV>Loading Active Directory module for Windows PowerShell with default drive 'AD:'</AV>
	<AI>0</AI><Nil /><PI>-1</PI><PC>100</PC><T>Processing</T><SR>-1</SR><SD> </SD></PR></MS></Obj><Obj S="progress" RefId="5">
	<TNRef RefId="0" /><MS><I64 N="SourceId">1</I64><PR N="Record"><AV>Loading Active Directory module for Windows PowerShell with default drive 'AD:'</AV>
	<AI>0</AI><Nil /><PI>-1</PI><PC>100</PC><T>Completed</T><SR>-1</SR><SD> </SD></PR></MS></Obj>
	<S S="Error">Set-ADOrganizationalUnit : A parameter cannot be found that matches parameter _x000D__x000A_</S>
	<S S="Error">name 'Path'._x000D__x000A_</S><S S="Error">At line:1 char:101_x000D__x000A_</S>
	<S S="Error">+ ... e description" -Path "DC=yourdomain,DC=com" _x000D__x000A_</S>
	<S S="Error">-ProtectedFromAccidentalDeletion $tr ..._x000D__x000A_</S>
	<S S="Error">+                    ~~~~~_x000D__x000A_</S>
	<S S="Error">    + CategoryInfo          : InvalidArgument: (:) [Set-ADOrganizationalUnit], _x000D__x000A_</S
	><S S="Error">    ParameterBindingException_x000D__x000A_</S>
	<S S="Error">    + FullyQualifiedErrorId : NamedParameterNotFound,Microsoft.ActiveDirectory _x000D__x000A_</S>
	<S S="Error">   .Management.Commands.SetADOrganizationalUnit_x000D__x000A_</S><S S="Error"> _x000D__x000A_</S>
	</Objs>`

func TestDecodeXMLCLI(t *testing.T) {
	expected := `Set-ADOrganizationalUnit : A parameter cannot be found that matches parameter name 'Path'.At line:1 char:101
... e description" -Path "DC=yourdomain,DC=com" -ProtectedFromAccidentalDeletion $tr ...
                   ~~~~~
CategoryInfo          : InvalidArgument: (:) [Set-ADOrganizationalUnit], ParameterBindingException
FullyQualifiedErrorId : NamedParameterNotFound,Microsoft.ActiveDirectory .Management.Commands.SetADOrganizationalUnit`

	msg, err := decodeXMLCli(clixmlError)
	if err != nil {
		t.Fatal(err)
	}
	if msg != expected {
		t.Errorf("actual result did not match the expected one:\nactual: ---%s---\nexpected: ---%s---", msg, expected)
	}
}
