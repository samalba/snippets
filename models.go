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
	Id                 int64
	Name               string `sql:"type:varchar(100);"`
	Description        string `sql:"type:varchar(255);"`
	Url                string `sql:"type:varchar(100);"`
	Email              string `sql:"type:varchar(100);"`
	Teams              []Team
	OrganizationAdmins []OrganizationAdmin
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          time.Time
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

type OrganizationAdmin struct {
	Id             int64
	OrganizationId int64
	UserId         int64
}

type User struct {
	Id                 int64
	Login              string `sql:"type:varchar(100);"`
	Name               string `sql:"type:varchar(100);"`
	Company            string `sql:"type:varchar(100);"`
	Email              string `sql:"type:varchar(100);"`
	AvatarURL          string `sql:"type:varchar(100);"`
	Location           string `sql:"type:varchar(100);"`
	SuperAdmin         bool
	TeamMemberships    []TeamMembership    `json:"-"`
	OrganizationAdmins []OrganizationAdmin `json:"-"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          time.Time
}

type Snippet struct {
	Id        int64
	Week      int64 // Week number
	Year      int64
	Content   string `sql:"type:text;"`
	Published bool
	UserId    int64 // ForeignKey to a User
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

var gormDb *gorm.DB

func init() {
	db, err := gorm.Open("sqlite3", "snippets.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot init the db: %s\n", err)
		os.Exit(1)
	}
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	// Create if not exists
	db.CreateTable(&User{})
	db.CreateTable(&TeamMembership{})
	db.CreateTable(&Team{})
	db.CreateTable(&Organization{})
	db.CreateTable(&Snippet{})
	// Indexes
	db.Model(&User{}).AddUniqueIndex("idx_user_login", "login")
	// Migrations
	db.AutoMigrate(&User{}, &TeamMembership{}, &Team{}, &Organization{}, &Snippet{})
	gormDb = &db
}

func getDB() *gorm.DB {
	return gormDb
}
