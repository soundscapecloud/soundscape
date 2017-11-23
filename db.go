package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	//  _ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	//  _ "github.com/jinzhu/gorm/dialects/mssql"
	//"fmt"
)

// User ...
/*type User struct {
	ID       uint
	Username string
	Password string
}*/

var db *gorm.DB

func dbInit() {
	var err error
	// Switch database
	if os.Getenv("MYSQL_DB") != "" {
		db, err = gorm.Open("mysql", os.Getenv("MYSQL_USER")+":"+os.Getenv("MYSQL_PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":"+os.Getenv("MYSQL_PORT")+")/"+os.Getenv("MYSQL_DB")+"?charset=utf8&parseTime=True&loc=Local")
	} else {
		db, err = gorm.Open("sqlite3", "./streamlist.db")
	}
	if err != nil {
		panic("failed to connect database")
	}
	//defer db.Close()
	db.SingularTable(true)
	//db.AutoMigrate(&User{})
	db.AutoMigrate(&Media{})
	db.AutoMigrate(&List{})
	db.AutoMigrate(&User{})

	// Add / verify foreign keys
	db.Exec("ALTER TABLE `list_media` ADD CONSTRAINT `fk_list_id` FOREIGN KEY (`list_id`) REFERENCES `list` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, ADD CONSTRAINT `fk_media_id` FOREIGN KEY (`media_id`) REFERENCES `media` (`id`) ON UPDATE CASCADE ON DELETE CASCADE;")
}
