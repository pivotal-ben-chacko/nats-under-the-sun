package main

import (
  "encoding/json"
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
	msgJSON := NatCommand{}

// Connect to NATS using client credentials
fmt.Println("Connecting to NATS using client credentials...")
var nc1, err = nats.Connect("nats://10.193.69.11:4222",
			nats.ClientCert("./nats_cert.pem", "./nats_private.key"),
			nats.RootCAs("./nats_ca_cert.pem"))

if err != nil {
    fmt.Println("Oh boy!")
} else {
    fmt.Println("We gooood....")
}

// Connect to Director NATS
fmt.Println("Connecting to NATS using director credentials...")
var nc2, err2 = nats.Connect("nats://10.193.69.11:4222",
                        nats.ClientCert("./nats_director_cert.pem", "./nats_director_private.key"),
                        nats.RootCAs("./nats_director_ca_cert.pem"))

if err2 != nil {
    fmt.Println("Oh boy!")
} else {
    fmt.Println("We gooood....")
}


var startwg sync.WaitGroup
startwg.Add(1)
// Subscribe to agent directives from director
fmt.Println("Subscribing to agent directives from director")
nc1.Subscribe("agent.24fd2d5b-afc4-43b6-8554-d95d344e81e7", func(m *nats.Msg) {
	fmt.Printf("Received message from Director: %s\n\n", string(m.Data))

	err := json.Unmarshal([]byte(m.Data), &msgJSON)
	if err != nil {
		panic(err)
	}
});

// Subscribe to agent messages to director
fmt.Println("Subscribing to agent messages to director")
nc2.Subscribe("director.>", func(m *nats.Msg) {
	fmt.Printf("Received message from agent: %s\n\n\n", string(m.Data))
});
startwg.Wait()

}
