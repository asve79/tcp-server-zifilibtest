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
		if !onlySend {
			_, err := conn.Read(make([]byte, 0))
			if err != io.EOF {
				fmt.Println("connection closed....", err)
				break
			}

			conn.SetReadDeadline(time.Now().Add(time.Duration(delayTime) * time.Second))
			nums, _ := conn.Read(buf)
			//if err != nil {
			//	fmt.Println("Error reading:", err.Error())
			//	continue
			//}

			if nums > 0 {
				fmt.Println("Receved:", buf[:nums])
			} else {
				fmt.Println("Nothing receved")
			}
		}

		if !onlyReceve {
			if maxBlockSizePrt == minBlockSizePrt {
				arrsize = maxBlockSizePrt
			} else {
				arrsize = rand.Intn(maxBlockSizePrt-minBlockSizePrt) + minBlockSizePrt
			}

			fmt.Println("Send values: Connection " + strconv.Itoa(id) + ". Send iteration " + strconv.Itoa(count) + " + " +
				strconv.Itoa(arrsize) + " byte(s)")
			_, err := conn.Write([]byte("Connection " + strconv.Itoa(id) + ". Send iteration " + strconv.Itoa(count) + " : "))
			if err != nil {
				fmt.Println("Error send data. Close connection.")
				conn.Close()
				break
			}
			for j := 0; j < arrsize; j++ {
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

		time.Sleep(time.Duration(delayTime) * time.Second)
		count++

		//ОБрываем соединение если случайное число из 100 равно 1
		if flagRandomDisconnection {
			if rand.Intn(100) == 1 {
				conn.Close()
			}
		}
	}
}
