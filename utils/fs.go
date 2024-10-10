package utils

import "os"

func WriteToFile(filename, content string) error {
	file, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LogError("file-writer", "could not close file: "+err.Error())
		}
	}(file)

	_, err = file.WriteString(content)

	if err != nil {
		return err
	}

	return nil
}
