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
// account.go
// @author Sidharth Mishra
// @created Fri Mar 30 2018 19:18:42 GMT-0700 (PDT)
// @last-modified Sat Apr 14 2018 10:58:47 GMT-0700 (PDT)
//

package account

import (
	"fmt"

	"github.com/sidmishraw/gostm/stm"
)

// ------------------------------------------------------------------------

// Account is the representation for an account.
// It is the domain object. We split it into identity and state.
type Account struct {
	Details *details // bank account details -- identity
	state   stm.TVar // account state -- stored and managed by the STM
	stm     *stm.STM // the STM managing the state of this account
}

// details is the account details -- immutable or identity of the account domain object.
type details struct {
	Name string // name of the acccount
}

// state is the state of the account.
type state struct {
	amt int // amount of money in the account
}

// newAccState creates a new account state.
func newAccState(amt int) *state {
	s := new(state)
	s.amt = amt
	return s
}

// ------------------------------------------------------------------------

// MakeCopy makes state conform to the stm.Value interface.
// Hence, account's state can be stored and managed by the STM.
func (s *state) MakeCopy() stm.Value {
	ns := new(state)
	ns.amt = s.amt
	return ns
}

// IsEqual checks the equality between two States.
// This implementation makes the state type conform to the stm.Value interface,
// enabling it to be stored and managed by the STM.
func (s *state) IsEqual(v stm.Value) bool {
	vv, ok := v.(*state) // type assert state is the Value
	if !ok {
		return false
	}

	if vv.amt == s.amt {
		return true
	}

	return false
}

// NewAccount creates a new account for the given name and initial balance.
func NewAccount(name string, initialAmt int, stm *stm.STM) *Account {
	acc := new(Account)

	acc.Details = new(details)
	acc.Details.Name = name
	acc.stm = stm

	acc.state = acc.stm.NewTVar(newAccState(initialAmt))

	return acc
}

// Deposit adds the amount to the account's current balance, resulting in
// increasing the current balance. It is an operation that modifies the
// account's state. Hence, it must be delegated to the STM as it is managing
// the state of this account.
func (acc *Account) Deposit(amt int) {
	// The task of updating the state has been delegated to the
	// STM that is managing the state. This ensures consistency,
	// and atomicity.
	//
	acc.stm.Perform(func(t *stm.Transaction) bool {
		accState := t.Read(acc.state).(*state)

		accState.amt = accState.amt + amt

		return t.Write(acc.state, accState)
	})
}

// Withdraw removes the amount from the account's current balance, resulting in
// decrease in the current balance. It is an operation that modifies the
// account's state. Hence, it must be delegated to the STM as it is managing
// hte account's state.
func (acc *Account) Withdraw(amt int) {
	acc.stm.Perform(func(t *stm.Transaction) bool {
		accState := t.Read(acc.state).(*state)

		if accState.amt < amt {
			return false
		}
		accState.amt = accState.amt - amt

		return t.Write(acc.state, accState)
	})
}

// Transfer transfers the desired amount from this account to the destination account.
// Since, this operation is only complete when both the accounts have been modified/updated
// it needs to be atomic. Moreover, as this operation modifies the states of the accounts
// it is delegated to the STM.
func (acc *Account) Transfer(dest *Account, amt int) {
	acc.stm.Perform(func(t *stm.Transaction) bool {
		srcState := t.Read(acc.state).(*state)
		destState := t.Read(dest.state).(*state)

		if srcState.amt < amt {
			return false
		}

		srcState.amt = srcState.amt - amt
		destState.amt = destState.amt + amt

		return t.Write(acc.state, srcState) && t.Write(dest.state, destState)
	})
}

// ToString gives a string representation of the account, just used for debugging.
func (acc *Account) ToString() string {
	return fmt.Sprintf(`{details: "%v", state: "%v"}`, acc.Details, acc.state)
}

// TODO:: Code review.
// TODO:: Check conformance with the QSTM pattern.
