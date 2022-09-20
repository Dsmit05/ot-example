package config

import "fmt"

type Config struct {
	CollectorURL string
	Database
}

type Database struct {
	Host     string
	Port     int
	Table    string
	User     string
	Password string
}

func (d Database) GetConnStringDB() string {
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Table,
	)
}
