package coffemaker

import (
	"fmt"
	"sync"
	"time"
)

// CoffeeMachine type
type CoffeeMachine struct {
	beverages map[string]map[string]int
	lock      *sync.Mutex
	outlets   chan string
	stock     map[string]int
}

// Init coffee machine
func (cm *CoffeeMachine) Init(beverages map[string]map[string]int, outlets map[string]int, stock map[string]int) error {
	cm.beverages = beverages
	outletCount, ok := outlets["count_n"]
	if !ok {
		return fmt.Errorf("invalid json")
	}
	// Creating bufferred channel for number of outlets
	cm.outlets = make(chan string, outletCount)
	cm.stock = stock
	// Initiazling mutex lock to synchronize stock consumption and refill
	cm.lock = &sync.Mutex{}
	return nil
}

// Refill ingredients
func (cm *CoffeeMachine) Refill(ingredient string, amount int) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	// Synchronized block
	existingAmount, ok := cm.stock[ingredient]
	if !ok {
		cm.stock[ingredient] = amount
	} else {
		cm.stock[ingredient] = existingAmount + amount
	}
}

// ServeFromOutlet ...
func (cm *CoffeeMachine) ServeFromOutlet(drink string) string {
	go cm.serve(drink, cm.outlets)
	return <-cm.outlets
}

// Serve the drink
func (cm *CoffeeMachine) serve(drink string, displayMessage chan string) {
	message, err := cm.validateAndReserveDrinkIngredients(drink)

	// The Mutex Lock has been released as the ingredients are already reserved for the drink

	if err != nil {
		displayMessage <- err.Error()
		return
	}
	// Beverage preparation time
	time.Sleep(time.Second * 1)
	displayMessage <- message
}

func (cm *CoffeeMachine) validateAndReserveDrinkIngredients(drink string) (string, error) {
	ingredients, available := cm.beverages[drink]
	if !available {
		return "", fmt.Errorf("%s is not available", drink)
	}
	cm.lock.Lock()
	defer cm.lock.Unlock()
	// Synchronized block
	for item, amount := range ingredients {
		totalAmount, available := cm.stock[item]
		if !available {
			return "", fmt.Errorf("%s cannot be prepared because %s is not available", drink, item)
		}
		if totalAmount < amount {
			return "", fmt.Errorf("%s cannot be prepared because item %s is not sufficient", drink, item)
		}
	}
	for item, amount := range ingredients {
		totalAmount, _ := cm.stock[item]
		totalAmount = totalAmount - amount
		cm.stock[item] = totalAmount
	}
	return fmt.Sprintf("%s is prepared", drink), nil
}
