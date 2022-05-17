package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/peterbourgon/ff/v3"
	"github.com/pkg/errors"
	"os"
)

type Student struct {
	Firstname      string
	Lastname       string
	GitTutorialUrl string
}

func loadCSV(path string) (error, [][]string) {
	f, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "cannot open file %s", path), nil
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
	if err != nil {
		return errors.Wrapf(err, "could not read %s as a csv", path), nil
	}
	return nil, records
}

func castCSVEntries(entries [][]string) (error, []Student) {
	students := []Student{}
	for _, entry := range entries[1:] {
		firstName := entry[0]
		if firstName == "" {
			return errors.Errorf("invalid entry[0] as firstname %+v", entry), nil
		}
		lastName := entry[1]
		if lastName == "" {
			return errors.Errorf("invalid entry[1] as lastname %+v", entry), nil
		}
		gitTutorialUrl := entry[3]
		students = append(students, Student{
			Firstname:      firstName,
			Lastname:       lastName,
			GitTutorialUrl: gitTutorialUrl,
		})
	}
	return nil, students
}

func main() {
	fs := flag.NewFlagSet("osef", flag.ExitOnError)
	var (
		studentsCSV = fs.String("studentsCSV", "", "")
	)
	ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())

	if *studentsCSV == "" {
		panic("need a csv")
	}
	err, data := loadCSV(*studentsCSV)
	if err != nil {
		fmt.Printf("Could not load CSV: %+v\n", err)
		panic("")
	}
	err, students := castCSVEntries(data)
	if err != nil {
		fmt.Printf("Could not parse CSV: %+v\n", err)
		panic("")
	}

	csvFile, err := os.Create("out_tmp.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	w := csv.NewWriter(csvFile)
	defer w.Flush()

	records := [][]string{}
	for _, s := range students {
		out, err := gradeTutorialFromUrl(s.GitTutorialUrl, s.Firstname, s.Lastname)
		if err != nil {
			fmt.Printf("Could not grade tuto: %+v\n", err)
			panic("")
		}
		fmt.Printf("%+v - %+v\n", s, out)
		records = append(records, []string{s.Firstname, fmt.Sprint(out.Grade), out.Comment})
	}

	w.WriteAll(records)

}
