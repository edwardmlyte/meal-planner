package controller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/chvck/meal-planner/model/recipe"
	"github.com/chvck/meal-planner/model/user"
	"github.com/chvck/meal-planner/store"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// RecipeIndex is the HTTP handler for the recipe index endpoint
func RecipeIndex(w http.ResponseWriter, r *http.Request) {
	db := store.Database()
	u := context.Get(r, "user").(user.User)
	recipes, err := recipe.All(db, u.ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Could not retrieve recipes", http.StatusNotFound)
		return
	}

	JSONResponse(*recipes, w)
}

// RecipeByID is the HTTP handler for fetching a single recipe
func RecipeByID(w http.ResponseWriter, r *http.Request) {
	db := store.Database()
	vars := mux.Vars(r)
	u := context.Get(r, "user").(user.User)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	recipes, err := recipe.One(db, id, u.ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Could not retrieve recipe", http.StatusNotFound)
		return
	}

	JSONResponse(recipes, w)
}

// RecipeCreate is the HTTP handler for creating a recipe
func RecipeCreate(w http.ResponseWriter, r *http.Request) {
	db := store.Database()
	var re recipe.Recipe
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		log.Println(err.Error())
		http.Error(w, "Invalid recipe", http.StatusBadRequest)
		return
	} else {
		if err := json.Unmarshal(body, &re); err != nil {
			log.Println(err.Error())
			http.Error(w, "Invalid recipe", http.StatusBadRequest)
			return
		}
	}

	u := context.Get(r, "user").(user.User)
	err := recipe.Create(db, re, u.ID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Could not create recipe", http.StatusInternalServerError)
		return
	}

}