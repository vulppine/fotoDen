package tool

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadInput(message string) string {
	if message != "" {
		fmt.Printf("%s: ", message)
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	return scanner.Text()
}

func ReadInputAsArray(message string, sep string) []string {
	return strings.Split(ReadInput(message+"[seperator: "+sep+"] "), sep)
}

func ReadInputAsBool(message string, cond string) bool {
	res := ReadInput(message + "[" + cond + "] ")
	if res == cond {
		return true
	}

	return false
}

func ReadInputAsInt(message string) (int, error) {
	var i int
	var err error

	n := ReadInput(message)
	if n == "" {
		return 0, nil
	}

	i, err = strconv.Atoi(n)
	if checkError(err) {
		return 0, err
	}

	return i, nil
}

func ReadInputAsFloat(message string) (float32, error) {
	var f float64
	var err error

	n := ReadInput(message)
	if n == "" {
		return 0, nil
	}

	f, err = strconv.ParseFloat(n, 32)
	if checkError(err) {
		return 0, err
	}

	return float32(f), nil
}
