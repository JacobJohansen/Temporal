package models

import (
	"log"

	"github.com/jinzhu/gorm"
)

type Upload struct {
	gorm.Model
	Hash             string `gorm:not null;`
	Type             string `gorm:not null;` //  file, pin
	HoldTimeInMonths int    `gorm:not null;`
	UploadAddress    string `gorm:not null;`
}

type UploadManager struct {
	DB *gorm.DB
}

// NewUploadManager is used to generate an upload manager interface
func NewUploadManager(db *gorm.DB) *UploadManager {
	return &UploadManager{DB: db}
}

// FindUploadsByHash is used to return all instances of uploads matching the
// given hash
func (um *UploadManager) FindUploadsByHash(hash string) []*Upload {

	uploads := []*Upload{}

	um.DB.Find(&uploads).Where("hash = ?", hash)

	return uploads
}

// GetUploadByHashForUploader is used to retrieve the last (most recent) upload for a user
func (um *UploadManager) GetUploadByHashForUploader(hash string, uploaderAddress string) []*Upload {
	var uploads []*Upload
	um.DB.Find(&uploads).Where("hash = ? AND uploader_address = ?", hash, uploaderAddress)
	return uploads
}

// GetUploads is used to return all  uploads
func (um *UploadManager) GetUploads() []*Upload {
	var uploads []*Upload
	um.DB.Find(uploads)
	return uploads
}

// GetUploadsForAddress is used to retrieve all uploads by an address
func (um *UploadManager) GetUploadsForAddress(address string) []*Upload {
	var uploads []*Upload
	um.DB.Where("upload_address = ?", address).Find(&uploads)
	return uploads
}

// AddPinHash is used to upload a pin hash to our database
func (um *UploadManager) AddPinHash(hash string, uploaderAddress string, holdTimeInMonths int) {
	var upload Upload
	upload.HoldTimeInMonths = holdTimeInMonths
	upload.UploadAddress = uploaderAddress
	upload.Hash = hash
	upload.Type = "pin"
	um.DB.Create(&upload)
}

// AddFileHash is used to add the hash of a file to our database
func (um *UploadManager) AddFileHash(hash string, uploaderAddress string, holdTimeInMonths int) {
	var upload Upload
	upload.HoldTimeInMonths = holdTimeInMonths
	upload.UploadAddress = uploaderAddress
	upload.Hash = hash
	upload.Type = "file"
	um.DB.Create(&upload)
}

// OpenDBConnection is used to open a database connection
func OpenDBConnection() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./ipfs_database.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
