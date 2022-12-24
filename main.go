package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Server struct {
		DbHost string `yaml:"host"`
		DbPort string `yaml:"port"`
	}

	User struct {
		DbUser     string `yaml:"user"`
		DbPassword string `yaml:"password"`
	} `yaml:"login"`
}

type LoginValues struct {
	Username string
	Password string
}

const htmlLogin = `<html>
<body>
    <form action="/login" method="POST">
        <label for="username">Username:</label><br>
        <input type="text" id="username" name="username"><br>
        <label for="password">Password:</label><br>
        <input type="password" id="password" name="password"><br><br>
        <input type="submit" value="Submit">
    </form> 
</body>
</html>`

func readConfig(cfg *DBConfig) error {
	conf, err := os.ReadFile("conf.yml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(conf), &cfg)
	if err != nil {
		return err
	}
	if cfg.Server.DbHost == "" || cfg.Server.DbPort == "" || cfg.User.DbUser == "" || cfg.User.DbPassword == "" {
		err := errors.New("empty fields")
		return err
	}
	fmt.Println(cfg)
	return err
}

func connectDB(cfg *DBConfig) *sql.DB {
	connection := cfg.User.DbUser + ":" + cfg.User.DbPassword + "@tcp(" + cfg.Server.DbHost + ":" + cfg.Server.DbPort + ")/test"
	db, err := sql.Open("mysql", connection)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func createTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (username VARCHAR(25) PRIMARY KEY NOT NULL, password VARCHAR(25) NOT NULL)")
	return err
}

func (lv *LoginValues) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Get the username and password from the form
		lv.Username = r.FormValue("username")
		lv.Password = r.FormValue("password")

		checker, err := checkLogin(lv.Username, lv.Password)
		if err != nil {
			fmt.Println(err)
		}
		// Check the username and password against a database or other method
		if checker {
			http.Redirect(w, r, "/mfa", http.StatusSeeOther)
		} else {
			fmt.Fprint(w, "Incorrect username or password. Please try again.")
		}
	} else {
		fmt.Fprint(w, htmlLogin)
	}
}

func checkLogin(username string, password string) (bool, error) {
	var correct bool
	correct = true
	return correct, nil

}

func newfunc(lv *LoginValues) {
	fmt.Println(lv.Username)
}

func main() {
	var config DBConfig
	var loginValues LoginValues
	readConfig(&config)
	db := connectDB(&config)
	err := createTable(db)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/login", loginValues.login)
	http.ListenAndServe(":8080", nil)
}
