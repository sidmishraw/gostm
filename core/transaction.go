/**
* transaction.go
*
* @author Sidharth Mishra
* @description transaction of the STM
* @created Thu Oct 05 2017 00:29:41 GMT-0700 (PDT)
* @copyright
* @last-modified Thu Oct 05 2017 18:56:19 GMT-0700 (PDT)
 */

package core

import (
	"log"
	"sync"
	"time"
)

// Max sleep wait time in Milliseconds
const MAX_SLEEP_WAIT_TIME int = 1000

// transaction is a transaction
type transaction struct {
	record    *record
	manager   *statemanager
	operation func() bool
}

// Record is the metadata of a transaction
// It holds the readSet, writeSet and oldValues
type record struct {
	status    bool
	version   int
	name      string
	writeSet  []string
	readSet   []string
	oldValues map[string]*State
}

// Run : initiates the transaction
func (t *transaction) Run(waitGroup *sync.WaitGroup) {
	
	log.Printf("transaction:: %s has started...\n", t.record.name)
	
	for !t.record.status {
		
		log.Printf("Initiating transaction:: %s\n", t.record.name)
		log.Printf("Taking ownership of `writeSet` members of the transaction:: %s\n", t.record.name)
		
		// take ownership of `writeSet` members
		ownershipStatus := t.takeOwnership()
		
		if !ownershipStatus {
			
			// failed to take ownership, rollback and goto sleep and then retry
			// from the beginning
			log.Printf("MODERATE:: transaction %s has failed to take ownership of all of its writeSet members, retrying after sometime", t.record.name)
			
			// since there has been no modification to the writeSet members
			// yet, the only change that needs to be reverted is
			// the release of writeSet members' ownership.
			t.releaseOwnership()
			
			log.Printf("transaction:: %s going to wait for %d before retrying...\n", t.record.name, MAX_SLEEP_WAIT_TIME)
			
			// sleeps for MAX_SLEEP_WAIT_TIME ms
			time.Sleep(time.Duration(MAX_SLEEP_WAIT_TIME))
			
			continue
		}
		
		// transaction takes ownership of write set,now onto backing up
		log.Printf("transaction:: %s has taken ownership successfully, now moving on to taking backups\n", t.record.name)
		
		// take backup
		t.takeBackup()
		
		log.Printf("transaction:: %s has taken backup, starting transaction operation\n", t.record.name)
		
		// apply the transaction's operational logic to the writeSet and readSet members
		operationStatus := t.operation()
		
		if !operationStatus {
			
			// failed to operate successfully, this transaction is flawed,
			// bailing out
			log.Printf("CRITICAL:: transaction:: %s has faulty operational logic, bailing out after rolling back\n", t.record.name)
			
			t.rollback()
			
			t.record.status = false
			
			break
		}
		
		log.Printf("transaction:: %s operation completed, moving to commit changes...\n", t.record.name)
		
		// commit changes
		commitStatus := t.commit()
		
		if !commitStatus {
			
			// failed to commit changes to the writeSet, hence rolling back
			// and then retrying
			log.Printf("MODERATE:: transaction:: %s couldn't commit its changes, rolling back and retrying...\n", t.record.name)
			
			t.rollback()
			
			time.Sleep(time.Duration(MAX_SLEEP_WAIT_TIME))
			
			continue
		}
		
		log.Printf("transaction:: %s has successfully committed its changes made to the writeSet members, marking transaction as completed.", t.record.name)
		
		// since the commit was successful, the transaction releases all its
		// writeSet members of its ownership and marks itself as complete
		t.releaseOwnership()
		
		// marks itself as complete
		t.record.status = true
	}
	
	log.Printf("transaction:: %s has ended...\n", t.record.name)
	
	// notify that the task that was added to waitGroup has been
	// done and hence should removed(?)
	(*waitGroup).Done()
}

