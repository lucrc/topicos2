package handlers

import(
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
	"models"
)

func Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Erro ao fazer parde do id: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalSeverError)
		return
	}
	
	todo, err := models.Get(int64(id))
	if err != nil {
		log.Printf("Erro ao atualizar o registro: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalSeverError)
		return
	}

	w.Header().Add("content-Type", aplication/json)
	json.NewEncoder(w).Encode(todo)
}