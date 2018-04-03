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
// @last-modified Mon Apr 02 2018 18:43:11 GMT-0700 (PDT)
//

package account

import "github.com/sidmishraw/gostm/stm"

// Account is the representation for an account.
type Account struct {
	Name string  // name of the account
	Amt  float64 // amount of money in the account
}

// NewAccount creates a new account for the given name and initial balance.
func NewAccount(name string, initialAmt float64) Account {
	acc := new(Account)
	acc.Name = name
	acc.Amt = initialAmt
	return *acc
}

// Equals implements the Any interface. Test the equality between 2 accounts.
func (acc1 Account) Equals(acc stm.Any) bool {
	acc2 := acc.(Account)
	if acc1.Name == acc2.Name && acc1.Amt == acc2.Amt {
		return true
	}
	return false
}
