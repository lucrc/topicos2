package middleware

import (
    "database/sql"
    "encoding/json" // package para codificar e decodificar o json em struct e vice-versa
    "fmt"
    "go-postgres/models" // models package onde o esquema do usuário é definido
    "log"
    "net/http" // usado para acessar o objeto de solicitação e resposta da API
    "os"       // usado para ler a variável de ambiente
    "strconv"  // package usado para converter string em tipo int

    "github.com/gorilla/mux" // usado para obter os parâmetros da rota

    "github.com/joho/godotenv" // package usado para ler o arquivo .env
    _ "github.com/lib/pq"      // driver postgres golang
)

// formato de respostas
type response struct {
    ID      int64  `json:"id,omitempty"`
    Message string `json:"message,omitempty"`
}

// criar conexão com postgres db
func createConnection() *sql.DB {
    // carregar arquivo .env
    err := godotenv.Load(".env")

    if err != nil {
        log.Fatalf("Erro ao carregar o arquivo .env")
    }

    // Abrir a conexão
    db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

    if err != nil {
        panic(err)
    }

    // verificando a conexão
    err = db.Ping()

    if err != nil {
        panic(err)
    }

    fmt.Println("Conectado com sucesso!")
    // retorno da conexão
    return db
}

// CreateUser cria um usuário no banco de dados postgres
func CreateUser(w http.ResponseWriter, r *http.Request) {
    // define o cabeçalho para o tipo de conteúdo x-www-form-urlencoded
    //Permitir que todas as origens resolvam o problema do cors
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // crie um usuário vazio do tipo models.User
    var user models.User

    // decodificar a solicitação json para o usuário
    err := json.NewDecoder(r.Body).Decode(&user)

    if err != nil {
        log.Fatalf("Não foi possível decodificar o corpo da solicitação.  %v", err)
    }

    // chame a função de inserir usuário e passe o usuário
    insertID := insertUser(user)

    // formatar um objeto de resposta
    res := response{
        ID:      insertID,
        Message: "Usuário criado com sucesso",
    }

    // envio da resposta
    json.NewEncoder(w).Encode(res)
}

// GetUser retornará um único usuário por seu id
func GetUser(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // obtenha o ID do usuário dos parâmetros da solicitação, a chave é "id"
    params := mux.Vars(r)

    // converter o tipo de id de string para int
    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }

    // chama a função getUser com ID de usuário para recuperar um único usuário
    user, err := getUser(int64(id))

    if err != nil {
        log.Fatalf("Não foi possível obter o usuário. %v", err)
    }

    // envie a resposta
    json.NewEncoder(w).Encode(user)
}

// GetAllUser retornará todos os usuários
func GetAllUser(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // obter todos os usuários no banco de dados
    users, err := getAllUsers()

    if err != nil {
        log.Fatalf("Não foi possível obter todos os usuários. %v", err)
    }

    // enviar todos os usuários como resposta
    json.NewEncoder(w).Encode(users)
}

// UpdateUser atualiza detalhes do usuário no banco de dados postgres
func UpdateUser(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // obtenha o ID do usuário dos parâmetros de solicitação, a chave é "id"
    params := mux.Vars(r)

    // converter o tipo de id de string para int
    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Não foi possível converter a string em int.  %v", err)
    }

    // crie um usuário vazio do tipo models.User
    var user models.User

    // decodificar a solicitação json para o usuário
    err = json.NewDecoder(r.Body).Decode(&user)

    if err != nil {
        log.Fatalf("Não foi possível decodificar o corpo da solicitação.  %v", err)
    }

    // chamar update user para atualizar o usuário
    updatedRows := updateUser(int64(id), user)

    // formatar a string da mensagem
    msg := fmt.Sprintf("Usuário atualizado com sucesso. Total de linhas/registros afetados %v", updatedRows)

    // formatar a mensagem de resposta
    res := response{
        ID:      int64(id),
        Message: msg,
    }

    // envie a resposta
    json.NewEncoder(w).Encode(res)
}

// DeleteUser exclui detalhes do usuário no banco de dados postgres
func DeleteUser(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // obtenha o ID do usuário dos parâmetros de solicitação, a chave é "id"
    params := mux.Vars(r)

    // converta o id em string para int
    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Não foi possível converter a string em int.  %v", err)
    }

    // chame o deleteUser, converta o int em int64
    deletedRows := deleteUser(int64(id))

    // formatar a string da mensagem
    msg := fmt.Sprintf("Usuário atualizado com sucesso. Total de linhas/registros afetados %v", deletedRows)

    // formatar a mensagem de resposta
    res := response{
        ID:      int64(id),
        Message: msg,
    }

    // send the response
    json.NewEncoder(w).Encode(res)
}

