package botgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type BotSocket struct {
	bot        *Bot
	connection *websocket.Conn

	stopSignal chan bool
	errChan    chan error

	lastSChan chan int
	lastS     int

	readChan  chan Payload
	writeChan chan Payload

	messageChan chan Message
}

const BOT_SOCKET_CLOSED_ERROR string = "Bot socket closed."

func CreateBotSocket(bot *Bot) (botSocket *BotSocket, err error) {
	botSocket = new(BotSocket)

	botSocket.stopSignal = make(chan bool, 1)
	botSocket.errChan = make(chan error, 1)

	botSocket.lastSChan = make(chan int, 1)

	botSocket.readChan = make(chan Payload, 10)
	botSocket.writeChan = make(chan Payload, 10)

	botSocket.messageChan = make(chan Message, 10)

	var gateway Gateway
	gateway, err = bot.GetGateway()
	if err != nil {
		return nil, err
	}

	botSocket.connection, _, err = websocket.DefaultDialer.Dial(gateway.Url, nil)
	if err != nil {
		return nil, err
	}

	botSocket.connection.SetReadLimit(int64(BOT_MAX_MESSAGE_SIZE))

	var payload Payload
	err = botSocket.connection.ReadJSON(&payload)
	if err != nil {
		return nil, err
	}

	if payload.Op != 10 {
		return nil, errors.New("Cannot connect to the botSocket server.")
	}

	var d struct {
		Token      string `json:"token"`
		Intents    int    `json:"intents"`
		Shard      []int  `json:"shard"`
		Properties []any  `json:"properties"`
	}
	d.Token = bot.Authorization
	d.Intents = 1 << 12
	d.Shard = []int{0, 1}
	d.Properties = nil

	payload.Op = 2
	payload.D = d
	payload.S = 0
	payload.T = ""
	err = botSocket.connection.WriteJSON(payload)
	if err != nil {
		return nil, err
	}

	err = botSocket.connection.ReadJSON(&payload)
	if err != nil {
		return nil, err
	}

	if payload.Op != 0 {
		return nil, errors.New("Identify failed.")
	}

	go botSocket.heartBeatLoop()
	go botSocket.readLoop()
	go botSocket.writeLoop()
	go botSocket.mainLoop()

	return botSocket, nil
}

func (botSocket *BotSocket) spreadStop() {
	select {
	case botSocket.stopSignal <- true:
	default:
	}
}

func (botSocket *BotSocket) setLastS(lastS int) {
	botSocket.lastSChan <- lastS
}

func (botSocket *BotSocket) getLastS() (lastS int) {
	select {
	case botSocket.lastS = <-botSocket.lastSChan:
	default:
	}

	lastS = botSocket.lastS
	return lastS
}

func (botSocket *BotSocket) heartBeatLoop() {
	for true {
		select {
		case _ = <-botSocket.stopSignal:
			botSocket.spreadStop()
			return
		default:
		}

		var payload Payload
		payload.Op = 1
		payload.D = botSocket.getLastS()

		botSocket.writeChan <- payload

		time.Sleep(BOT_HEART_BEAT_PERIOD)
	}
}

func (botSocket *BotSocket) writeLoop() {
	var err error

	for true {
		select {
		case _ = <-botSocket.stopSignal:
			botSocket.spreadStop()
			return
		case v := <-botSocket.writeChan:
			err = botSocket.connection.WriteJSON(v)
			if err != nil {
				botSocket.Close(err)
				return
			}
		}
	}
}

func (botSocket *BotSocket) readLoop() {
	var err error

	for true {
		select {
		case _ = <-botSocket.stopSignal:
			botSocket.spreadStop()
			return
		default:
			var payload Payload
			err = botSocket.connection.ReadJSON(&payload)
			if err != nil {
				botSocket.Close(err)
				return
			}

			fmt.Println(payload)
			botSocket.readChan <- payload
		}
	}
}

func getFromAny(dst any, src any) (err error) {
	var bytes []byte

	bytes, err = json.Marshal(src)
	if err != nil {
		return err
	}

	json.Unmarshal(bytes, &dst)
	if err != nil {
		return err
	}

	return
}

func (botSocket *BotSocket) mainHandler(payload Payload) (err error) {
	switch payload.Op {
	case 0:
		botSocket.setLastS(payload.S)

		if payload.T == "DIRECT_MESSAGE_CREATE" {
			var message Message
			err = getFromAny(&message, payload.D)
			if err != nil {
				return err
			}

			botSocket.messageChan <- message
		} else {
			return errors.New("Invalid notification type.")
		}
	case 11:
	default:
		return errors.New("Invalid op code")
	}
	return nil
}

func (botSocket *BotSocket) mainLoop() {
	var err error

	for true {
		select {
		case _ = <-botSocket.stopSignal:
			botSocket.spreadStop()
			return
		case payload := <-botSocket.readChan:
			err = botSocket.mainHandler(payload)
			if err != nil {
				botSocket.Close(err)
				return
			}
		}
	}
}

func (botSocket *BotSocket) Close(err error) {
	select {
	case botSocket.errChan <- err:
	default:
	}

	botSocket.spreadStop()
	botSocket.connection.Close()
}

func (botSocket *BotSocket) ReadError() (err error) {
	err = <-botSocket.errChan
	return err
}

func (botSocket *BotSocket) ReadMessage() (message Message, err error) {
	select {
	case _ = <-botSocket.stopSignal:
		botSocket.spreadStop()
		err = errors.New(BOT_SOCKET_CLOSED_ERROR)
	case message = <-botSocket.messageChan:
	}
	return message, err
}
