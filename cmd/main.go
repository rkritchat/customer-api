package main

import (
	"ex_produce/internal/customer"
	"ex_produce/internal/repository"
	"ex_produce/internal/router"
	"fmt"
	"net/http"
)

func main() {
	//init database
	customerRepo := repository.NewCustomer()

	//init service
	service := customer.NewService(customerRepo)

	//init router
	r := router.InitRouter(service)

	//start server
	fmt.Println("start on port 9000")
	_ = http.ListenAndServe(":9000", r)
}
