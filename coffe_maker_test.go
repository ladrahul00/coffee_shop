package coffemaker

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

type Machine struct {
	Outlets    map[string]int            `json:"outlets"`
	ItemsStock map[string]int            `json:"total_items_quantity"`
	Beverages  map[string]map[string]int `json:"beverages"`
}

type MachineInput struct {
	Machine Machine `json:"machine"`
}

type TestCase struct {
	coffeeMachine    CoffeeMachine
	drink            string
	expectedMessages []string
}

func containsInArray(arr []string, searchString string) bool {
	for _, str := range arr {
		if str == searchString {
			return true
		}
	}
	return false
}

func initTest(t *testing.T, fileName string) CoffeeMachine {
	coffeMachine := CoffeeMachine{}
	dataBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	data := MachineInput{}
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	coffeMachine.Init(data.Machine.Beverages, data.Machine.Outlets, data.Machine.ItemsStock)

	t.Log("Coffee Machine Init Successfull")
	return coffeMachine
}

func prepareTestCasesForInput1Sync(coffeMachine CoffeeMachine) []TestCase {
	testCases := []TestCase{}
	testCases = append(testCases, TestCase{
		coffeeMachine:    coffeMachine,
		drink:            "green_tea",
		expectedMessages: []string{"green_tea cannot be prepared because green_mixture is not available"},
	})
	testCases = append(testCases, TestCase{
		coffeeMachine:    coffeMachine,
		drink:            "hot_tea",
		expectedMessages: []string{"hot_tea is prepared"},
	})
	testCases = append(testCases, TestCase{
		coffeeMachine:    coffeMachine,
		drink:            "hot_coco",
		expectedMessages: []string{"hot_coco is not available"},
	})
	testCases = append(testCases, TestCase{
		coffeeMachine:    coffeMachine,
		drink:            "hot_coffee",
		expectedMessages: []string{"hot_coffee is prepared"},
	})
	testCases = append(testCases, TestCase{
		coffeeMachine: coffeMachine,
		drink:         "black_tea",
		expectedMessages: []string{
			"black_tea cannot be prepared because item hot_water is not sufficient",
			"black_tea cannot be prepared because item sugar_syrup is not sufficient",
		},
	})
	return testCases
}

func prepareTestCasesForInput1ASync(coffeMachine CoffeeMachine) []TestCase {
	testCases := []TestCase{}
	testCases = append(testCases, TestCase{
		coffeeMachine:    coffeMachine,
		drink:            "green_tea",
		expectedMessages: []string{"green_tea cannot be prepared because green_mixture is not available"},
	})
	testCases = append(testCases, TestCase{
		coffeeMachine:    coffeMachine,
		drink:            "hot_tea",
		expectedMessages: []string{"hot_tea is prepared"},
	})
	testCases = append(testCases, TestCase{
		coffeeMachine:    coffeMachine,
		drink:            "hot_coco",
		expectedMessages: []string{"hot_coco is not available"},
	})
	testCases = append(testCases, TestCase{
		coffeeMachine: coffeMachine,
		drink:         "hot_coffee",
		expectedMessages: []string{
			// Adding all expected messages as this test case is going to to async
			"hot_coffee is prepared",
			"hot_coffee cannot be prepared because item hot_milk is not sufficient",
		},
	})
	testCases = append(testCases, TestCase{
		coffeeMachine: coffeMachine,
		drink:         "black_tea",
		expectedMessages: []string{
			// Adding all expected messages as this test case is going to to async
			"black_tea cannot be prepared because item hot_water is not sufficient",
			"black_tea cannot be prepared because item sugar_syrup is not sufficient",
			"black_tea is prepared",
		},
	})
	return testCases
}

func refillStock(coffeMachine CoffeeMachine, ingredient string, amount int, t *testing.T) {
	previousStock, ok := coffeMachine.stock[ingredient]
	if !ok {
		t.Errorf("%s not found", ingredient)
	}
	coffeMachine.Refill(ingredient, amount)
	currentStock, ok := coffeMachine.stock[ingredient]
	if !ok {
		t.Errorf("%s not found", ingredient)
	}
	expectedAmount := previousStock + amount
	if currentStock != expectedAmount {
		t.Errorf("refill FAILED: current %s stock amount does not match expected stock amount", ingredient)
	} else {
		t.Logf("refill SUCCESS: %s stock refill successfull", ingredient)
	}
}

func serveDrink(testCase TestCase, t *testing.T) {
	actualMsg := testCase.coffeeMachine.ServeFromOutlet(testCase.drink)
	if !containsInArray(testCase.expectedMessages, actualMsg) {
		t.Errorf("serving %s FAILED, \nexpected: '%s', \nactual: '%s'", testCase.drink, testCase.expectedMessages[0], actualMsg)
	} else {
		t.Logf("serving %s PASSED, \nmessage: '%s'", testCase.drink, actualMsg)
	}
}

// TestServeDrink ...
func TestServeDrinkSync(t *testing.T) {
	coffeMachine := initTest(t, "input_data1.json")
	testCases := prepareTestCasesForInput1Sync(coffeMachine)
	for _, testCase := range testCases {
		serveDrink(testCase, t)
	}
}

// TestServeDrinkWithRefill ...
func TestServeDrinkWithRefillSync(t *testing.T) {
	coffeMachine := initTest(t, "input_data1.json")
	testCases := prepareTestCasesForInput1Sync(coffeMachine)
	for _, testCase := range testCases {
		serveDrink(testCase, t)
	}
	refillStock(coffeMachine, "hot_water", 500, t)
	refillStock(coffeMachine, "sugar_syrup", 500, t)
	serveDrink(TestCase{
		coffeeMachine:    coffeMachine,
		drink:            "black_tea",
		expectedMessages: []string{"black_tea is prepared"},
	}, t)
}

// TestServeDrinkAsync ...
func TestServeDrinkAsync(t *testing.T) {
	coffeMachine := initTest(t, "input_data1.json")

	testCases := prepareTestCasesForInput1ASync(coffeMachine)
	for _, testCase := range testCases {
		t.Run("testRunAsync", func(t *testing.T) {
			t.Parallel()
			serveDrink(testCase, t)
		})
	}
}
