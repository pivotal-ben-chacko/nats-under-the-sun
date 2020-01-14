package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	nats "github.com/nats-io/nats.go"
)

type NatCommand struct {
	Protocol  int           `json:"protocol"`
	Method    string        `json:"method"`
	Arguments []interface{} `json:"arguments"`
	ReplyTo   string        `json:"reply_to"`
}

func main() {

	server := flag.String("server", "", "NATs server to target")
	agent := flag.String("agent", "", "Agent ID to filter on")
	flag.Parse()

	msgJSON := NatCommand{}

	// Connect to NATS using client credentials
	fmt.Println("Connecting to NATS using client credentials...")
	var nc1, err = nats.Connect(*server+":4222",
		nats.ClientCert("./nats_cert.pem", "./nats_private.key"),
		nats.RootCAs("./nats_ca_cert.pem"))

	if err != nil {
		fmt.Println("Oh boy!")
	} else {
		fmt.Println("We gooood....")
	}

	// Connect to Director NATS
	fmt.Println("Connecting to NATS using director credentials...")
	var nc2, err2 = nats.Connect(*server+":4222",
		nats.ClientCert("./nats_director_cert.pem", "./nats_director_private.key"),
		nats.RootCAs("./nats_director_ca_cert.pem"))

	if err2 != nil {
		fmt.Println("Oh boy!")
	} else {
		fmt.Println("We gooood....")
	}

	var startwg sync.WaitGroup
	startwg.Add(1)

	fmt.Println(*agent)
	if *agent != "" {
		// Subscribe to agent directives from director
		fmt.Println("Subscribing to agent directives from director")
		nc1.Subscribe(*agent, func(m *nats.Msg) {
			fmt.Printf("Received message from Director: %s\n\n", string(m.Data))

			err := json.Unmarshal([]byte(m.Data), &msgJSON)
			if err != nil {
				panic(err)
			}
		})
	} else {
		fmt.Println("Will not subscribe to director messages as agent id was not supplied...")
	}

	// Subscribe to agent messages to director
	fmt.Println("Subscribing to agent messages to director")
	nc2.Subscribe("director.>", func (m *nats.Msg) {
		fmt.Printf("Received message from agent: %s\n\n\n", string(m.Data))
	})

	startwg.Wait()
}




