package botgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Bot struct {
	Url           string
	Authorization string
	httpClient    http.Client
}

const BOT_TIMEOUT time.Duration = 3 * time.Second
const BOT_HEART_BEAT_PERIOD time.Duration = 10 * time.Second
const BOT_MAX_MESSAGE_SIZE int = 1 << 20

const BOT_SAND_BOX_URL string = "https://sandbox.api.sgroup.qq.com"
const BOT_URL string = "https://api.sgroup.qq.com"

func CreateBot(inSandBox bool, botAppId string, botToken string) (bot *Bot) {
	bot = new(Bot)
	if inSandBox {
		bot.Url = BOT_SAND_BOX_URL
	} else {
		bot.Url = BOT_URL
	}
	bot.Authorization = "Bot " + botAppId + "." + botToken
	bot.httpClient.Timeout = BOT_TIMEOUT
	return bot
}

func (bot *Bot) SendRequest(method string, path string, requestBody any, responseBody any) (err error) {
	var requestBodyBytes []byte
	requestBodyBytes, err = json.Marshal(requestBody)
	if err != nil {
		return err
	}

	var request *http.Request
	request, err = http.NewRequest(method, bot.Url+path, bytes.NewReader(requestBodyBytes))
	if err != nil {
		return err
	}

	request.Header.Add("accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", bot.Authorization)

	var response *http.Response
	response, err = bot.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer bot.httpClient.CloseIdleConnections()
	defer response.Body.Close()

	if responseBody == nil {
		return nil
	}

	var responseBodyBytes []byte
	responseBodyBytes, err = ReadAll(response.Body, 4<<10, BOT_MAX_MESSAGE_SIZE)
	if err == nil {
		err = errors.New("Response too long.")
	}
	if err.Error() == "EOF" {
		err = nil
	}
	if err != nil {
		return err
	}

	json.Unmarshal(responseBodyBytes, responseBody)
	if err != nil {
		return err
	}

	return nil
}

func (bot *Bot) GetGuilds() (guilds []Guild, err error) {
	err = bot.SendRequest("GET", "/users/@me/guilds", nil, &guilds)
	return guilds, err
}

func (bot *Bot) GetGateway() (gateway Gateway, err error) {
	bot.SendRequest("GET", "/gateway", nil, &gateway)
	return gateway, err
}

func (bot *Bot) SendDirectMessage(guildId string, messageId string, content string) (message Message, err error) {
	var requestBody struct {
		Content string `json:"content"`
		MsgId   string `json:"msg_id"`
	}
	requestBody.Content = content
	requestBody.MsgId = messageId
	bot.SendRequest("POST", "/dms/"+guildId+"/messages", requestBody, &message)
	return message, err
}
