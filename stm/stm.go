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
// stm.go
// @author Sidharth Mishra
// @created Thu Mar 29 2018 00:22:14 GMT-0700 (PDT)
// @last-modified Sat Apr 14 2018 11:00:55 GMT-0700 (PDT)
//

package stm

import (
	"log"
	"sync"
)

// STM is the single shared memory store that can only be modified by transactions.
type STM struct {
	memory     []*memoryCell // the collection of memory cells makes up the memory
	commitLock *sync.Mutex   // the commit lock needed for maintaining consistency -- serializability
}

// New makes and initializes a new STM instance.
func New() (stm *STM) {
	stm = new(STM)
	stm.memory = make([]*memoryCell, 0, 0)
	stm.commitLock = new(sync.Mutex)
	return stm
}

// NewTVar creates a new memory cell in the STM and returns the reference
// to the memory cell as a TVar instance.
func (stm *STM) NewTVar(data Value) TVar {
	memCell := newMemCell(data)
	stm.memory = append(stm.memory, memCell)
	return TVar(memCell)
}

// Perform accepts the transactional actions submitted to the STM and performs them.
func (stm *STM) Perform(actions ...func(*Transaction) bool) {
	for _, action := range actions {
		t := newTransaction(stm, action)
		t.execute()
	}
}

// newTransaction makes a new transaction for the given action.
func newTransaction(stm *STM, action func(*Transaction) bool) *Transaction {
	t := new(Transaction)
	t.version = 0
	t.isComplete = false
	t.action = action
	t.readQuarantine = make(map[*memoryCell]Value)
	t.writeQuarantine = make(map[*memoryCell]Value)
	t.stm = stm
	return t
}

// acquireCommitLock acquires the lock on the commit lock of this STM.
func (stm *STM) acquireCommitLock() {
	stm.commitLock.Lock()
}

// releaseCommitLock releases the lock on the commit lock of this STM.
func (stm *STM) releaseCommitLock() {
	stm.commitLock.Unlock()
}

// PrintState prints the current state of all the memory cells -- just a snapshot
func (stm *STM) PrintState() {
	log.Println("------- State of the STM -------- ")
	for _, memCell := range stm.memory {
		log.Println(memCell.toString())
	}
	log.Println("------- END -------- ")
}
