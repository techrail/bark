package tests

//
// import (
//	"encoding/json"
//	"log"
//	"testing"
//	"time"
//
//	"github.com/joho/godotenv"
//	"github.com/techrail/bark/db"
//	"github.com/techrail/bark/models"
// )
//
// func Init() {
//	err := godotenv.Load("../.env")
//	if err != nil {
//		log.Fatal("Error loading .env file")
//	}
// }
// func TestDBConnection(t *testing.T) {
//	Init()
//	_, err := db.ConnectToDatabase()
//	if err != nil {
//		log.Fatalf("Failed to connect to the database: %v", err)
//	}
// }
//
// func TestInsertLog(t *testing.T) {
//	Init()
//
//	db, err := db.ConnectToDatabase()
//
//	moreData, _ := json.Marshal(map[string]interface{}{
//		"a": "apple",
//		"b": "banana",
//	},
//	)
//	sampleLog := models.BarkLog{
//		LogTime:     time.Now(),
//		LogLevel:    0,
//		AppName: "test",
//		Code:        "1234",
//		Message:     "Test",
//		MoreData:    moreData}
//
//	if err != nil {
//		log.Fatalf("Failed to connect to the database: %v", err)
//	}
//	err = db.InsertLog(sampleLog)
//	if err != nil {
//		log.Fatal(err)
//	}
// }
//
// func TestInsertBatch(t *testing.T) {
//	Init()
//
//	db, err := db.ConnectToDatabase()
//	moreData, _ := json.Marshal(map[string]interface{}{
//		"a": "apple",
//		"b": "banana",
//	},
//	)
//	sampleLogs := []models.BarkLog{
//		// Id:          1234,
//		{LogTime: time.Now(),
//			LogLevel:    0,
//			AppName: "test",
//			Code:        "1234",
//			Message:     "Test",
//			MoreData:    moreData},
//		{LogTime: time.Now(),
//			LogLevel:    0,
//			AppName: "test",
//			Code:        "1234",
//			Message:     "Test",
//			MoreData:    moreData},
//		{LogTime: time.Now(),
//			LogLevel:    0,
//			AppName: "test",
//			Code:        "1234",
//			Message:     "Test",
//			MoreData:    moreData},
//		{LogTime: time.Now(),
//			LogLevel:    0,
//			AppName: "test",
//			Code:        "1234",
//			Message:     "Test",
//			MoreData:    moreData},
//	}
//	if err != nil {
//		log.Fatalf("Failed to connect to the database: %v", err)
//	}
//	err = db.InsertBatch(sampleLogs)
//	if err != nil {
//		log.Fatal(err)
//	}
// }
//
// func TestFetchLogs(t *testing.T) {
//	Init()
//
//	db, err := db.ConnectToDatabase()
//	if err != nil {
//		log.Fatal(err)
//	}
//	_, err = db.FetchLimitedLogs(10)
//	if err != nil {
//		log.Fatal(err)
//	}
// }
