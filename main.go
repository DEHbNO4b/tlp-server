package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	address string = "192.168.1.4:8082"
	login   string = `{ "id": 0, "stream": "ee46d1d0-e8e0-4058-be66-54846fd278cf" }`
)

func main() {
	//creating logger
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	for {
		//creating context
		ctx, cancel := context.WithCancel(context.Background())

		err := tlpConnect(ctx, log)
		log.Error("take error from tlpConnect func:", slog.String("error", err.Error()))
		time.Sleep(30 * time.Second)

		cancel()
	}
}
func tlpConnect(ctx context.Context, log *slog.Logger) error {

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("unable to create connection: %w", err)
	}
	defer conn.Close()

	//sending message "keep alive" to server tlp
	go sendLogin(ctx, log, conn)

	strokeChan := make(chan string)

	go lightningWriter(ctx, log, strokeChan)

	buf := make([]byte, 1024)
	b := []byte("9	KEEP	ALIVE")

	for {
		select {
		case <-ctx.Done():
			log.Info("tlpConnect have done work")
		default:
			// Прослушиваем ответ
			n, err := conn.Read(buf)
			if err != nil {
				log.Error("error ocurred while reading connection:", slog.String("error", err.Error()))
				time.Sleep(60 * time.Second)
				return err
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
}

func lightningWriter(ctx context.Context, log *slog.Logger, r <-chan string) {

	//file, err := os.OpenFile(s[2]+s[3]+s[4]+".txt", os.O_WRONLY|os.O_CREATE, 0666)
	var file *os.File
	var filename string
	var err error
	writeCounter := 0

	for {
		select {
		case <-ctx.Done():
			log.Info("lightningWriter have done work")
			return
		case stroke := <-r:
			s := strings.Split(stroke, "\t")
			if file == nil || filename != s[2]+s[3]+s[4]+".txt" {
				filename = s[2] + s[3] + s[4] + ".txt"
				file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
				if err != nil {
					log.Error("error ocurred in attempt to open file", slog.String("err", err.Error()))
					continue
				}
			}
			_, err := file.WriteString(stroke)
			if err != nil {
				log.Error("error ocurred while writting string to file", slog.String("err", err.Error()))
			}
			writeCounter++
			log.Info("have got a string", slog.Attr{
				Key:   "counter",
				Value: slog.StringValue(strconv.Itoa(writeCounter)),
			})

		}
	}
}

func sendLogin(ctx context.Context, log *slog.Logger, conn net.Conn) {
	for {
		select {
		case <-ctx.Done():
			log.Info("sendLogin have done work")
			return
		default:

			_, err := fmt.Fprintf(conn, login+"\r"+"\n")
			if err != nil {
				log.Error("error ocurred in sendLogin func", slog.String("err", err.Error()))
				// return
			}
			time.Sleep(10 * time.Second)
		}
	}
}
