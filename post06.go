package post06

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type UserData struct {
	ID          int
	Username    string
	Name        string
	Suname      string
	Description string
}

var (
	Hostname = ""
	Port     = "2345"
	Username = ""
	Password = ""
	Database = ""
)

func openConnection() (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, Username, Password, Database)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func exists(username string) int {
	username = strings.ToLower(username)
	db, err := openConnection()
	if err != nil {
		return -1
	}
	defer db.Close()
	UserID := -1
	statement := fmt.Sprintf(`SELECT "id" FROM "users" where username = '%s'`, username)
	rows, err := db.Query(statement)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		UserID = id

	}
	defer rows.Close()
	return UserID
}

func AddUser(d UserData) int {
	d.Username = strings.ToLower(d.Username)
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()
	userID := exists(d.Username)
	if userID != -1 {
		fmt.Println("User already exists:", d.Username)
		return -1
	}
	insertStatement := `insert into "users"(username)values ($1)`
	_, err = db.Exec(insertStatement, d.Username)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	userID = exists(d.Username)
	if userID == -1 {
		return userID
	}

	insertStatement = `insert into "userdata"("userid","name", "suname","discription") values ($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, userID, d.Name, d.Suname, d.Description)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}
	return userID
}

func deleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	db.Close()
	statement := fmt.Sprintf(`SELECT "username" FROM "users" where id = %d`, id)
	rows, err := db.Query(statement)

	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}

	}
	defer rows.Close()
	if exists(username) != id {
		return fmt.Errorf("User with ID %d does not exist", id)
	}
	deleteStatement := `delete from "userdata" where userid=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}
	return nil
}
func ListUsers() ([]UserData, error) {
	Data := []UserData{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT "id","username","name","surname","description"
		FROM "users","userdata"
		WHERE users.id = userdata.userid`)
	if err != nil {
		return Data, err
	}

	for rows.Next() {
		var id int
		var username string
		var name string
		var surname string
		var description string
		err = rows.Scan(&id, &username, &name, &surname, &description)
		temp := UserData{ID: id, Username: username, Name: name, Suname: surname, Description: description}
		Data = append(Data, temp)
		if err != nil {
			return Data, err
		}
	}
	defer rows.Close()
	return Data, nil
}
func UpdateUser(d UserData) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID == -1 {
		return errors.New("User does not exist")
	}
	d.ID = userID
	updateStatement := `update "userdata" set "name"=$1, "surname"=$2, "description"=$3 where "userid"=$4`
	_, err = db.Exec(updateStatement, d.Name, d.Suname, d.Description, d.ID)
	if err != nil {
		return err
	}

	return nil
}
