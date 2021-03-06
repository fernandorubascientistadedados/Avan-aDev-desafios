package main

import (
	"net/http"
	"net/url"
	"log"
	"io/ioutil"
	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"github.com/wesleywillians/go-rabbitmq/queue"
	uuid "github.com/satori/go.uuid"
)

type Order struct {
	ID uuid.UUID
	Coupon   string
	CcNumber string
}

type Result struct {
	Status string
}

const (
	InvalidCoupon = "invalid"
	ValidCoupon = "valid"
	ConnectionError = "connection error"
)

func NewOrder() Order {
	return Order{ID: uuid.NewV4()}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}

func main()  {
	messageChannel := make(chan amqp.Delivery)

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	for msg := range messageChannel {
		Process(msg)
	}
}

func Process(msg amqp.Delivery) {
	order := NewOrder()

	json.Unmarshal(msg.Body, &order)

	resultCoupon := MakeHttpCall("http://localhost:9092", order.Coupon)

	switch resultCoupon.Status {
	case InvalidCoupon:
		log.Println("Order:", order.ID,"- invalid coupon! ❌")	
	case ConnectionError:
		msg.Reject(false)
		log.Println("Order:", order.ID,"- connection error! ⏳")
	case ValidCoupon:
		log.Println("Order:", order.ID,"- processed! ✔️")
	}

}

func MakeHttpCall(urlMicroservice string, coupon string) Result {

	values := url.Values{}
	values.Add("coupon", coupon);

	res, err := http.PostForm(urlMicroservice, values)

	if err != nil {
		result := Result{Status: ConnectionError}

		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result
}