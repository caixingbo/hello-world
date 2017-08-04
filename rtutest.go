// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license.  See the LICENSE file for details.

package main

import (
	"log"
	"os"
	"testing"

	"github.com/goburrow/modbus"
	"time"
	//"encoding/binary"
	//"bytes"
	//"encoding/binary"

	//"fmt"
)

const (
	rtuDevice = "COM3"
	tcpDevice = "192.168.0.210:502"
)

func main()  {
	t := &testing.T{}
	//go TestRTUClientAdvancedUsage(t)
	go TestTCPClientAdvancedUsage(t)
	time.Sleep(20*time.Second)

	log.Printf("test1")
}


func TestRTUClientAdvancedUsage(t *testing.T) {
	handler := modbus.NewRTUClientHandler(rtuDevice)
	handler.BaudRate = 19200
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 40
	handler.Logger = log.New(os.Stdout, "rtu: ", log.LstdFlags)
	err := handler.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer handler.Close()

	client := modbus.NewClient(handler)
	results, err := client.ReadInputRegisters(0, 3)
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
	for index := range results {
		println("results:",results[index])
	}
	results, err = client.ReadWriteMultipleRegisters(0, 2, 2, 2, []byte{1, 2, 3, 4})
	if err != nil || results == nil {
		t.Fatal(err, results)
	}
}

func toHex(ten int16) (hex []int16, length int) {
	var m int16 = 0

	hex = make([]int16, 0)
	length = 0;

	for{
		m = ten / 16
		ten = ten % 16

		if(m == 0){
			hex = append(hex, ten)
			length++
			break
		}

		hex = append(hex, m)
		length++;
	}
	return
}

func TestTCPClientAdvancedUsage(t *testing.T) {
	handler := modbus.NewTCPClientHandler(tcpDevice)
	handler.Timeout = 2 * time.Second
	handler.SlaveId = 1
	handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)
	handler.Connect()
	defer handler.Close()

	client := modbus.NewClient(handler)


	/*
	log.Println("read begin...")

	results, err := client.ReadHoldingRegisters(11, 1)
	//results, err := client.ReadHoldingRegisters(0, 10)
	//b:=[]byte{1,1}
	//b:=[]byte{0x00,0x00,0x03,0xe8}
	//b_buf:=bytes.NewBuffer(b)




	b_buf:=bytes.NewBuffer(results)
	var x int16
	binary.Read(b_buf,binary.BigEndian,&x)
	log.Printf("int ...%v",x)

	a := x
	hex,length := toHex(a)

	for i:=0; i < length; i++ {
		if(hex[i] >= 10){
			fmt.Printf("mac:%c",'A'+hex[i]-10)
			mac := 'A'+hex[i]-10
			log.Printf("mac1 ...%v",mac)
		} else{
			fmt.Printf("mac:%c",hex[i])
			mac := hex[i]
			log.Printf("mac1 ...%v",mac)
		}
	}
*/

	//results, err := client.ReadInputRegisters(10600, 6)
	/*
	var i = 0
	for i<23 {
		results, _ := client.ReadCoils(1100+ uint16(i),1)
		i = i+1
		for index := range results {
			log.Println("coil:",results[index],i)
		}
	}
	*/
	//results, err := client.ReadCoils(1100,16)

	//log.Println("read end...")
	//results, err := client.ReadInputRegisters(10800, 13)
/*
	if err != nil || results == nil {
		log.Printf("err,%v",err)
		t.Fatal(err, results)
	}

	log.Printf("read...%v",results)

	for index := range results {
		log.Println("results1:",int(results[index]),index)
	}
*/
	//_, err := client.WriteSingleCoil(uint16(1010),0xFF00)
	_, err := client.WriteMultipleCoils(uint16(1010), 1, []byte{1})
	log.Println("begin")
	if err != nil {
		log.Printf("err,%v", err)
	}

	//time.Sleep(500*time.Microsecond)
	time.Sleep(1*time.Second)
	log.Println("end")
	_, err = client.WriteMultipleCoils(uint16(1010), 1, []byte{0})
	//_, err = client.WriteSingleCoil(uint16(1010),0x0000)
	if err != nil {
		log.Printf("err,%v", err)
	}

	/*


	//results, err = client.WriteMultipleRegisters(0, 3, []byte{0,10,0,11,0,15})
	results, err = client.WriteMultipleCoils(0, 2, []byte{3})

	if err != nil || results == nil {
		t.Fatal(err, results)
	}

	for index := range results {
		println("results2:",results[index])
	}
	*/
}
