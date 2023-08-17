package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"tencent-bot/internal/pkg/botgo"
	"time"
)

const BotAppId string = "102061205"
const BotToken string = "Ilw6R4JPoiIhOmw6AvsOk5AUDlJCOawi"

func Reply(message botgo.Message) (replyContent string) {
	var err error
	var content string = message.Content
	var parts []string = strings.Split(content, " ")
	switch parts[0] {
	case "h":
		if len(parts) != 1 {
			break
		}
		return `发送 “h” 获取使用帮助
发送 “d [n]” 获取 1 ~ n 的随机整数（例：d 3）`
	case "d":
		const invalidUsageString string = "用法错误，正确用法：d [n]"

		if len(parts) != 2 {
			return invalidUsageString
		}

		var maxValue int
		maxValue, err = strconv.Atoi(parts[1])
		if err != nil {
			return invalidUsageString
		}

		if maxValue <= 0 {
			return "n 必须为正整数"
		}

		return "" + strconv.Itoa(rand.Intn(maxValue)+1)
	}
	return "小助手还不能明白您在说什么哦，发送 “h” 来获取小助手的用法。"
}

func RunBot() (err error) {
	var client http.Client
	client.Timeout = 3 * time.Second

	var bot *botgo.Bot
	bot = botgo.CreateBot(false, BotAppId, BotToken)

	var botSocket *botgo.BotSocket
	botSocket, err = botgo.CreateBotSocket(bot)
	if err != nil {
		return err
	}

	defer botSocket.Close(nil)

	var message botgo.Message
	for true {
		message, err = botSocket.ReadMessage()
		if err != nil {
			err = botSocket.ReadError()
			return err
		}

		_, err = bot.SendDirectMessage(message.GuildId, message.Id, Reply(message))
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	rand.Seed(time.Now().Unix())
	var err error
	err = RunBot()
	if err != nil {
		fmt.Println(err.Error())
	}
}
