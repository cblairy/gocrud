package main

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"strings"
)

func createModel(modelName string, modelFields map[string]string) {
	fileContent := generateFileContent(modelName, modelFields)

	// Créer le répertoire
	err := os.MkdirAll("./model/", os.ModePerm)
	if err != nil {
		fmt.Println("Directory creation error:", err)
		return
	}

	// Créer le fichier dans /model/<nom du modèle>
	filePath := fmt.Sprintf("./model/%s_model.go", strings.ToLower(modelName))
	err = createFile(filePath, fileContent)
	if err != nil {
		fmt.Println("Error during file creation:", err)
		return
	}
	fmt.Printf("The file has been created successfully: %s\n", filePath)
}

// generateFileContent génère le contenu du fichier de modele
func generateFileContent(modelName string, fields map[string]string) string {
	var content string
	caser := cases.Title(language.English)

	content += "package model\n\n"

	dateFieldExists := false
	for _, field := range fields {
		if field == "date" {
			dateFieldExists = true
			break
		}
	}

	if dateFieldExists {
		content += `import "time"`
	}

	content += fmt.Sprintf("\n\ntype %sModel struct {\n", caser.String(modelName))

	// Ajouter les champs à l'interface de type
	for fieldName, fieldType := range fields {
		content += fmt.Sprintf("\t%s %s\n", caser.String(fieldName), goType(fieldType))
	}

	content += "}\n"

	return content
}

func goType(fieldType string) string {
	switch fieldType {
	case "float":
		return "float32"
	case "date":
		return "*time.Time"
	default:
		return fieldType
	}
}
