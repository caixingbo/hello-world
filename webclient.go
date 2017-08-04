package main

import (
	"github.com/gorilla/websocket"
	"strconv"
	"log"
	"fmt"
)

func main() {
	dialer := websocket.Dialer{ /* set fields as needed */ }
	ws, _, err := dialer.Dial("ws://127.0.0.1:7058/ws", nil)

	if err != nil {
		log.Printf("1%v",err)
	}
	msg := []byte("{\"MsgID\":\"" + strconv.Itoa(20080) + "\",\"Name\":\"lishi\"," +
		"\"Token\":\"TWpBeE4yUjNaR0YwWVdObGJuUmxjZz09YkdsemFHaz0=\"}")

	fmt.Println(string(msg))

	if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("2%v",err)
	}
	_, p, err := ws.ReadMessage()
	if err != nil {
		log.Printf("3%v",err)
	}

	fmt.Println("receive message ", string(p))

	msg = []byte("{\"MsgID\":\"" + strconv.Itoa(20310) + "\",\"Name\":\"lishi\"}")
	fmt.Println(string(msg))

	if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		log.Printf("4%v",err)
	}

	if err != nil {
		log.Printf("5%v",err)
	}
	_, p, err = ws.ReadMessage()
	if err != nil {
		log.Printf("6%v",err)
	}
	fmt.Println("receive message ", string(p))
}
