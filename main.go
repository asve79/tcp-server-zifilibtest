package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {

	//Берем параметры из коммандной строки
	portNumPtr := flag.Int("port", 3333, "port number")
	delayTime := flag.Int("delay", 2, "delay as seconds")
	minBlockSizePrt := flag.Int("minblocksize", 0, "genetated block size. minimum.")
	maxBlockSizePrt := flag.Int("maxblocksize", 0, "genetated block size. maximum.")
	flagOnlySend := flag.Bool("onlysenddata", false, "Send data only mode")
	flagOnlyReceve := flag.Bool("onlyrecevedata", false, "Receve data only mode")
	flagRandomDataSend := flag.Bool("randomdatasend", false, "Generate random data to send")
	flagRandomDisconnection := flag.Bool("randomdisconnection", false, "Emulate lost connection")
	hostPtr := flag.String("host", "localhost", "listen address")
	connTypePtr := flag.String("type", "tcp", "type tcp/udp")

	flag.Parse()

	//Предосторочности
	if *minBlockSizePrt > *maxBlockSizePrt {
		fmt.Println("Check you paramaters! minblocksize > maxblocksize")
		os.Exit(1)
	}

	if *flagOnlySend && *flagOnlyReceve {
		fmt.Println("You cannot set both receve and send modes at one time")
		os.Exit(1)
	}

	fmt.Printf("Receve only mode : %t\nSend only mode : %t\n", *flagOnlyReceve, *flagOnlySend)
	fmt.Printf("Random lost connection emulate : %t\n", *flagRandomDisconnection)
	fmt.Printf("Random data generator : %t\n", *flagRandomDataSend)

	fmt.Print("Generate data length : ")
	if *minBlockSizePrt == *minBlockSizePrt {
		fmt.Println(strconv.Itoa(*minBlockSizePrt))
	} else {
		fmt.Println("from " + strconv.Itoa(*minBlockSizePrt) + " to " + strconv.Itoa(*maxBlockSizePrt))
	}

	var id int
	id = 1
	l, err := net.Listen(*connTypePtr, *hostPtr+":"+strconv.Itoa(*portNumPtr))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	//Вечный цикл ожидания
	fmt.Println("Listening on " + *hostPtr + ":" + strconv.Itoa(*portNumPtr))
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Обработчик соединения
		go handleRequest(conn, id, *delayTime, *minBlockSizePrt, *maxBlockSizePrt,
			*flagOnlySend, *flagOnlyReceve, *flagRandomDataSend, *flagRandomDisconnection)
		id++

	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, id int, delayTime int, minBlockSizePrt int,
	maxBlockSizePrt int, onlySend bool, onlyReceve bool, flagRandomDataSend bool,
	flagRandomDisconnection bool) {

	buf := make([]byte, 16384)

	var arrsize int
	count := 1
	for {
		if !onlySend { //Если не режим "только отправлять""

			//Ставим таймаут чтобы бесконечно не хдать прихода данных
			conn.SetReadDeadline(time.Now().Add(time.Duration(delayTime) * time.Second))

			nums, err := conn.Read(buf)

			if nums > 0 {
				fmt.Println("Receved:", buf[:nums])
			} else {
				if err == io.EOF { //Если обрыв соединения
					fmt.Println("Lose connection " + strconv.Itoa(id))
					conn.Close()
					break
				} else { //А это если просто ниченго не пришло
					fmt.Println("Nothing receved")
				}
			}
		}

		if !onlyReceve { //Если не режим "только отправка"
			if maxBlockSizePrt == minBlockSizePrt {
				arrsize = maxBlockSizePrt //Если размер блока фиксированой длины
			} else {
				arrsize = rand.Intn(maxBlockSizePrt-minBlockSizePrt) + minBlockSizePrt //Если размер блока рандомной длины
			}

			fmt.Println("Send values: Connection " + strconv.Itoa(id) + ". Send iteration " + strconv.Itoa(count) + " + " +
				strconv.Itoa(arrsize) + " byte(s)")
			_, err := conn.Write([]byte("Connection " + strconv.Itoa(id) + ". Send iteration " + strconv.Itoa(count) + " : "))
			if err != nil {
				fmt.Println("Error send data. Close connection.")
				conn.Close()
				break
			}
			for j := 0; j < arrsize; j++ { //Этим циклом генерируем данные в буфер
				sbyte := []byte("*")
				if flagRandomDataSend {
					rand.Read(sbyte)
				}
				_, err := conn.Write([]byte(sbyte))
				if err != nil {
					fmt.Println("Error send data. Close connection.")
					conn.Close()
					break
				}
			}
		}

		time.Sleep(time.Duration(delayTime) * time.Second) //Ждем сколько-то секунд
		count++

		//ОБрываем соединение если включен режим имитации обрава и случайное число из 100 равно 1
		if flagRandomDisconnection {
			if rand.Intn(100) == 1 {
				conn.Close()
			}
		}
	}
}