// Takes ownership of all the `MemCells` referenced in the transaction's `writeSet`.
func (t *transaction) takeOwnership() bool {
	
	writeSet := &t.record.writeSet
	
	maxOwnershipCount := len(*writeSet)
	
	// the status of the ownership phase, by default we assume that it fails
	// it is only successful if all the writeSet members are successfully
	// owned by the transaction
	status := false
	
	for i := 0; i < len(*writeSet) && maxOwnershipCount > 0; i++ {
		
		// the variableName of the `Variable` that needs to be
		// owned by this transaction
		variableName := (*writeSet)[i]
		currentOwner := t.manager.GetOwner(variableName)
		
		if nil == currentOwner {
			
			// set owner -- successfully took ownership
			t.manager.SetOwner(variableName, t)
			
			log.Printf("transaction:: %s took ownership of Variable:: %s\n", t.record.name, variableName)
			
			maxOwnershipCount--
		}
	}
	
	if maxOwnershipCount <= 0 {
		
		// all the members of the writeSet were owned successfully by this
		// transaction
		status = true
	}
	
	return status
}

// <p>
// Takes the backup of all the members in the read set and write set.
// In case of a rollback, the `<i>writeSet</i>` member contents are
// * restored, from the <i>record.oldValues</i>.
// * <p>
// * While committing, the values of the read set and oldValues is checked, if
// * they are different commit fails.
func (t *transaction) takeBackup() {
	
	writeSet := &t.record.writeSet
	readSet := &t.record.readSet
	
	for i := 0; i < len(*writeSet); i++ {
		
		// Backing up `writeSet` members
		variableName := (*writeSet)[i]
		
		// get the current state of the `Variable` or `MemCell` for creating
		// the backup, the solution is to use fully qualified name of my State struct.
		currentState := t.manager.Read(variableName)
		
		if nil != currentState {
			
			// if we got something for the current state, we need to store
			// it as backup
			t.record.oldValues[variableName] = currentState
		}
	}
	
	for i := 0; i < len(*readSet); i++ {
		
		variableName := (*readSet)[i]
		
		currentState := t.manager.Read(variableName)
		
		if nil != currentState {
			
			t.record.oldValues[variableName] = currentState
		}
	}
}

// <p>
// * Rolls back all changes made by the transaction and releases ownerships of
// * the writeSet members as well.
func (t *transaction) rollback() {
	
	log.Printf("Initiating rollback for transaction: %s\n", t.record.name)
	
	writeSet := &t.record.writeSet
	
	for i := 0; i < len(*writeSet); i++ {
		
		variableName := (*writeSet)[i]
		
		// fetch the backup
		backup := t.record.oldValues[variableName]
		
		// restore the backup
		t.manager.Write(variableName, backup)
	}
	
	// release all  the writeSet members from ownership
	t.releaseOwnership()
	
	log.Printf("Rollback complete for transaction:: %s\n", t.record.name)
}

//  * <p>
//  * Commits the changes made by the transaction to its writeSet members after
//  * referring to the state's of its readSet members.
//  *
//  * @return true if commit was successful else returns false
func (t *transaction) commit() bool {
	
	log.Printf("Initiating commit for transaction:: %s", t.record.name)
	
	status := true
	readSet := &t.record.readSet
	
	for i := 0; i < len(*readSet); i++ {
		
		variableName := (*readSet)[i]
		
		log.Printf("readset mem:: %s \n", variableName)
		
		currentState := t.manager.Read(variableName)
		
		backup := t.record.oldValues[variableName]
		
		log.Println(currentState)
		log.Println(backup)
		
		if nil != currentState {
			
			// there is some non-null state in the statetable.
			if nil == backup {
				
				// this means that when the backup was taken, the readSet
				// member was un-initiaized but now it has some non-null
				// state this means that it has been modified in some way
				// and this
				// might not be good since the new state of the readSet
				// member might cause some
				// problem with states of the writeSet members.
				
				// commit failed
				status = false
				break
			} else if backup != currentState {
				
				// backup is not null, now check if their values are equal
				status = false
				break
			} else {
				
				// backup of the readSet member matches its current state
				status = true
			}
		} else {
			
			// currentstate is empty or null, if old state was not null,
			// then there has been a change in state
			if nil == backup {
				
				// means that both currentState of the readSet member and
				// its backup were empty
				status = true
			} else {
				
				status = false
			}
		}
	}
	
	log.Printf("Completing commit for transaction:: %s\n", t.record.name)
	
	return status
}

