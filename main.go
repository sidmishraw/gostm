/**
* main.go
*
* @author Sidharth Mishra
* @description
* @created Thu Oct 05 2017 00:24:16 GMT-0700 (PDT)
* @copyright
* @last-modified Thu Oct 05 2017 19:58:22 GMT-0700 (PDT)
 */

package main

import "github.com/sidmishraw/ggfirst/core"
import "log"
import (
	"sync"
)

func main() {
	
	//test1()
	test3()
}

// first single threaded test for the stm
func test1() {
	
	// for demo, accountbalance
	type AccountBalance struct {
		Balance float64
	}
	
	// a sample deposit function
	deposit := func(accountName string, amount float64) {
		
		log.Printf("Depositing amount:: %f into bank account:: %s\n", amount, accountName)
		
		manager := core.Manager()
		
		b := *(manager.Read(accountName))
		
		if nil == b {
			b = AccountBalance{Balance: 0.0}
		}
		
		bNew := core.State(AccountBalance{Balance: b.(AccountBalance).Balance + amount})
		
		manager.Write(accountName, &bNew)
		
		log.Printf("Deposited amount:: %f into bank account:: %s, new amount:: %f\n", amount, accountName, bNew.(AccountBalance).Balance)
	}
	
	// a sample withdraw function
	withdraw := func(accountName string, amount float64) {
		
		log.Printf("Withdrawing amount:: %f from bank account:: %s\n", amount, accountName)
		
		manager := core.Manager()
		
		b := *(manager.Read(accountName))
		
		if nil == b {
			b = AccountBalance{Balance: 0.0}
		}
		
		bNew := core.State(AccountBalance{Balance: b.(AccountBalance).Balance - amount})
		
		manager.Write(accountName, &bNew)
		
		log.Printf("Withdrew amount:: %f into bank account:: %s, new amount:: %f\n", amount, accountName, bNew.(AccountBalance).Balance)
	}
	
	setupbankaccounts := func() {
		
		manager := core.Manager()
		
		if manager == nil {
			println("aaa")
		}
		// create bank accounts to operate on
		manager.MakeVariable("Account#1")
		manager.MakeVariable("Account#2")
		
		// make initial states
		b1 := core.State(AccountBalance{Balance: 500.0})
		b2 := core.State(AccountBalance{Balance: 1500.0})
		manager.Write("Account#1", &b1)
		manager.Write("Account#2", &b2)
	}
	
	log.Printf("Initiating bankDriver...")
	
	ts := new(core.Transactions)
	
	setupbankaccounts()
	
	m := core.Manager()
	
	log.Printf("Initially:: Account#1:: %f and Account#2:: %f", (*m.Read("Account#1")).(AccountBalance).Balance, (*m.Read("Account#2")).(AccountBalance).Balance)
	
	t := ts.NewTransaction("T1").
		AddWriteSetMembers("Account#1", "Account#2").
		AddReadSetMembers("Account#1", "Account#2").
		AddOperationLogic(
		func() bool {
			deposit("Account#1", 500.0)
			withdraw("Account#2", 500.0)
			
			return true
		}).
		Get()
	
	// using a waitgroup for waiting
	// similar to Java's Thread.join()
	wg := sync.WaitGroup{}
	
	wg.Add(1)
	
	go t.Run(&wg)
	
	wg.Wait()
	
	log.Printf("Finally:: Account#1:: %f and Account#2:: %f\n", (*m.Read("Account#1")).(AccountBalance).Balance, (*m.Read("Account#2")).(AccountBalance).Balance)
}

