package lib

import (
	"bufio"
	"encoding/csv"
	"errors"
	"os"
)

func LinesInFile(fileName string) ([]string, error) {
	f, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	result := []string{}

	for scanner.Scan() {
		line := scanner.Text()

		result = append(result, line)
	}

	return result, err
}

func ReadCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	defer f.Close()

	if err != nil {
		return nil, errors.New("Unable to read input file " + filePath + "\n" + err.Error())
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()

	if err != nil {
		return nil, errors.New("Unable to parse file as CSV for " + filePath + "\n" + err.Error())
	}

	return records, err
}

func WriteCsvFile(filePath string, records [][]string) error {
	f, err := os.Create(filePath)
	defer f.Close()

	if err != nil {
		return errors.New("Unable to open output file " + filePath + "\n" + err.Error())
	}

	csvWriter := csv.NewWriter(f)
	defer csvWriter.Flush()

	for _, record := range records {
		err = csvWriter.Write(record)

		if err != nil {
			return errors.New("Unable to write to output file " + filePath + "\n" + err.Error())
		}
	}

	return err
}
