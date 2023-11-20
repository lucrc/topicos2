package handlers

import(
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
	"models"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Erro ao fazer parse do id: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalSeverError)
		return
	}

	
	rows, err := models.Delete(int64(id))
	if err != nil {
		log.Printf("Erro ao removere o registro: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalSeverError)
		return
	}

	if rows > 1 {
		log.Printf("ErroR: foram removidos %d registros", rows)
	}
	resp := map[string]any{
		"Error": false,
		"Message": "Registro removido com sucesso"
	}
	w.Header().Add("content-Type", aplication/json)
	json.NewEncoder(w).Encode(resp)
}