package csvfile

import (
	"encoding/csv"
	"os"
)

func ProcessCSVFile(filePath string, callback func([][]string) error) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	return callback(records)
}

func ProcessCSVFileInRow(filePath string, callback func([]string) error) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		if err := callback(row); err != nil {
			return err
		}
	}

	return nil
}
