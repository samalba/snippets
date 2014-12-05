package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Organization struct {
	Id          int64
	Name        string `sql:"type:varchar(100);"`
	Description string `sql:"type:varchar(255);"`
	Url         string `sql:"type:varchar(100);"`
	Email       string `sql:"type:varchar(100);"`
	Teams       []Team
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type Team struct {
	Id              int64
	Name            string `sql:"type:varchar(100);"`
	Description     string `sql:"type:varchar(255);"`
	Url             string `sql:"type:varchar(100);"`
	Email           string `sql:"type:varchar(100);"`
	OrganizationId  int64  // ForeignKey to Organization
	TeamMemberships []TeamMembership
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}

// Many-To-Many relationship between Teams and Users
type TeamMembership struct {
	Id     int64
	TeamId int64 // ForeignKey to a Team
	UserId int64 // ForeignKey to a User
}

type User struct {
	Id              int64
	Name            string `sql:"type:varchar(100);"`
	Company         string `sql:"type:varchar(100);"`
	Email           string `sql:"type:varchar(100);"`
	TeamMemberships []TeamMembership
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       time.Time
}

type Snippet struct {
	Id        int64
	Week      int64 // Week number
	Year      int64
	UserId    int64 // ForeignKey to a User
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

var gormDb gorm.DB

func init() {
	gormDb, err := gorm.Open("sqlite3", "snippets.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot init the db: %s\n", err)
	}
	gormDb.DB().Ping()
	gormDb.DB().SetMaxIdleConns(10)
	gormDb.DB().SetMaxOpenConns(100)
	// Create if not exists
	gormDb.CreateTable(&User{})
	gormDb.CreateTable(&TeamMembership{})
	gormDb.CreateTable(&Team{})
	gormDb.CreateTable(&Organization{})
	// Migrations
	gormDb.AutoMigrate(&User{}, &TeamMembership{}, &Team{}, &Organization{})
}

func getDB() gorm.DB {
	return gormDb
}
