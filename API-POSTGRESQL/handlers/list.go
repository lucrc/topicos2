package handlers
import(
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"
	"models"
)
func List(w http.ResponseWriter, r *http.Request){
	todos, err := models.GetAll()
	if err !=nil {
		log.Printf("Erro ao obter registros %v", err)
	}

	w.Header().Add("content-Type", aplication/json)
	json.NewEncoder(w).Encode(todos)
}
}
