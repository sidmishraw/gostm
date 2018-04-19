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
// memorycell.go
// @author Sidharth Mishra
// @created Thu Mar 29 2018 00:23:26 GMT-0700 (PDT)
// @last-modified Sun Apr 15 2018 09:57:38 GMT-0700 (PDT)
//

package stm

import (
	"fmt"
	"sync"

	"github.com/satori/go.uuid" // for UUID support
)

// memoryCell represents a memory cell where the data is stored.
type memoryCell struct {
	id          string        // The unique identity of the memory cell. Helps in getting it hashed
	data        Value         // The contents of the memory cell.
	memCellLock *sync.RWMutex // A read-write lock for obtaining more granular locking.
}

// newMemCell is a memory cell constructor. It creates and initializes a new memory cell.
func newMemCell(data Value) (memCell *memoryCell) {
	memCell = new(memoryCell)
	memCell.id = uuid.NewV4().String() // generates a new v4 UUID string for ID
	memCell.data = data
	memCell.memCellLock = new(sync.RWMutex)
	return memCell
}

// read the data contained in the memory cell.
// First acquire a read lock on the memory cell and then read the contents.
func (memCell *memoryCell) read() Value {
	memCell.memCellLock.RLock()         // acquire the read lock on the memCell
	defer memCell.memCellLock.RUnlock() // defer the unlock of the memCell lock
	return memCell.data.MakeCopy()      // return the copy of the data
}

// write the newData into the memory cell updating the contents of the memory cell.
// First acquire a write lock and then update the contents.
func (memCell *memoryCell) write(newData Value) {
	memCell.memCellLock.Lock()         // acquire the write lock on the memCell
	defer memCell.memCellLock.Unlock() // defer the release of the memCell lock
	memCell.data = newData             // update the contents of the memory cell
}

// toString gives back a string representation of the memory cell instance.
func (memCell *memoryCell) toString() string {
	// return fmt.Sprintf("MemoryCell#(%s)", memCell.id)
	return fmt.Sprintf(`{"id": %v, "data": %v}`, memCell.id, memCell.data)
}
