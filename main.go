package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

// Poziv:	telnet localhost 8888

func main() {

	//  Provera pitanja
	questions := get_questions(15)

	points := 0
	for _, q := range questions {
		println(q.question_string())
		var resp string
		fmt.Scan(&resp)
		if resp == q.Correct_answer {
			println("Correct")
			points += q.Points
		} else {
			println("Incorrect, correct is " + q.Correct_answer)
		}
		println("Points = " + strconv.Itoa(points))
	}

	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("server started on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		go s.newClient(conn)
	}
}
