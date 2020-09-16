package db

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // import global
	"github.com/kimtaek/gamora/pkg/slack"
	"math"
	"os"
	"time"
)

// Model global basic model struct
type Model struct {
	ID        uint64     `form:"id" json:"id" gorm:"primary_key"`
	CreatedBy uint64     `json:"-"`
	UpdatedBy uint64     `json:"-"`
	DeletedBy uint64     `json:"-"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

// Configure config for db
type Configure struct {
	Database string `env:"MYSQL_DATABASE" envDefault:"database"`
	Username string `env:"MYSQL_USERNAME" envDefault:"user"`
	Password string `env:"MYSQL_PASSWORD" envDefault:"password"`
	Host     string `env:"MYSQL_HOST" envDefault:"localhost"`
	Port     string `env:"MYSQL_PORT" envDefault:"3306"`
}

// Config global defined db config
var Config Configure

// DB global defined db config
var DB *gorm.DB

// Setup init db
func Setup() {
	_ = env.Parse(&Config)
	connection, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Config.Username,
		Config.Password,
		Config.Host,
		Config.Port,
		Config.Database,
	))

	if err != nil {
		slack.SendMessage(slack.Message{
			Text: "Database: " + err.Error(),
		})
		os.Exit(1)
	}

	connection.LogMode(true)
	connection.DB().SetConnMaxLifetime(time.Minute * 3)
	DB = connection

	_, _ = color.New(color.FgWhite).Println(time.Now().Format(time.RFC3339), "[info]", "[database connected!]")
}

// Connection get db connection
func Connection() *gorm.DB {
	return DB
}

// CloseDB close db connection
func CloseDB() {
	defer DB.Close()
}

// PaginationParam for using pagination
type PaginationParam struct {
	DB          *gorm.DB
	Page        int    `form:"page" json:"page"`
	Limit       int    `form:"limit" json:"limit"`
	OrderBy     string `form:"orderBy" json:"orderBy"`
	OrderBySort string `form:"orderBySort" json:"orderBySort"`
}

// Pagination for using pagination
type Pagination struct {
	Data interface{} `json:"data"`
	Meta struct {
		Total     int `json:"total"`
		TotalPage int `json:"totalPage"`
		Offset    int `json:"offset"`
		Limit     int `json:"limit"`
		Page      int `json:"page"`
		PrevPage  int `json:"prevPage"`
		NextPage  int `json:"nextPage"`
	} `json:"meta"`
}

// Paginate for using pagination
func Paginate(p *PaginationParam, dataSource interface{}) *Pagination {
	db := p.DB

	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 25
	}

	done := make(chan bool, 1)
	var pagination Pagination
	var count int
	var offset int

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	go totalCount(db, dataSource, done, &count)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	db.Limit(p.Limit).Offset(offset).Find(dataSource)
	<-done
	pagination.Meta.Total = count
	pagination.Data = dataSource
	pagination.Meta.Page = p.Page

	pagination.Meta.Offset = offset
	pagination.Meta.Limit = p.Limit
	pagination.Meta.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		pagination.Meta.PrevPage = p.Page - 1
	} else {
		pagination.Meta.PrevPage = p.Page
	}

	if p.Page == pagination.Meta.TotalPage {
		pagination.Meta.NextPage = p.Page
	} else {
		pagination.Meta.NextPage = p.Page + 1
	}
	return &pagination
}

// totalCount for using pagination
func totalCount(db *gorm.DB, countDataSource interface{}, done chan bool, count *int) {
	db.Model(countDataSource).Count(count)
	done <- true
}
