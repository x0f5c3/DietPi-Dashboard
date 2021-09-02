package lib

import (
	"log"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type system struct {
	CPU  float64 `json:"cpu"`
	RAM  MemData `json:"ram"`
	Swap MemData `json:"swap"`
}

type processlist struct {
	Processes []ProcessData `json:"processes"`
}

type softwarelist struct {
	Software []DPSoftwareData `json:"software"`
	Response string           `json:"response,omitempty"`
}

type request struct {
	Page string      `json:"page"`
	Do   string      `json:"do"`
	Args interface{} `json:"args"`
}

func readSocket(c *websocket.Conn, m chan<- request, n chan<- struct{}) {
	var req request
	firstmessage := true
	for {
		err := c.ReadJSON(&req)
		if err != nil {
			log.Println("Couldn't get data from frontend:", err)
			close(m)
			close(n)
			break
		}
		if req.Do == "" {
			if !firstmessage {
				n <- struct{}{}
			} else {
				firstmessage = false
			}
		}
		m <- req
	}
}

func ServeWebsockets(w http.ResponseWriter, r *http.Request) {
	log.Println("Request to /ws (websockets)")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Couldn't upgrade connection to websockets:", err)
		return
	}
	m := make(chan request)
	n := make(chan struct{})
	go readSocket(c, m, n)
	defer c.Close()
	for i := range m {
		switch i.Page {
		case "/":
		main:
			for {
				stats := system{CPU(), RAM(), Swap()}
				err := c.WriteJSON(stats)
				if err != nil {
					log.Println("Couldn't send message to frontend:", err)
				}
				select {
				case <-n:
					break main
				default:
				}
			}
		case "/process":
		process:
			for {
				err := c.WriteJSON(processlist{Processes()})
				if err != nil {
					log.Println("Couldn't send message to frontend:", err)
				}
				select {
				case <-n:
					break process
				default:
				}
			}
		case "/software":
			err := c.WriteJSON(softwarelist{DPSoftware(), ""})
			if err != nil {
				log.Println("Couldn't send message to frontend:", err)
			}
		software:
			for {
				select {
				case data := <-m:
					argArr := []string{"/boot/dietpi/dietpi-software", data.Do}
					for _, element := range data.Args.([]interface{}) {
						argArr = append(argArr, strconv.Itoa(int(element.(float64))))
					}
					cmd := &exec.Cmd{
						Path: "/boot/dietpi/dietpi-software",
						Args: argArr,
					}
					out, _ := cmd.Output()
					err := c.WriteJSON(softwarelist{DPSoftware(), string(out)})
					if err != nil {
						log.Println("Couldn't send message to frontend:", err)
					}
				case <-n:
					break software
				}
			}
		}
	}
}