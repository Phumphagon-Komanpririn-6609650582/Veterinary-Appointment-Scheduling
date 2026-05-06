package main

import (
	"database/sql"
	"log"
	"veterinary-api/controllers"
	"veterinary-api/repositories"
	"veterinary-api/services"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3" // Driver สำหรับต่อ SQLite
)

func main() {
	//เชื่อมต่อ SQLite
	db, err := sql.Open("sqlite3", "veterinary.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	//ประกอบร่าง API
	vetRepo := repositories.NewVetRepository(db)
	vetService := services.NewVetService(vetRepo)
	vetController := controllers.NewVetController(vetService)

	//ตั้งค่า Gin Web Server
	r := gin.Default()

	//Routes
	r.GET("/api/vets", vetController.GetAllVets)

	//รันเซิร์ฟเวอร์ที่พอร์ต 8080
	log.Println("Server is running on http://localhost:8080")
	r.Run(":8080")
}
