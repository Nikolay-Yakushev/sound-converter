package checkers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "fmt"
)

var maxFileSize int = 1024 * 1024 *20

func ErrCheck(err error) {
    if err != nil{
        panic(err)
    }
}

func CheckIsDocument(update *tgbotapi.Update) bool {
    if update.Message.Document == nil {
        return false
    }
    return true
}

func IsOkDocSize(update *tgbotapi.Update) bool {
    fmt.Printf("File size %d\n", update.Message.Document.FileSize)
    if update.Message.Document.FileSize <= maxFileSize {
        return true
    }
    return false
}
 