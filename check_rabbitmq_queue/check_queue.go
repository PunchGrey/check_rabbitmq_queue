package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

//Queue структура описывающая очередь
type Queue struct {
	Messages               int    `json:"messages"`
	MessagesUnacknowledged int    `json:"messages_unacknowledged"`
	Name                   string `json:"name"`
}

//ArrQueue слайс, который парсится из get запроса
type ArrQueue []Queue

var nagiousNotify ArrQueue

func main() {

	var (
		host                       = flag.String("h", "127.0.0.1", "This is the host with rabbitmq")
		port                       = flag.Int("p", 15672, "Rabbitmq is using this port")
		user                       = flag.String("u", "admin", "This is the user for rabbitmq")
		pass                       = flag.String("pw", "1234", "This is the password for rabbitmq")
		nameQueue                  = flag.String("q", "celery, coverage_tasks", "This is the queue name")
		errMessages                = flag.Int("em", 28, "Error messages in the queue of rabbitmq")
		warnMessages               = flag.Int("wm", 15, "Warn messages in the queue of rabbitmq")
		errMessagesUnacknowledged  = flag.Int("emu", 5, "Error messages in the queue of rabbitmq")
		warnMessagesUnacknowledged = flag.Int("wmu", 2, "Warn messages in the queue of rabbitmq")
	)
	flag.Parse()

	url := fmt.Sprintf("http://%s:%s@%s:%d/api/queues", *user, *pass, *host, *port)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		fmt.Println(resp.Status)
		os.Exit(2)
	}

	var result ArrQueue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		fmt.Println(err)
		os.Exit(2)
	}
	resp.Body.Close()

	//	fmt.Println(&result)
	if *nameQueue != "ALL" {
		for _, item := range result {
			arrNamesQue := strings.Split(*nameQueue, ",")
			for _, qname := range arrNamesQue {
				if item.Name == strings.TrimSpace(qname) {
					nagiousNotify = append(nagiousNotify, item)
					break
				}
			}
		}
	} else {
		nagiousNotify = result
	}

	for _, item := range nagiousNotify {
		if item.Messages > *errMessages {
			printNagiosNotify(nagiousNotify)
			os.Exit(2)
		}
		if item.MessagesUnacknowledged > *errMessagesUnacknowledged {
			printNagiosNotify(nagiousNotify)
			os.Exit(2)
		}
		if item.Messages > *warnMessages {
			printNagiosNotify(nagiousNotify)
			os.Exit(1)
		}
		if item.MessagesUnacknowledged > *warnMessagesUnacknowledged {
			printNagiosNotify(nagiousNotify)
			os.Exit(1)
		}
		printNagiosNotify(nagiousNotify)
		os.Exit(0)
	}
}

func printNagiosNotify(nagiousNotify ArrQueue) {
	for _, item := range nagiousNotify {
		fmt.Printf("name queue %s: %d messages,  %d MessagesUnacknowledged \n", item.Name, item.Messages, item.MessagesUnacknowledged)
	}
}
