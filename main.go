package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Arguments map[string]string

type itemStruct struct {
	Id    string
	Email string
	Age   int
}

type itemStructArray []itemStruct

func normalizeJson(in []byte) []byte {
	in = []byte(strings.Replace(string(in), "Age", "age", -1))
	in = []byte(strings.Replace(string(in), "Email", "email", -1))
	in = []byte(strings.Replace(string(in), "Id", "id", -1))
	return in
}

func parseFile(file *os.File) (res itemStructArray) {
	var text []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	if err := json.Unmarshal([]byte(strings.Join(text, "")), &res); err != nil {
		return nil
	}
	return
}

func Perform(args Arguments, writer io.Writer) error {

	id := args["id"]
	operation := args["operation"]
	itemString := args["item"]
	fileName := args["fileName"]
	var item itemStruct
	var content itemStructArray
	var f *os.File
	var err error

	flagNotSpecified := "-id flag has to be specified"

	if !(len(fileName) > 0) {
		return fmt.Errorf("-fileName flag has to be specified")
	} else {
		f, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
		defer f.Close()
		if err != nil {
			os.Exit(1)
		}
	}
	if !(len(operation) > 0) {
		return fmt.Errorf("-operation flag has to be specified")
	}

	switch operation {
	case "list":
		content = parseFile(f)
		data, err := json.Marshal(content)
		if err != nil {
			log.Fatal(err)
		}
		writer.Write(normalizeJson(data))

	case "add":
		if !(len(itemString) > 0) {
			return fmt.Errorf("-item flag has to be specified")
		}
		if err := json.Unmarshal([]byte(itemString), &item); err != nil {
			return nil
		}

		content = parseFile(f)

		for _, v := range content {
			if strings.TrimSpace(item.Id) == strings.TrimSpace(v.Id) {
				writer.Write([]byte("Item with id " + strings.TrimSpace(item.Id) + " already exists"))
				return nil
			}
		}

		content = append(content, item)
		data, err := json.Marshal(content)
		if err != nil {
			log.Fatal(err)
		}

		f.Write(normalizeJson(data))

	case "findById":
		if !(len(id) > 0) {
			return fmt.Errorf(flagNotSpecified)
		}
		content = parseFile(f)
		for _, v := range content {
			if v.Id == id {
				data, err := json.Marshal(v)
				if err != nil {
					log.Fatal(err)
				}
				writer.Write(normalizeJson(data))
				return nil
			}
		}
		writer.Write([]byte(""))
		return nil

	case "remove":
		if !(len(id) > 0) {
			return fmt.Errorf(flagNotSpecified)
		}
		content = parseFile(f)
		var newContent itemStructArray
		for i, v := range content {
			if v.Id == id {
				newContent = append(content[:i], content[i+1:]...)
				data, err := json.Marshal(newContent)
				if err != nil {
					log.Fatal(err)
				}
				f.Truncate(0)
				_, err = f.Seek(0, 0)
				f.Write(normalizeJson(data))
				return nil
			}
		}
		writer.Write([]byte("Item with id " + id + " not found"))
		return nil

	default:
		return fmt.Errorf("Operation " + operation + " not allowed!")

	}

	return nil

}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs() (res Arguments) {
	useId := flag.String("id", "", "object id")
	useOperation := flag.String("operation", "", "operation")
	useItem := flag.String("item", "", "item")
	useFileName := flag.String("fileName", "", "fileName")
	flag.Parse()

	res["id"] = *useId
	res["operation"] = *useOperation
	res["item"] = *useItem
	res["fileName"] = *useFileName

	return
}
