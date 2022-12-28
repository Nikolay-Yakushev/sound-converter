package receivers

import (
	"path/filepath"
	"encoding/binary"
	"net/http"
	"os"
	"io"
	checkers "github.com/Nikolay-Yakushev/sound-converter/checkers"
	loggers "github.com/Nikolay-Yakushev/sound-converter/loggers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
//todo remove when docker
const TELEGRAM_TOKEN = "5903498628:AAEUyizBqXKO76eq7J9JJ9kyEuY6U7dRgYY"
const BASE_DIRECTORY = "/home/driver220v/Nsound/sound-converter/"
const DEBUG = true
const UPDATE_TIMEOUT = 30

func createFile(update *tgbotapi.Update) *os.File {
    userDir := filepath.Join(
        BASE_DIRECTORY,
        "USERS",
        update.Message.Chat.UserName)

    err := os.MkdirAll(userDir, os.ModePerm)
    if err == nil{
        loggers.InfLogger.Printf("Directory: %s already exist", userDir)
    }

    userFileName := filepath.Join(userDir, update.Message.Document.FileName)
    file, err := os.Create(userFileName)
    checkers.ErrCheck(err)
    return file

}

func readFile(resp* http.Response, messages chan []byte) {
    var bufferSize int = 2048
    for {
        buffer := make([]byte, bufferSize)
        _, err := io.ReadAtLeast(resp.Body, buffer, bufferSize)
        if err != nil {

            if err == io.EOF{
                break
            }

            if err == io.ErrUnexpectedEOF{
                messages <- buffer
                break

            }else{
                loggers.ErrLogger.Printf("Error has occured:%s", err.Error())
                break
            }

        }else{
            messages <- buffer
        }
        
    }
    close(messages)
}

func SaveMessageData(update tgbotapi.Update, bot *tgbotapi.BotAPI) {

    fileConfig := tgbotapi.FileConfig{FileID: update.Message.Document.FileID}
    file, err := bot.GetFile(fileConfig)
    checkers.ErrCheck(err)

    resp, err := http.Get(file.Link(TELEGRAM_TOKEN))
    checkers.ErrCheck(err)
    defer resp.Body.Close()

    loggers.InfLogger.Printf("Received Filename %s", update.Message.Document.FileName)
    if resp.StatusCode == http.StatusOK{

        messages := make(chan []byte, 8192)
        go readFile(resp, messages)

        file := createFile(&update)
        defer file.Close()

        for msg := range messages{
            binary.Write(file, binary.BigEndian, msg)
        }
        
        loggers.InfLogger.Printf("Finished saving %s", update.Message.Document.FileName)
    } 
}

