package models

import(
	"db"
)

func GetAll() (todos []Todo, err errro) {
	conn, err := OpenConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	rows := conn.Query(`SELECT * FROM todos` )

	if err != nil {
		return
	}

	for rows.Next(){
		var todo Todo
	
		err = rows.Scan(&todo.ID, & todo.Title, &todo.Description, todo.Done)
		if err != nill {
			continue
		}
		todos = append(todos, todo)
	}	
	return
}