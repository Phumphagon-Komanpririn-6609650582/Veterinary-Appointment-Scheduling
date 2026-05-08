package main

import (
	"database/sql"
	"log"
	"veterinary-api/controllers"
	"veterinary-api/middlewares"
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

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authController := controllers.NewAuthController(authService)

	//Vet-API
	vetRepo := repositories.NewVetRepository(db)
	vetService := services.NewVetService(vetRepo)
	vetController := controllers.NewVetController(vetService)

	//Slots-API
	slotRepo := repositories.NewSlotRepository(db)
	slotService := services.NewSlotService(slotRepo)
	slotController := controllers.NewSlotController(slotService)

	//ตั้งค่า Gin Web Server
	r := gin.Default()

	//Routes
<<<<<<< HEAD
	r.POST("/api/login", authController.Login)
	r.GET("/api/vets", middlewares.RequireAuth, vetController.GetAllVets)
	r.GET("/api/vets/:id/slots", middlewares.RequireAuth, slotController.GetAvailableSlots)
=======
	r.GET("/api/vets", vetController.GetAllVets)
	r.GET("/api/vets/:id/slots", slotController.GetAvailableSlots)
	r.GET("/api/slots", slotController.GetAllAvailableSlots)
>>>>>>> Develop

	//รันเซิร์ฟเวอร์ที่พอร์ต 8080
	log.Println("Server is running on http://localhost:8080")
	r.Run(":8080")
}
