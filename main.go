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
// main.go
// @author Sidharth Mishra
// @created Thu Mar 29 2018 00:21:30 GMT-0700 (PDT)
// @last-modified Mon Apr 02 2018 19:53:37 GMT-0700 (PDT)
//

package main

import (
	"fmt"
	"log"

	"github.com/sidmishraw/gostm/account"
	"github.com/sidmishraw/gostm/stm"
)

func main() {
	STM := stm.New()

	acc1 := STM.NewTVar(account.NewAccount("account1", 100))
	acc2 := STM.NewTVar(account.NewAccount("account2", 500))

	// 1st transaction that transfers 100 from account2 to account1
	t := STM.NewTransaction(func(t *stm.Transaction) bool {
		log.Println("Started execution of action ...")
		defer log.Println("Finished execution action ...")

		tacc1 := t.Read(acc1)
		tacc2 := t.Read(acc2)

		log.Println("tacc1 = ", tacc1)
		log.Println("tacc2 = ", tacc2)

		a1 := tacc1.(account.Account) // get the actual account 1 instance
		a2 := tacc2.(account.Account) // get the actual account 2 instance

		log.Println("a1 = ", a1)
		log.Println("a2 = ", a2)

		a1.Amt = a1.Amt + 100 // deposit 100
		a2.Amt = a2.Amt - 100 // withdraw 100

		log.Println("After a1 = ", a1)
		log.Println("After a2 = ", a2)

		return t.Write(acc1, a1) && t.Write(acc2, a2)
	})

	// 2nd transaction that transfers 10 from account1 to account2
	tt := STM.NewTransaction(func(t *stm.Transaction) bool {
		log.Println("Started execution of action ...")
		defer log.Println("Finished execution action ...")

		tacc1 := t.Read(acc1)
		tacc2 := t.Read(acc2)

		log.Println("tacc1 = ", tacc1)
		log.Println("tacc2 = ", tacc2)

		a1 := tacc1.(account.Account) // get the actual account 1 instance
		a2 := tacc2.(account.Account) // get the actual account 2 instance

		log.Println("a1 = ", a1)
		log.Println("a2 = ", a2)

		a1.Amt = a1.Amt - 10 // withdraw 10
		a2.Amt = a2.Amt + 10 // deposit 10

		log.Println("After a1 = ", a1)
		log.Println("After a2 = ", a2)

		return t.Write(acc1, a1) && t.Write(acc2, a2)
	})

	t.Execute()
	tt.Execute()

	var k int
	fmt.Scan(&k)

	// final consistent state is going to be [190, 410].
	tLog := STM.NewTransaction(func(t *stm.Transaction) bool {

		if bal1 := t.Read(acc1).(account.Account); bal1.Amt != 190 {
			log.Println("Failed consistency test")
		}

		if bal2 := t.Read(acc2).(account.Account); bal2.Amt != 410 {
			log.Println("Failed consistency test")
		}

		log.Println("Acc1 = ", t.Read(acc1))
		log.Println("Acc2 = ", t.Read(acc2))

		return true
	})

	tLog.Execute()

	fmt.Scan(&k)
}
