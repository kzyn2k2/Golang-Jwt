package main

import (
	"fmt"
	"net/http"

	"github.com/Kzynyi/ass/internal/routes"
	"github.com/Kzynyi/ass/internal/service/user"
)

func main() {
	fmt.Println("Server started on port 8000")
	fmt.Println("Enter postgres username:")
	fmt.Scanln(&user.Username)
	fmt.Println("Enter postgres password:")
	fmt.Scanln(&user.Password)
	http.ListenAndServe(":8000", routes.GetNewMux())

}
