package models

type ShoppingItem struct {
	ID      int
	Name    string
	StoreID int // ID магазина, к которому привязан товар
}
