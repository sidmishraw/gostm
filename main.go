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
// @last-modified Thu Apr 19 2018 12:50:58 GMT-0700 (PDT)
//

package main

import (
	"fmt"

	"github.com/sidmishraw/gostm/account"
	"github.com/sidmishraw/gostm/stm"
)

func main() {
	STM := stm.New()

	acc1 := account.NewAccount("account1", 100, STM)
	acc2 := account.NewAccount("account1", 500, STM)

	// Initial state of the STM [100, 500].
	//
	STM.PrintState()

	// 1st transaction that transfers 100 from account2 to account1
	//
	acc2.Transfer(acc1, 100)

	// 2nd transaction that transfers 10 from account1 to account2
	//
	acc1.Transfer(acc2, 10)

	var k int
	fmt.Scan(&k)

	// final consistent state is going to be [190, 410].
	//
	STM.PrintState()
}