// 3 threaded test for the stm
func test3() {
	
	// for demo, accountbalance
	type AccountBalance struct {
		Balance float64
	}
	
	// a sample deposit function
	deposit := func(accountName string, amount float64) {
		
		log.Printf("Depositing amount:: %f into bank account:: %s\n", amount, accountName)
		
		manager := core.Manager()
		
		b := *(manager.Read(accountName))
		
		if nil == b {
			b = AccountBalance{Balance: 0.0}
		}
		
		bNew := core.State(AccountBalance{Balance: b.(AccountBalance).Balance + amount})
		
		manager.Write(accountName, &bNew)
		
		log.Printf("Deposited amount:: %f into bank account:: %s, new amount:: %f\n", amount, accountName, bNew.(AccountBalance).Balance)
	}
	
	// a sample withdraw function
	withdraw := func(accountName string, amount float64) {
		
		log.Printf("Withdrawing amount:: %f from bank account:: %s\n", amount, accountName)
		
		manager := core.Manager()
		
		b := *(manager.Read(accountName))
		
		if nil == b {
			b = AccountBalance{Balance: 0.0}
		}
		
		bNew := core.State(AccountBalance{Balance: b.(AccountBalance).Balance - amount})
		
		manager.Write(accountName, &bNew)
		
		log.Printf("Withdrew amount:: %f into bank account:: %s, new amount:: %f\n", amount, accountName, bNew.(AccountBalance).Balance)
	}
	
	setupbankaccounts := func() {
		
		manager := core.Manager()
		
		// create bank accounts to operate on
		manager.MakeVariable("Account#1")
		manager.MakeVariable("Account#2")
		manager.MakeVariable("Account#3")
		
		// make initial states
		b1 := core.State(AccountBalance{Balance: 500.0})
		b2 := core.State(AccountBalance{Balance: 1500.0})
		b3 := core.State(AccountBalance{Balance: 100.0})
		
		manager.Write("Account#1", &b1)
		manager.Write("Account#2", &b2)
		manager.Write("Account#3", &b3)
	}
	
	log.Printf("Initiating bankDriver...")
	
	ts := new(core.Transactions)
	
	setupbankaccounts()
	
	m := core.Manager()
	
	log.Printf("Initially:: Account#1:: %f and Account#2:: %f", (*m.Read("Account#1")).(AccountBalance).Balance, (*m.Read("Account#2")).(AccountBalance).Balance)
	
	t := ts.NewTransaction("T1").
		AddWriteSetMembers("Account#1", "Account#2").
		AddReadSetMembers("Account#1", "Account#2").
		AddOperationLogic(
		func() bool {
			deposit("Account#1", 500.0)
			withdraw("Account#2", 500.0)
			
			return true
		}).
		Get()
	
	t2 := ts.NewTransaction("T2").
		AddWriteSetMembers("Account#1", "Account#2").
		AddReadSetMembers("Account#1", "Account#2").
		AddOperationLogic(
		func() bool {
			deposit("Account#1", 50.0)
			withdraw("Account#2", 50.0)
			
			return true
		}).
		Get()
	
	// Read from A1 and if bal > 100 the value transfer 500 from A2 -> A3
	t3 := ts.NewTransaction("T3").
		AddWriteSetMembers("Account#2", "Account#3").
		AddReadSetMembers("Account#1", "Account#2", "Account#3").
		AddOperationLogic(
		func() bool {
			
			if b1 := (*(core.Manager().Read("Account#1"))).(AccountBalance).Balance; b1 > 100.0 {
				
				deposit("Account#3", 500.0)
				withdraw("Account#2", 500.0)
			}
			
			return true
		}).
		Get()
	
	// using a waitgroup for waiting
	// similar to Java's Thread.join()
	wg := sync.WaitGroup{}
	
	wg.Add(3)
	
	go t.Run(&wg)
	go t2.Run(&wg)
	go t3.Run(&wg)
	
	wg.Wait()
	
	log.Printf("Finally:: Account#1:: %f and Account#2:: %f and Account#3:: %f\n", (*m.Read("Account#1")).(AccountBalance).Balance, (*m.Read("Account#2")).(AccountBalance).Balance, (*m.Read("Account#3")).(AccountBalance).Balance)
}
