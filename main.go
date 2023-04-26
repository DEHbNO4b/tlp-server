package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var login string = `{ "id": 0, "stream": "ee46d1d0-e8e0-4058-be66-54846fd278cf" }`

//var login string = `{ "id": 0, "stream": "4b205c0a-2fd6-4e5e-b4d4-d31eb0e43918" }`

type date struct {
	year  int
	month int
	day   int
}

func main() {
	conn, err := net.Dial("tcp", "192.168.1.4:8082")
	checkError(err)
	defer conn.Close()
	go func() {
		for {
			fmt.Fprintf(conn, login+"\r"+"\n")
			time.Sleep(10 * time.Second)
		}
	}()
	strokeChan := make(chan string)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go lightningWriter(ctx, strokeChan)

	buf := make([]byte, 1024)
	b := []byte("9	KEEP	ALIVE")

	for {
		// Прослушиваем ответ
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			time.Sleep(60 * time.Second)
			continue
		}

		//пишем ответ в БД
		//пишем ответ в файл

		if string(buf[:len(b)]) != string(b) {
			stroke := string(buf[:n])
			strokeChan <- stroke
		}
		//пишем ответ для АСПД
		fmt.Print("Message from server: " + string(buf[:n]))
	}
}
func lightningWriter(ctx context.Context, r <-chan string) {

	//file, err := os.OpenFile(s[2]+s[3]+s[4]+".txt", os.O_WRONLY|os.O_CREATE, 0666)
	var file *os.File
	var filename string
	var err error
	writeCounter := 0
	for {
		select {
		case <-ctx.Done():
			return
		case stroke := <-r:
			s := strings.Split(stroke, "\t")
			if file == nil || filename != s[2]+s[3]+s[4]+".txt" {
				filename = s[2] + s[3] + s[4] + ".txt"
				file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
				checkError(err)
			}
			n, err := file.WriteString(stroke)
			checkError(err)
			writeCounter++
			fmt.Printf("writeCounter = %d, n = %d \n", writeCounter, n)

		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	}
}
