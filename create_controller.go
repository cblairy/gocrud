package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func createController(modelName string, pluralModelName string) {

	fileContentController := generateFileContentController(modelName, pluralModelName)

	err := os.MkdirAll("./controller/", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory ./controller:", err)
		return
	}

	filePath := fmt.Sprintf("./controller/%s_controller.go", strings.ToLower(modelName))
	err = createFile(filePath, fileContentController)
	if err != nil {
		fmt.Println("Error during file creation:", err)
		return
	}
	fmt.Printf("The file has been successfully created: %s\n", filePath)
}

func generateFileContentController(modelName string, pluralModelName string) string {
	caser := cases.Title(language.English)
	upperModelName := caser.String(modelName)
	upperPluralModelName := caser.String(pluralModelName)
	moduleName, errModuleName := getModuleName()
	if errModuleName != nil {
		log.Fatal(errModuleName)
	}
	content := fmt.Sprintf(`package controller

	import (
		"database/sql"
		"encoding/json"
		"fmt"
		"log"
		"net/http"
		"%s/model"
		"strings"
		"reflect"

		"github.com/julienschmidt/httprouter"
	)
	
	type %sController struct {
		DB *sql.DB 
	}
	
	/** GETALL */
	func (c *%sController) GetAll%s(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		rows, err := c.DB.Query("SELECT * FROM %s")
		if err != nil {
			log.Println("Erreur lors de la récupération de tous les éléments:", err)
			http.Error(w, "Erreur lors de la récupération de tous les éléments", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
	
		var results []model.%sModel
	
		for rows.Next() {
			var result model.%sModel
			err := rows.Scan(&result)
			if err != nil {
				log.Println("Erreur lors du scan des résultats:", err)
				http.Error(w, "Erreur lors de la récupération des éléments", http.StatusInternalServerError)
				return
			}
			results = append(results, result)
		}
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
	
	/** GETBYID */
	func (c *%sController) Get%sById(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "Paramètre id manquant", http.StatusBadRequest)
			return
		}
	
		row := c.DB.QueryRow("SELECT * FROM %s WHERE id=?", id)
	
		var result model.%sModel
		err := row.Scan(&result)
		if err != nil {
			log.Println("Erreur lors de la récupération de l'élément par ID:", err)
			http.Error(w, "Erreur lors de la récupération de l'élément par ID", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
	
	/** CREATE */
	func (c *%sController) Create%s(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var model model.%sModel
		err := json.NewDecoder(r.Body).Decode(&model)
		if err != nil {
			log.Println("Erreur lors de la lecture du corps de la demande:", err)
			http.Error(w, "Erreur lors de la lecture du corps de la demande", http.StatusBadRequest)
			return
		}
	
		var columns []string
		var values []interface{}

		modelType := reflect.TypeOf(model)
		modelValue := reflect.ValueOf(model)

		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			columns = append(columns, field.Name)
			values = append(values, modelValue.Field(i).Interface())
		}

		query := fmt.Sprintf("INSERT INTO %s (%%s) VALUES (%%s)", strings.Join(columns, ", "), strings.Join(make([]string, len(columns)), ", "))

		_, err = c.DB.Exec(query, values...)
		if err != nil {
			log.Println("Erreur lors de l'insertion d'un nouvel élément:", err)
			http.Error(w, "Erreur lors de la création d'un nouvel élément", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
	
	/** UPDATE */
	func (c *%sController) Update%s(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "Paramètre id manquant", http.StatusBadRequest)
			return
		}
	
		err := r.ParseForm()
		if err != nil {
			log.Println("Erreur lors de l'analyse du formulaire:", err)
			http.Error(w, "Erreur lors de l'analyse du formulaire", http.StatusInternalServerError)
			return
		}
	
		var setValues []string
		for key, values := range r.Form {
			if len(values) > 0 {
				setValues = append(setValues, fmt.Sprintf("%%s = ?", key))
			}
		}
	
		query := fmt.Sprintf("UPDATE %s SET %%s WHERE id = ?", strings.Join(setValues, ", "))
		values := append(r.Form["value"], id)

		_, err = c.DB.Exec(query, values)
		if err != nil {
			log.Println("Erreur lors de la mise à jour de l'élément par ID:", err)
			http.Error(w, "Erreur lors de la mise à jour de l'élément par ID", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
	
	/** DELETE */
	func (c *%sController) Delete%s(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "Paramètre id manquant", http.StatusBadRequest)
			return
		}
	
		_, err := c.DB.Exec("DELETE FROM %s WHERE id=?", id)
		if err != nil {
			log.Println("Erreur lors de la suppression de l'élément par ID:", err)
			http.Error(w, "Erreur lors de la suppression de l'élément par ID", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

`, moduleName, // tri par ligne, %%s ne compte pas
		upperModelName,
		upperModelName, upperPluralModelName,
		strings.ToLower(modelName),
		upperModelName,
		upperModelName,
		upperModelName, upperModelName,
		modelName,
		upperModelName,
		upperModelName, upperModelName,
		upperModelName,
		upperModelName,
		upperModelName,
		upperModelName,
		strings.ToLower(modelName),
		upperModelName, upperModelName,
		upperModelName)

	return content
}
