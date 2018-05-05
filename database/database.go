package database

import (
	"log"
	"os"

	"github.com/RTradeLtd/RTC-IPFS/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	rollbar "github.com/rollbar/rollbar-go"
)

/*
	roll.Token = "POST_SERVER_ITEM_ACCESS_TOKEN"
	//roll.Environment = "production" // defaults to "development"

	r := gin.Default()
	r.Use(rollbar.Recovery(true))

	r.Run(":8080")
	func l(err error) {
	token := os.Getenv("ROLLBAR_TOKEN")
	rollbar.SetToken(token)
	rollbar.SetServerRoot("github.com/RTradeLtd/RTC-IPFS") // path of project (required for GitHub integration and non-project stacktrace collapsing)

	rollbar.Error(err)

	rollbar.Wait()
}
*/
var db *gorm.DB

func rollbarError(err error) {
	token := os.Getenv("ROLLBAR_TOKEN")
	rollbar.SetToken(token)
	rollbar.SetServerRoot("github.com/RTradeLtd/RTC-IPFS")
	rollbar.Error(err)
	rollbar.Wait()
}

func RunMigrations() {
	var uploads models.Upload
	db := OpenDBConnection()
	db.AutoMigrate(uploads)
}

// OpenDBConnection is used to create a database connection
func OpenDBConnection() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./ipfs_database.db")
	if err != nil {
		rollbarError(err)
		log.Fatal(err)
	}
	return db
}

// CloseDBConnection is used to close a db
func CloseDBConnection(db *gorm.DB) {
	db.Close()
}

// GetUploads is used to retrieve all uploads
func GetUploads() []*models.Upload {
	var uploads []*models.Upload
	db = OpenDBConnection()
	db.Find(&uploads)
	return uploads
}

// AddHash his used to add a hash to our database
func AddHash(c *gin.Context) {
	var upload models.Upload
	hash := c.Param("hash")
	upload.Hash = hash
	upload.Type = "pin"
	db := OpenDBConnection()
	db.Create(&upload)
	db.Close()
}

// AddFileHash is used to add the hash of a file to our database
func AddFileHash(hash string) {
	var upload models.Upload
	upload.Hash = hash
	upload.Type = "file"
	db := OpenDBConnection()
	db.AutoMigrate(&upload)
	db.Create(&upload)
	db.Close()
}