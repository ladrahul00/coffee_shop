package coffemaker

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

// InitTest ...
func InitTest(t *testing.T, fileName string) CoffeeMachine {
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

func containsInArray(arr []string, searchString string) bool {
	for _, str := range arr {
		if str == searchString {
			return true
		}
	}
	return false
}

type TestCase struct {
	coffeeMachine    CoffeeMachine
	drink            string
	expectedMessages []string
}

func serveDrink(testCase TestCase, t *testing.T) {
	msg := make(chan string)
	go testCase.coffeeMachine.Serve(testCase.drink, msg)
	actualMsg := <-msg
	if !containsInArray(testCase.expectedMessages, actualMsg) {
		t.Errorf("serving %s FAILED, \nexpected: '%s', \nactual: '%s'", testCase.drink, testCase.expectedMessages[0], actualMsg)
	} else {
		t.Logf("serving %s PASSED, \nmessage: '%s'", testCase.drink, actualMsg)
	}
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
			"hot_coffee is prepared",
			"hot_coffee cannot be prepared because item hot_milk is not sufficient",
		},
	})
	testCases = append(testCases, TestCase{
		coffeeMachine: coffeMachine,
		drink:         "black_tea",
		expectedMessages: []string{
			"black_tea cannot be prepared because item hot_water is not sufficient",
			"black_tea cannot be prepared because item sugar_syrup is not sufficient",
			"black_tea is prepared",
		},
	})
	return testCases
}

// TestServeDrink ...
func TestServeDrinkSync(t *testing.T) {
	coffeMachine := InitTest(t, "input_data1.json")
	testCases := prepareTestCasesForInput1Sync(coffeMachine)
	for _, testCase := range testCases {
		serveDrink(testCase, t)
	}
}

// TestServeDrinkWithRefill ...
func TestServeDrinkWithRefillSync(t *testing.T) {
	coffeMachine := InitTest(t, "input_data1.json")
	testCases := prepareTestCasesForInput1Sync(coffeMachine)
	for _, testCase := range testCases {
		serveDrink(testCase, t)
	}
	prevHotWaterStock, ok := coffeMachine.stock["hot_water"]
	if !ok {
		t.Error("hot_water not found")
	}
	coffeMachine.Refill("hot_water", 500)
	currentHotWaterStock, ok := coffeMachine.stock["hot_water"]
	if !ok {
		t.Error("hot_water not found")
	}
	if currentHotWaterStock != (prevHotWaterStock + 500) {
		t.Error("refill FAILED: current hot_water stock does not match expected stock")
	} else {
		t.Log("refill SUCCESS: hot_water stock refill successfull")
		coffeMachine.Refill("sugar_syrup", 500)
	}

	serveDrink(TestCase{
		coffeeMachine:    coffeMachine,
		drink:            "black_tea",
		expectedMessages: []string{"black_tea is prepared"},
	}, t)
}

// TestServeDrinkAsync ...
func TestServeDrinkAsync(t *testing.T) {
	coffeMachine := InitTest(t, "input_data1.json")

	testCases := prepareTestCasesForInput1ASync(coffeMachine)
	for _, testCase := range testCases {
		t.Run("testRunAsync", func(t *testing.T) {
			t.Parallel()
			serveDrink(testCase, t)
		})
	}
}
