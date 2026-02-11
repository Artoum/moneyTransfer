package main

import "fmt"

type User struct {
	Id      int
	Name    string
	Balance float64
}

func (u *User) deposit(amount float64) {
	u.Balance += amount
	fmt.Printf("Баланс пользователя %v успешно пополнен на %.2f \nНа счету у %v: %.2f \n",
		u.Name, amount, u.Name, u.Balance)
}

func (u *User) withdraw(amount float64) error {
	if u.Balance < amount {
		return fmt.Errorf("У пользователя %v недостаточно средств для снятия или перевода.", u.Name)
	}
	u.Balance -= amount
	fmt.Printf("Со счета %v успешно снято %.2f\nОстаток средств у %v: %2.f \n",
		u.Name, amount, u.Name, u.Balance)

	return nil
}

func main() {
	many := []*User{}

	alice := User{101, "Alice", 0}
	tom := User{102, "Tom", 0}

	many = append(many, &alice, &tom)

	alice.deposit(1000)
	tom.deposit(1500)

	err := alice.withdraw(1500)
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}
	err = tom.withdraw(1000)
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}
}
