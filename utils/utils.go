package utils

import (
    "os"
    "bufio"
    "fmt"
    "io/ioutil"
)

func ReadFirstLineFromFile(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read the first line
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// No lines found in the file
	return "", fmt.Errorf("empty file")
}

func WriteFile(filePath, content string) error {
    // Write the content to the file
    err := ioutil.WriteFile(filePath, []byte(content), 0644)
    if err != nil {
        return err
    }
    return nil
}
