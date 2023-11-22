package controller

	import (
		"database/sql"
		"encoding/json"
		"fmt"
		"log"
		"net/http"
		"crud/model"
		"strings"
		"reflect"

		"github.com/julienschmidt/httprouter"
	)
	
	type UserController struct {
		DB *sql.DB 
	}
	
	/** GETALL */
	func (c *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		rows, err := c.DB.Query("SELECT * FROM user")
		if err != nil {
			log.Println("Erreur lors de la récupération de tous les éléments:", err)
			http.Error(w, "Erreur lors de la récupération de tous les éléments", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
	
		var results []model.UserModel
	
		for rows.Next() {
			var result model.UserModel
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
	func (c *UserController) GetUserById(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "Paramètre id manquant", http.StatusBadRequest)
			return
		}
	
		row := c.DB.QueryRow("SELECT * FROM user WHERE id=?", id)
	
		var result model.UserModel
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
	func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var model model.UserModel
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

		query := fmt.Sprintf("INSERT INTO User (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(make([]string, len(columns)), ", "))

		_, err = c.DB.Exec(query, values...)
		if err != nil {
			log.Println("Erreur lors de l'insertion d'un nouvel élément:", err)
			http.Error(w, "Erreur lors de la création d'un nouvel élément", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
	
	/** UPDATE */
	func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
				setValues = append(setValues, fmt.Sprintf("%s = ?", key))
			}
		}
	
		query := fmt.Sprintf("UPDATE user SET %s WHERE id = ?", strings.Join(setValues, ", "))
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
	func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "Paramètre id manquant", http.StatusBadRequest)
			return
		}
	
		_, err := c.DB.Exec("DELETE FROM User WHERE id=?", id)
		if err != nil {
			log.Println("Erreur lors de la suppression de l'élément par ID:", err)
			http.Error(w, "Erreur lors de la suppression de l'élément par ID", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

