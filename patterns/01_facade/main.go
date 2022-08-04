package main

import (
	"errors"
	"fmt"
)

// Это шаблон проектирования программного обеспечения, обычно используемый в объектно-ориентированном программировании.
// По аналогии с фасадом в архитектуре, фасад — это объект , который служит внешним интерфейсом, маскирующим более сложный базовый или
// структурный код. Фасад может:
// - улучшить читабельность и удобство использования программной библиотеки , маскируя взаимодействие с более сложными компонентами
// за единым (и часто упрощенным) API
// - предоставить контекстно-зависимый интерфейс для более общей функциональности (в комплекте с контекстно-зависимой проверкой ввода)
// - служить отправной точкой для более широкого рефакторинга монолитных или тесно связанных систем в пользу более слабосвязанного кода

// Шаблон обычно используется, когда
// для доступа к сложной системе требуется простой интерфейс,
// система очень сложна или трудна для понимания,
// точка входа необходима для каждого уровня многоуровневого программного обеспечения, или
// абстракции и реализации подсистемы тесно связаны.
type Product struct {
	Name  string
	Price float64
}

type Shop struct {
	Name     string
	Products []Product
}

func (s *Shop) Sell(user User, product string) error {
	if err := user.Card.CheckBalance(); err != nil {
		return err
	} else {
		fmt.Printf("[Shop] %s has positive balance\n", user.Name)
	}
	for _, prod := range s.Products {
		if prod.Name != product {
			continue
		}
		if prod.Price > user.GetBalance() {
			return errors.New("[Shop] Not enough funds, it so expensive\n")
		}
		fmt.Printf("[Shop] Product %s - sold \n", prod.Name)
	}
	return nil
}

type Bank struct {
	Name  string
	Cards []Card
}

func (b *Bank) UserBalance(cardnumber string) error {
	for _, card := range b.Cards {
		if card.Name != cardnumber {
			continue
		}
		if card.Balance <= 0 {
			return errors.New("[Bank] Not enough funds on the bank card")
		}
	}
	return nil // [Bank] its all right with Balance
}

type Card struct {
	Name    string
	Balance float64
	Bank    *Bank
}

func (c *Card) CheckBalance() error {
	return c.Bank.UserBalance(c.Name)
}

type User struct {
	Name string
	Card *Card
}

func (u *User) GetBalance() float64 {
	return u.Card.Balance
}

var (
	bank = Bank{
		Name:  "Tink",
		Cards: []Card{},
	}

	card1 = Card{
		Name:    "CRD-1",
		Balance: 150,
		Bank:    &bank,
	}
	card2 = Card{
		Name:    "CRD-2",
		Balance: 3,
		Bank:    &bank,
	}
	card3 = Card{
		Name:    "CRD-3",
		Balance: -13,
		Bank:    &bank,
	}

	user1 = User{
		Name: "Tom",
		Card: &card1,
	}
	user2 = User{
		Name: "Sam",
		Card: &card2,
	}
	user3 = User{
		Name: "Bob",
		Card: &card3,
	}

	prod = Product{
		Name:  "Book",
		Price: 120,
	}

	shop = Shop{
		Name: "SHOP",
		Products: []Product{
			prod,
		},
	}
)

func main() {
	bank.Cards = append(bank.Cards, card1, card2, card3)
	if err := shop.Sell(user3, prod.Name); err != nil {
		fmt.Println(err)
		return
	}
}
