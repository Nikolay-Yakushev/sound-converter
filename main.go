package main

import (
	"bufio"
	"context"
	"log"
	"os"

	checkers "github.com/Nikolay-Yakushev/sound-converter/checkers"
	loggers "github.com/Nikolay-Yakushev/sound-converter/loggers"
	receivers "github.com/Nikolay-Yakushev/sound-converter/receivers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//todo remove when docker
const TELEGRAM_TOKEN = "5903498628:AAEUyizBqXKO76eq7J9JJ9kyEuY6U7dRgYY"
const BASE_DIRECTORY = "/home/driver220v/Nsound/sound-converter/"
const DEBUG = true
const UPDATE_TIMEOUT = 30


func formTextResponse(update *tgbotapi.Update, bot *tgbotapi.BotAPI, replyText string) {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, replyText)
    msg.ReplyToMessageID = update.Message.MessageID
    if _, err := bot.Send(msg); err != nil {
        panic(err)
    }
}


func handleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	switch {
        case update.Message != nil:

            var replyText string
            document := checkers.CheckIsDocument(&update)

            if !document{
                replyText  = "File should be document."
                formTextResponse(&update, bot, replyText)
                break
            }
            if !checkers.IsOkDocSize(&update){
                replyText  = "File is biger than 20mb."
                formTextResponse(&update, bot, replyText)
                break
            }
            go receivers.SaveMessageData(update, bot)
        }
	}


func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI) {
	for {
		select {
		    case <-ctx.Done():
                if err := ctx.Err(); err != nil {
                    loggers.ErrLogger.Printf("receiveUpdates err: %s\n", err)
                }
			    return
            case update := <-updates:
                go handleUpdate(update, bot)
		} 
    }
}


func main() {
	bot, err := tgbotapi.NewBotAPI(TELEGRAM_TOKEN)
    checkers.ErrCheck(err)

    bot.Debug = DEBUG
	loggers.InfLogger.Printf(
        "Authorized on account %s. DEBUG options is set to `%t`", 
        bot.Self.UserName, DEBUG)

    ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

    updateConfig := tgbotapi.NewUpdate(0)
    updateConfig.Timeout = UPDATE_TIMEOUT

	updates := bot.GetUpdatesChan(updateConfig)
    go receiveUpdates(ctx, updates, bot)

    log.Println("Start listening for updates. Press enter to stop")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

    cancel()
           
}