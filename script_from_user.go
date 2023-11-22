package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var fieldTypes = []string{"int", "float", "bool", "string", "date"}

func getModel(modelName string) map[string]string {
	modelFields := getFieldsFromUser()

	// Afficher les données saisies
	fmt.Println("Data entered for the model", modelName, ":")
	for fieldName, fieldType := range modelFields {
		fmt.Printf("%s: %s\n", fieldName, fieldType)
	}
	return modelFields
}

func getModelName() string {
	fmt.Print("Enter model name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	modelName := scanner.Text()
	return modelName
}

func getPluralModelName() string {
	fmt.Print("Enter the model name a second time (plural): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	modelName := scanner.Text()
	return modelName
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func getFieldsFromUser() map[string]string {
	fields := make(map[string]string)

	for {
		var fieldName = ""

		for len(fieldName) == 0 {
			fieldName = getUserInput("Enter the field name (or type 'end' to end): ")
		}
		if strings.ToLower(fieldName) == "end" {
			break
		}

		// Afficher une liste de types prédéfinis pour que l'utilisateur choisisse
		fieldType := getUserInput("Enter the field type: [int, float, bool, string, date] ")

		if !contains(fieldTypes, fieldType) {
			fmt.Println("Invalid type. Please choose from the predefined types.")
			continue
		}

		fields[fieldName] = fieldType
	}

	return fields
}

func contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}
