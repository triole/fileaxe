package fileaxe

import (
	"fmt"
	"log"
)

func (fa FileAxe) askForConfirmation(filename string, action string) bool {
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	fmt.Printf(
		"Type %s to confirm %s of %q  ",
		okayResponses, action, filename,
	)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		return fa.askForConfirmation(filename, action)
	}
}

func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}