//------------------------- funções de manipulador----------------
// insira um usuário no banco de dados
func insertUser(user models.User) int64 {

    // crie a conexão postgres db
    db := createConnection()

    // feche a conexão db
    defer db.Close()

    // cria a consulta sql de inserção
    // retornar o userid retornará o id do usuário inserido
    sqlStatement := `INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid`

    // o id inserido será armazenado neste id
    var id int64

    //executa a instrução SQL
    // A função de scaneamento salvará o ID de inserção no id
    err := db.QueryRow(sqlStatement, user.Name, user.Location, user.Age).Scan(&id)

    if err != nil {
        log.Fatalf("Não foi possível executar a consulta. %v", err)
    }

    fmt.Printf("Inseriu um único registro %v", id)

    // retornar o id inserido
    return id
}

// obtenha um usuário do banco de dados por seu ID de usuário
func getUser(id int64) (models.User, error) {
    // crie a conexão postgres db
    db := createConnection()

    // feche a conexão db
    defer db.Close()

    // crie um usuário de models.User type
    var user models.User

    // create the select sql query
    sqlStatement := `SELECT * FROM users WHERE userid=$1`

    // crie a consulta sql selecionada
    row := db.QueryRow(sqlStatement, id)

    // desempacotar o objeto de linha para o usuário
    err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

    switch err {
    case sql.ErrNoRows:
        fmt.Println("Nenhuma linha foi retornada!")
        return user, nil
    case nil:
        return user, nil
    default:
        log.Fatalf("Não foi possível verificar a linha. %v", err)
    }

    // retornar usuário vazio em caso de erro
    return user, err
}

// obtenha todos os usuários do banco de dados.
func getAllUsers() ([]models.User, error) {
    // crie a conexão postgres db
    db := createConnection()

    // feche a conexão db
    defer db.Close()

    var users []models.User

    // crie a consulta sql selecionada
    sqlStatement := `SELECT * FROM users`

    // execute a instrução SQL
    rows, err := db.Query(sqlStatement)

    if err != nil {
        log.Fatalf("Não foi possível executar a consulta. %v", err)
    }

    // feche a declaração
    defer rows.Close()

    // iterar sobre as linhas
    for rows.Next() {
        var user models.User

        // desempacotar o objeto de linha para o usuário
        err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

        if err != nil {
            log.Fatalf("Não foi possível verificar a linha. %v", err)
        }

        // anexe o usuário na fatia de usuários
        users = append(users, user)

    }

    // return empty user on error
    return users, err
}

// atualizar usuário no banco de dados
func updateUser(id int64, user models.User) int64 {

    // crie a conexão postgres db
    db := createConnection()

    // feche a conexão db
    defer db.Close()

    // crie a consulta sql de atualização
    sqlStatement := `UPDATE users SET name=$2, location=$3, age=$4 WHERE userid=$1`

    // execute a instrução SQL
    res, err := db.Exec(sqlStatement, id, user.Name, user.Location, user.Age)

    if err != nil {
        log.Fatalf("Não foi possível executar a consulta. %v", err)
    }

    // verifique quantas linhas afetadas
    rowsAffected, err := res.RowsAffected()

    if err != nil {
        log.Fatalf("Erro ao verificar as linhas afetadas. %v", err)
    }

    fmt.Printf("Total de linhas/registros afetados %v", rowsAffected)

    return rowsAffected
}

// excluir usuário no banco de dados
func deleteUser(id int64) int64 {

    // crie a conexão postgres db
    db := createConnection()

    // feche a conexão db
    defer db.Close()

    // crie a consulta sql de exclusão
    sqlStatement := `DELETE FROM users WHERE userid=$1`

    // execute a instrução SQL
    res, err := db.Exec(sqlStatement, id)

    if err != nil {
        log.Fatalf("Não foi possível executar a consulta. %v", err)
    }

    // verifique quantas linhas afetadas
    rowsAffected, err := res.RowsAffected()

    if err != nil {
        log.Fatalf("Erro ao verificar as linhas afetadas. %v", err)
    }

    fmt.Printf("Total de linhas/registro afetadod %v", rowsAffected)

    return rowsAffected
}