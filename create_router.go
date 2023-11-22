package main

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"log"
	"os"
	"strings"
)

func createRouter(modelName string, pluralModelName string) {
	filePath := "./router.go"
	caser := cases.Title(language.English)
	upperModelName := caser.String(modelName)

	// Vérifier si le fichier existe déjà
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("The file %s already exist.\n", filePath)
		appendNewRoutes(filePath, modelName, pluralModelName, upperModelName)

		return
	}

	fileContent := generateFileContentRouter(modelName, pluralModelName)

	// Créer le fichier dans ./router.go
	errorFileCreate := createFile(filePath, fileContent)
	if errorFileCreate != nil {
		fmt.Println("Error during file creation:", errorFileCreate)
		return
	}

	fmt.Printf("The file has been created successfully: %s\n", filePath)
}

func generateFileContentRouter(modelName string, pluralModelName string) string {
	moduleName, errModuleName := getModuleName()
	if errModuleName != nil {
		log.Fatal(errModuleName)
	}
	content := fmt.Sprintf(`package main
	
import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"%s/controller"
)

func router() {
	router := httprouter.New()

	%v

	log.Fatal(http.ListenAndServe(":8080", router))
}`, moduleName, createRoutes(modelName, pluralModelName))

	return content
}

func createRoutes(modelName string, pluralModelName string) string {
	lowerModelName := strings.ToLower(modelName)
	lowerPluralModelName := strings.ToLower(pluralModelName)
	caser := cases.Title(language.English)
	upperPluralModelName := caser.String(pluralModelName)
	upperModelName := caser.String(modelName)

	content := fmt.Sprintf(`	/********* %s routes **********/
	%sControllerInstance := controller.%sController{}
	router.GET("/%s/", %sControllerInstance.GetAll%s)
	router.GET("/%s/:id", %sControllerInstance.Get%sById)
	router.POST("/%s/create", %sControllerInstance.Create%s)
	router.PUT("/%s/:id", %sControllerInstance.Update%s)
	router.DELETE("/%s/:id", %sControllerInstance.Delete%s)`,
		upperModelName, // trié par ligne
		lowerModelName, upperModelName,
		lowerPluralModelName, lowerModelName, upperPluralModelName,
		lowerModelName, lowerModelName, upperModelName,
		lowerModelName, lowerModelName, upperModelName,
		lowerModelName, lowerModelName, upperModelName,
		lowerModelName, lowerModelName, upperModelName,
	)
	return content
}

func appendNewRoutes(filePath string, modelName string, pluralModelName string, upperModelName string) {
	addContent := fmt.Sprintf("\n%sControllerInstance := controller.%sController{}\n", modelName, upperModelName)
	// Ouvrir le fichier en mode lecture-écriture
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Lire tout le contenu du fichier
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error writting file:", err)
		return
	}

	// Trouver la position où les dernières routes ont été écrites
	lastRoutesIndex := strings.LastIndex(string(content), "/**")
	if lastRoutesIndex == -1 {
		fmt.Println("Impossible to find the position of the last roads.")
		return
	}

	// Déplacer le curseur à la position de la dernière occurrence
	currentPositionPointer, seekErr := file.Seek(int64(lastRoutesIndex), 0)
	if seekErr != nil {
		fmt.Println("Cursor movement error:", seekErr)
		return
	}

	// Ajouter de nouvelles routes CRUD au contenu existant
	newRoutes := createRoutes(modelName, pluralModelName)

	contentAfterLine50, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Écrire le contenu modifié dans le fichier avant les premieres routes
	_, err = file.WriteAt([]byte(addContent+newRoutes+"\n"+string(contentAfterLine50)), currentPositionPointer)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Printf("New routes successfully added to: %s\n", filePath)
}
