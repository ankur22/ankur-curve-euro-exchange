package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ankur22/ankur-curve-euro-exchange/internal/dao"
)

func main() {
	fmt.Println("Hello World")

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	dao := dao.CreateNewFerAPI("https://api.exchangeratesapi.io/latest", client)
}
