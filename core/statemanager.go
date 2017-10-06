/**
* statemanager.go
*
* @author Sidharth Mishra
* @description The StateManager
* @created Thu Oct 05 2017 00:32:40 GMT-0700 (PDT)
* @copyright
* @last-modified Thu Oct 05 2017 20:11:44 GMT-0700 (PDT)
 */

package core

import "log"

// StateManager is the manager responsible for maintain the object-state mappings,
// uses STM to achieve this.
// Extending the idea of `MemCell` assuming that the `MemCell` is the object that we were
// continuously updating, we break it apart into 2 parts:
// * Immutable part -- Variable
// * Mutable part -- State
// and `associate` both parts with each other using a `stateTable`
// for now, let the `Variable` names be unique.
// * `memory`: is the collection of variablenames - variables (MemCells)
// * `stateTable`: is the mapping between variablename - state (MemCell - it's State)
// * `stm`: The `<i>stm</i>` `associates` <i>Variable</i> to `<i>transaction</i>` that
// owns it during a particular time frame using the `<i>Variable</i>`'s `name`
// as key.
// So, `<b>effectively</b>` the collection of `MemCell`s is now represented by
// the 2 Maps `stateTable` and `memory` (logically) and the `<i>stm</i>` still
// represents the relation between the `<i>MemCell</i>`s and the
// <i>transaction</i>s.
// The <i>stm</i>, <i>memory</i> and <i>stateTable</i> are managed by the
// <i>StateManager</i>.
type statemanager struct {
	memory     map[string]*variable
	stateTable map[string]*State
	stm        map[string]*transaction
}

var manager *statemanager

// Manager : gives the pointer to the StateManager
func Manager() *statemanager {
	if nil == manager {
		log.Printf("Creating new StateManager...\n")
		manager = new(statemanager)
		manager.memory = make(map[string]*variable)
		manager.stateTable = make(map[string]*State)
		manager.stm = make(map[string]*transaction)
	}
	return manager
}

/******** STM OPERATIONS START ***************/

// GetOwner() : Gets the owner of the `MemCell`
func (manager *statemanager) GetOwner(variableName string) *transaction {
	return manager.stm[variableName]
}

// SetOwner(): Sets the owner of the `MemCell`
func (manager *statemanager) SetOwner(variableName string, owner *transaction) {
	manager.stm[variableName] = owner
}

// ReleaseOwnership(): Removes the owner of the transaction reference for the `MemCell`
// releasing it from ownerhsip
func (manager *statemanager) ReleaseOwnership(variableName string) {
	manager.stm[variableName] = nil
}

/******** STM OPERATIONS END ***************/

/****** Object - State, stateTable related START *******/

// MakeVariable(): Makes you a brand new `<i>Variable</i>` or `MemCell` that is allocated in
// * the memory.(JK!)
func (manager *statemanager) MakeVariable(variableName string, props ...*VParam) *variable {
	v := makeVariable(variableName, props...)
	//log.Printf("V: %s", v.ToString())
	//if nil == manager {
	//	log.Printf("lll")
	//}
	manager.memory[variableName] = v
	return v
}

// Read(): Fetches the current state of the variable from the `stateTable`
func (manager *statemanager) Read(variableName string) *State {
	log.Printf("Variable :: name: %s, has State: %s\n", variableName, *(manager.stateTable[variableName]))
	return manager.stateTable[variableName]
}

// Write(): Writes the new state of the <i>Variable</i>, updating its state in the
// `stateTable`.
// <br>
// This action symbolizes that the `MemCell`'s contents were updated.
func (manager *statemanager) Write(variableName string, newValue *State) {
	oldState := manager.stateTable[variableName]
	if nil != oldState {
		log.Printf("Updating Variable :: name: %s with current State: %s\n", variableName, *oldState)
	}
	manager.stateTable[variableName] = newValue
	log.Printf("Updated Variable :: name: %s to new State: %s\n", variableName, *newValue)
}

/****** Object - State, stateTable related END *******/
