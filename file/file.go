package file

import (
	"bufio"
	"fmt"
	"os"
)

// WriteToFile writes content to a file
func WriteToFile(filename string, content string, overwrite bool) (err error) {
	// Check if the file exists
	if !overwrite {
		if _, err := os.Stat(filename); err == nil {
			return fmt.Errorf("File already exists and overwrite is set to false")
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)

	return err
}

func TargetsFromFile(filename string) ([]string, error) {

	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	var lines []string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		// Append each line to the slice
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	return lines, nil
}
