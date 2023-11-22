package models

// Esquema de usuário da tabela users
type User struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    Location string `json:"location"`
    Age      int64  `json:"age"`
}