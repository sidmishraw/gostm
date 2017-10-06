/**
* memcell.go
*
* @author Sidharth Mishra
* @description contains the code related to MemCell representations
* @created Thu Oct 05 2017 18:18:19 GMT-0700 (PDT)
* @copyright Sidharth Mishra 2017
* @last-modified Thu Oct 05 2017 18:24:33 GMT-0700 (PDT)
 */

package core

import (
	"fmt"
)

// The variable represents the row identifier in the State table.
// This is also the immutable part of the MemCell.
// Separate `Object` from its `State` and `associate` them with each other with
// the help of a table like data structure called a `StateTable`.
// Object := 	Immutable part called `Variable` since `Object` is reserved in Java
//				|   Mutable part called `State`
// It is supposed to contain all the values that represent the <em> Object </em>
// uniquely in the world. It is associated with its <em> State </em> which is
// maintained in the StateTable.
type variable struct {
	name                string
	immutableProperties map[string]interface{}
}

// VParam is an entry in the Variable's immutable param list
// Each entry has a name and a value
type VParam struct {
	Name  string
	Value interface{}
}

// Variable is the new Variable
// @param name The name of the `Variable`
// @returns The variable
func makeVariable(name string, params ...*VParam) *variable {
	v := new(variable)
	v.name = name
	v.immutableProperties = make(map[string]interface{})
	// add in the params
	for _, p := range params {
		// log.Printf("pName: %s, p: %s", p.Name, p.Value)
		v.immutableProperties[p.Name] = p.Value
	}
	
	return v
}

// GetValue(): gets you the value of the property of the `Variable` if it exists
// otherwise you get nothing(nil)
func (v *variable) GetValue(name string) interface{} {
	return v.immutableProperties[name]
}

// ToString() gives the string representation of Variable
func (v *variable) ToString() string {
	return fmt.Sprintf(`
		{
			Variable: {
				"name": "%s",
				"immutableProperties": "%s"
			}
		}
`, v.name, v.immutableProperties)
}

// State represents the <em> State </em> represents the State of the `Object`.
// It is maintained in the `<i>StateTable</i>`. It is the `mutable` part of the `MemCell`.
type State interface {
}
