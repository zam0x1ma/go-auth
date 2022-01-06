package main

import (
    "go-auth/database"
    "go-auth/routes"
)

func main() {
    database.Connect()

    router := routes.Setup()
    
    router.Run(":8000")
}
