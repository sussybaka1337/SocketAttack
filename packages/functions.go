package packages

import (
	"bufio"
	cryptoRand "crypto/rand"
	"fmt"
	rand "math/rand"
	"net/url"
	"os"
	"strconv"
)

func ToInt(value string) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return result
}

func ToStr(value int) string {
	return strconv.Itoa(value)
}

func ParseTarget(target string) *url.URL {
	targetURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}
	if targetURL.Path == "" {
		targetURL.Path = "/"
	}
	return targetURL
}

func ReadLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(path + "not found")
	}
	defer file.Close()
	var values []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value := scanner.Text()
		values = append(values, value)
	}
	return values
}

func RandValue(values []string) string {
	length := len(values)
	accessIndex := rand.Intn(length)
	return values[accessIndex]
}

func Range(start int, end int) int {
	return rand.Intn(end-start+1) + start
}

func RandString(length int) string {
	buffer := make([]byte, length)
	cryptoRand.Read(buffer)
	return fmt.Sprintf("%x", buffer)[:length]
}
