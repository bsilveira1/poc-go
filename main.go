package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	// _ "github.com/lib/pq" // postgres driver
	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/rs/cors"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *sql.DB

func main() {
	var err error
	// Configuração para PostgreSQL
	// db, err = sql.Open("postgres", "host=postgres user=testuser password=testpassword dbname=testdb sslmode=disable")
	// if err != nil {
	// 	log.Fatalf("Erro ao conectar ao banco: %v", err)
	// }

	// Configuração para MySQL (comentada)
	db, err = sql.Open("mysql", "testuser:testpassword@tcp(mysql:3306)/testdb")
	if err != nil {
	    log.Fatalf("Erro ao conectar ao banco: %v", err)
	}

    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowCredentials: true,
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
    })

  	handler := c.Handler(http.DefaultServeMux)

	http.HandleFunc("/users", handleUsers)
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recebida requisição %s em /users", r.Method)

	switch r.Method {
	case "GET":
		getUsers(w, r)
	case "POST":
		createUser(w, r)
	case "PUT":
		updateUser(w, r)
	case "DELETE":
		deleteUser(w, r)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		log.Printf("Método %s não permitido", r.Method)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Buscando usuários no banco de dados...")

	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, "Erro ao buscar usuários", http.StatusInternalServerError)
		log.Printf("Erro ao buscar usuários: %v", err)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			http.Error(w, "Erro ao processar dados", http.StatusInternalServerError)
			log.Printf("Erro ao processar dados dos usuários: %v", err)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)

	log.Printf("Retornando %d usuários", len(users))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Criando um novo usuário...")

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Erro ao processar entrada", http.StatusBadRequest)
		log.Printf("Erro ao decodificar corpo da requisição: %v", err)
		return
	}

	_, err := db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
	if err != nil {
		http.Error(w, "Erro ao criar usuário", http.StatusInternalServerError)
		log.Printf("Erro ao inserir usuário no banco de dados: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("Usuário criado com sucesso: %s", user.Name)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Atualizando usuário...")

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Erro ao processar entrada", http.StatusBadRequest)
		log.Printf("Erro ao decodificar corpo da requisição: %v", err)
		return
	}

	_, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", user.Name, user.Email, user.ID)
	if err != nil {
		http.Error(w, "Erro ao atualizar usuário", http.StatusInternalServerError)
		log.Printf("Erro ao atualizar usuário no banco de dados: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Usuário atualizado com sucesso: %s", user.Name)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Deletando usuário...")

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		log.Printf("ID inválido fornecido: %v", err)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Erro ao deletar usuário", http.StatusInternalServerError)
		log.Printf("Erro ao deletar usuário com ID %d: %v", id, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Usuário com ID %d deletado com sucesso", id)
}
