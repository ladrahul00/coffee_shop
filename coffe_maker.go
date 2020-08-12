package coffemaker

import (
	"fmt"
	"sync"
	"time"
)

// Machine type
type Machine struct {
	Outlets    map[string]int            `json:"outlets"`
	ItemsStock map[string]int            `json:"total_items_quantity"`
	Beverages  map[string]map[string]int `json:"beverages"`
}

// MachineInput type
type MachineInput struct {
	Machine Machine `json:"machine"`
}

// CoffeeMachine type
type CoffeeMachine struct {
	beverages map[string]map[string]int
	outlets   int
	stock     map[string]int
	lock      *sync.Mutex
}

// Init coffee machine
func (cm *CoffeeMachine) Init(beverages map[string]map[string]int, outlets map[string]int, stock map[string]int) error {
	cm.beverages = beverages
	outletCount, ok := outlets["count_n"]
	if !ok {
		return fmt.Errorf("invalid json")
	}
	cm.outlets = outletCount
	cm.stock = stock
	cm.lock = &sync.Mutex{}
	return nil
}

// Serve the drink
func (cm *CoffeeMachine) Serve(drink string) string {
	ingredients, available := cm.beverages[drink]
	if !available {
		return fmt.Sprintf("%s is not available", drink)
	}
	cm.lock.Lock()
	defer cm.lock.Unlock()
	for item, amount := range ingredients {
		totalAmount, available := cm.stock[item]
		if !available {
			return fmt.Sprintf("%s cannot be prepared because %s is not available", drink, item)
		}
		if totalAmount < amount {
			return fmt.Sprintf("%s cannot be prepared because item %s is not sufficient", drink, item)
		}
	}
	for item, amount := range ingredients {
		totalAmount, _ := cm.stock[item]
		totalAmount = totalAmount - amount
		cm.stock[item] = totalAmount
	}
	// Beverage preparation time
	// Send machine instructions to dispatch drink
	time.Sleep(time.Second * 1)

	return fmt.Sprintf("%s is prepared", drink)
}

// Refill ingredients
func (cm *CoffeeMachine) Refill(ingredient string, amount int) {
	existingAmount, ok := cm.stock[ingredient]
	if !ok {
		cm.stock[ingredient] = amount
	} else {
		cm.stock[ingredient] = existingAmount + amount
	}
}
