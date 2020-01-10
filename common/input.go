package common

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Select ask user to enter the proper interger choice.
// TODO: to move to app code to make this module non interactive.
func Select(from, until int) int {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Enter the value (%d ... %d): ", from, until)
	for {
		ID, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Unable to retrieve the data input. %s", err)
		}

		ID = strings.Trim(ID, "\n")

		if selected, err := strconv.Atoi(ID); err == nil && selected >= from && selected <= until {
			return selected
		}
		fmt.Printf("Your choice must be between %d and %d. You entered '%s'. Please select proper one.\n", from, until, ID)
		fmt.Print("Enter the value: ")
	}
}

// GetNumber ask to enter a number
func GetNumber(mess string) (selected int) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(mess)
	for {
		ID, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Unable to retrieve the data input. %s", err)
		}
		selected, err = strconv.Atoi(strings.Trim(ID, " \n"))
		if err == nil {
			return selected
		}
		fmt.Printf("Please enter a number. %s\n", err)
		fmt.Print(mess)
	}
}
