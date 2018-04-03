//
//  BSD 3-Clause License
//
// Copyright (c) 2018, Sidharth Mishra
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//  list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//  this list of conditions and the following disclaimer in the documentation
//  and/or other materials provided with the distribution.
//
// * Neither the name of the copyright holder nor the names of its
//  contributors may be used to endorse or promote products derived from
//  this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// transaction.go
// @author Sidharth Mishra
// @created Fri Mar 30 2018 19:37:11 GMT-0700 (PDT)
// @last-modified Mon Apr 02 2018 19:37:16 GMT-0700 (PDT)
//

package stm

// Transaction is the only way to modify the memory cells in the STM.
type Transaction struct {
	Version         int                     // version of the transaction
	IsComplete      bool                    // flag showing if the transaction is running or is complete
	Action          func(*Transaction) bool // the action that this transaction executes
	readQuarantine  map[*memoryCell]Any     // the read quarantine
	writeQuarantine map[*memoryCell]Any     // the write quarantine
	stm             *STM                    // the reference to the STM this transaction intends to modify
}

// Reads the contents of the memory cell referenced by the `tVar`.
func (t *Transaction) Read(tVar TVar) Any {
	memCell := tVar.(*memoryCell)
	val := t.readQuarantine[memCell]
	if val == nil {
		val = memCell.read()
		t.readQuarantine[memCell] = val
	}
	return val
}

// Writes the new Data into the write quarantine. This will be flushed into the STM upon
// successful commit.
func (t *Transaction) Write(tVar TVar, newData Any) bool {
	memCell := tVar.(*memoryCell)
	t.writeQuarantine[memCell] = newData
	return true
}

// Execute executes this transaction as another thread.
func (t *Transaction) Execute() {
	go t.run()
}

// The actual execution logic of the transaction.
func (t *Transaction) run() {
	t.IsComplete = false
	for !t.IsComplete {
		if status := t.Action(t); !status {
			// failed to execute the action
			t.IsComplete = false
			t.rollback()
			continue
		}
		if status := t.commit(); !status {
			// failed to commit
			t.IsComplete = false
			t.rollback()
			continue
		}
		t.IsComplete = true
	}
	t.Version++
}

// rollback the transaction to the initial state so that it can retry.
func (t *Transaction) rollback() {
	t.readQuarantine = make(map[*memoryCell]Any)
	t.writeQuarantine = make(map[*memoryCell]Any)
}

// commit the write quarantine memory cells into the STM.
func (t *Transaction) commit() bool {
	t.stm.acquireCommitLock()       // acquire the commit lock on the STM
	defer t.stm.releaseCommitLock() // release the commit lock on the STM

	failCount := 0
	for memCell, value := range t.readQuarantine {
		currVal := memCell.read()
		if !currVal.Equals(value) {
			// data has changed by other transaction
			// commit has failed
			failCount++
		}
	}

	if failCount > 0 {
		return false // commit failed
	}

	// read quarantined values have been verified
	for memCell, value := range t.writeQuarantine {
		memCell.write(value) // update the values
	}

	return true // commit succeeded
}
