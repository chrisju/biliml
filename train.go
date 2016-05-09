package main

import (
	"fmt"
	"strconv"

	"github.com/sjwhitworth/golearn/base"
	. "github.com/sjwhitworth/golearn/linear_models"
)

func errexit(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	rawData, err := base.ParseCSVToInstances("data2.csv", true)
	errexit(err)

	lr := NewLinearRegression()

	//Do a training-test split
	trainData, testData := base.InstancesTrainTestSplit(rawData, 0.60)

	err1 := lr.Fit(trainData)
	errexit(err1)

	predictions, err2 := lr.Predict(testData)
	errexit(err2)

	_, rows := predictions.Size()

	total := 0.0
	m := 0.0
	//n := 0.0
	for i := 0; i < rows; i++ {
		actualValue, _ := strconv.ParseFloat(base.GetClass(testData, i), 64)
		expectedValue, _ := strconv.ParseFloat(base.GetClass(predictions, i), 64)

		if expectedValue <= 0 && actualValue == 0 {
			continue
		}
		if expectedValue < 0 {
			expectedValue = 0
		}
		if expectedValue > 200 || actualValue > 200 {
			d := expectedValue / actualValue
			if d > 1 {
				d = 1 / d
			}
			total += d
			m++
			fmt.Println(expectedValue, actualValue)
		}
	}
	fmt.Println(total / m)
}
