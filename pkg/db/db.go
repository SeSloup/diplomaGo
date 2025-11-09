package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

//Файл базы данных scheduler.db должен содержать таблицу scheduler с такими колонками:
//id — автоинкрементный идентификатор;
//date — дата задачи, которая будет храниться в формате YYYYMMDD или в Go-представлении 20060102;
//title — заголовок задачи;
//comment — комментарий к задаче;
//repeat — строковое поле не более 128 символов, которое будет содержать правила повторений для задачи. Формат правил будет описан в следующем шаге.

const schema = `CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(255) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);

	CREATE INDEX idx_scheduler_date ON scheduler(date);`

var DB *sql.DB

func Init(dbFile string) error {

	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	DB, err = sql.Open("sqlite", dbFile)

	if install {
		if _, err := DB.Exec(schema); err != nil {
			log.Fatal("failed to create schema: ", err.Error())
		}

	}

	return err
}