// * <p>
// * Releases the ownerships of all the writeSet member `MemCells`
func (t *transaction) releaseOwnership() {
	
	log.Printf("Initiating release of ownership of writeSet members of transaction:: %s\n", t.record.name)
	
	writeSet := &t.record.writeSet
	
	for i := 0; i < len(*writeSet); i++ {
		
		variableName := (*writeSet)[i]
		
		if nil != t.manager.GetOwner(variableName) && t.manager.GetOwner(variableName) == t {
			
			// release ownership only if this transaction owns it
			// this is to prevent race conditions(?)
			t.manager.ReleaseOwnership(variableName)
		}
	}
	
	log.Printf("Finished release of ownership of writeSet members of transaction:: %s\n", t.record.name)
}

/*** Book keeping methods START**/

// addWriteSetMembers : Adds the member <i>Variable</i> or `MemCell`s names to the writeSet of
//  * the transaction.
func (t *transaction) addWriteSetMembers(variableNames ...string) {
	
	for _, variableName := range variableNames {
		
		t.record.writeSet = append(t.record.writeSet, variableName)
	}
}

// addReadSetMembers :
// * Adds the member <i>Variable</i> or `MemCell`s names to the `readSet` of
// * the transaction.
func (t *transaction) addReadSetMembers(variableNames ...string) {
	
	for _, variableName := range variableNames {
		
		if !contains(t.record.writeSet, variableName) {
			
			// since the variables that are needed by the transaction in its
			// writeSet are going to be updated anyways
			// it would be a better idea to have them owned only once, hence the
			// variables that are already a part of the writeSet are not going
			// to be added to the readSet
			t.record.readSet = append(t.record.readSet, variableName)
		}
	}
}

// contains utility method
func contains(l []string, element string) bool {
	
	found := false
	
	for _, e := range l {
		
		if e == element {
			
			found = true
			
			break
		}
	}
	
	return found
}

/*** Book keeping methods END**/

/****** transaction utilities START **********/

// transactions count is the counter for generating IDs
// poorman's uuid xD
var transactionCount int

// Transactions is a utility for building transactions
type Transactions struct {
	t *transaction
}

// NewTransaction the first function in the chain for creating a new transaction
// Creates a new <i>Transaction</i> and sets the description of the
// * transaction, the operational logic and the reference to the
// * <i>StateManager</i> that is in charge of the world.
func (ts *Transactions) NewTransaction(name string) *Transactions {
	ts.t = new(transaction)
	ts.t.record = &record{
		status:    false,
		version:   transactionCount,
		name:      name,
		writeSet:  make([]string, 0),
		readSet:   make([]string, 0),
		oldValues: make(map[string]*State),
	}
	ts.t.manager = Manager()
	return ts
}

// AddWriteSetMembers :
// Adds the <i>Variable</i>s or `MemCell`s to the `writeSet` of the
// * transaction.
func (ts *Transactions) AddWriteSetMembers(variableNames ...string) *Transactions {
	
	ts.t.addWriteSetMembers(variableNames...)
	
	return ts
}

// AddReadSetMembers :
// Adds the <i>Variable</i>s or `MemCell`s to the `readSet` of the
// * transaction.
func (ts *Transactions) AddReadSetMembers(variableNames ...string) *Transactions {
	
	ts.t.addReadSetMembers(variableNames...)
	
	return ts
}

// AddOperationLogic :
// * Adds the transaction's operational logic
func (ts *Transactions) AddOperationLogic(operation func() bool) *Transactions {
	
	ts.t.operation = operation
	
	return ts
}

// Get :
// the terminal operation that gives the transaction instance
func (ts *Transactions) Get() *transaction {
	
	return ts.t
}

/****** Transaction utilities END **********/
