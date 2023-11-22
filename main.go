package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	modelName := getModelName()
	pluralModelName := getPluralModelName()
	modelFields := getModel(modelName)
	createModel(modelName, modelFields)
	createController(modelName, pluralModelName)
	createRouter(modelName, pluralModelName)
	//fileContent := " coucou "

	//db := database.DbConnection()

	//defer db.Close()
}
