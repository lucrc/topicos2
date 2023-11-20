package handlers

import(
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
	"models"
)

func Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		log.Printf("Erro ao fazer parde do id: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalSeverError)
		return
	}

	var todo models.Todo
	
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil{
		log.Printf("Erro ao fazer decode do json: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalSeverError)
		return
	}

	rows, err := models.Update(int64(id), todo)
	if err != nil {
		log.Printf("Erro ao atualizar o registro: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalSeverError)
		return
	}

	if rows > 1 {
		log.Printf("ErroR: foram atulizados %d registros", rows)
	}
	resp := map[string]any{
		"Error": false,
		"Message": "Dados atualizados com sucesso"
	}
	w.Header().Add("content-Type", aplication/json)
	json.NewEncoder(w).Encode(resp)
}