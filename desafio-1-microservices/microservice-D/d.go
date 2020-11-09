package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func main() {
	fmt.Println("init server: port: 9093")
	http.HandleFunc("/", Home)
	http.ListenAndServe(":9093", nil)
}

func Home(w http.ResponseWriter, r *http.Request) {

	type valid struct {
		Status   string
		Discount string
	}

	coupon := r.PostFormValue("coupon")
	isValid := r.PostFormValue("isValid")

	if isValid == "valid" {

		var re = regexp.MustCompile(`/\D/g`)
		Discount := re.ReplaceAllString(coupon, "")

		jsonData, err := json.Marshal(valid{Status: isValid, Discount: Discount})

		if err != nil {
			log.Fatal("Error processing json")
		}
		w.Write(jsonData)
	}

	jsonData, err := json.Marshal(valid{Status: isValid, Discount: ""})

	if err != nil {
		log.Fatal("Error processing json")
	}
	w.Write(jsonData)
}
