package main

import (
	"errors"
	"fmt"
	"sync"
)

type User struct {
	Id      string
	Name    string
	Balance float64
}

type Transaction struct {
	FromID string
	ToID   string
	Amount float64
}

type PaymentSystem struct {
	Users            map[string]*User
	TransactionQueue []Transaction
}

func (u *User) Deposit(amount float64) {
	u.Balance += amount
	fmt.Printf("Баланс пользователя %v успешно пополнен на %.2f \nНа счету у %v: %.2f \n",
		u.Name, amount, u.Name, u.Balance)
}

func (u *User) Withdraw(amount float64) error {
	if u.Balance < amount {
		return fmt.Errorf("У пользователя %v недостаточно средств для снятия или перевода.", u.Name)
	}
	u.Balance -= amount
	fmt.Printf("Со счета %v успешно снято %.2f\nОстаток средств у %v: %2.f \n",
		u.Name, amount, u.Name, u.Balance)

	return nil
}

func (ps *PaymentSystem) AddUser(user *User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if ps.Users == nil {
		ps.Users = make(map[string]*User)
	}

	ps.Users[user.Id] = user
	return nil
}

func (ps *PaymentSystem) AddTransaction(t Transaction) {
	if ps.TransactionQueue == nil {
		ps.TransactionQueue = make([]Transaction, 0)
	}

	ps.TransactionQueue = append(ps.TransactionQueue, t)
}

func (ps *PaymentSystem) ProcessTransactions(t Transaction) error {
	fromUser, exists := ps.Users[t.FromID]
	if !exists {
		return fmt.Errorf("пользователь %s не найден", t.FromID)
	}
	err := fromUser.Withdraw(t.Amount)
	if err != nil {
		return fmt.Errorf("Ошибка списания у %s %w: ", t.FromID, err)
	}

	toUser, exists := ps.Users[t.ToID]
	if !exists {
		return fmt.Errorf("пользователь %s не найден", t.ToID)
	}
	toUser.Deposit(t.Amount)
	return nil
}

func (ps *PaymentSystem) Worker(ch <-chan Transaction, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range ch {
		ps.ProcessTransactions(t)
	}
}

func main() {
	ps := &PaymentSystem{
		Users:            make(map[string]*User),
		TransactionQueue: []Transaction{},
	}

	var wg sync.WaitGroup
	ch := make(chan Transaction, len(ps.TransactionQueue))

	user1 := User{Id: "1", Name: "Alice", Balance: 1000}
	user2 := User{Id: "2", Name: "Bob", Balance: 500}

	ps.AddUser(&user1)
	ps.AddUser(&user2)

	ps.AddTransaction(Transaction{FromID: "1", ToID: "2", Amount: 200})
	ps.AddTransaction(Transaction{FromID: "2", ToID: "1", Amount: 50})

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go ps.Worker(ch, &wg)
	}

	for _, t := range ps.TransactionQueue {
		ch <- t
	}
	close(ch)
	wg.Wait()

	ps.TransactionQueue = []Transaction{}

	fmt.Printf("Баланс %v: %.2f\n", user1.Name, user1.Balance)
	fmt.Printf("Баланс %v: %.2f\n", user2.Name, user2.Balance)
}
