package resource

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/tesarwijaya/ouroboros/internal/config"
)

func NewSQLConnection(c *config.Config) (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.SqlDBHost, c.SqlDBPort, c.SqlDBUsername, c.SqlDBPassword, c.SqlDBName)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	fmt.Println("db ping..")
	if err != nil {
		return nil, err
	}
	fmt.Println("db ping success!")

	return db, err
}
