package main

import (
	"database/sql"
	"log"
	"veterinary-api/controllers"
	"veterinary-api/middlewares"
	"veterinary-api/repositories"
	"veterinary-api/services"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// เชื่อมต่อ SQLite
	db, err := sql.Open("sqlite3", "veterinary.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	appointmentRepo := repositories.NewAppointmentRepository(db)
	appointmentService := services.NewAppointmentService(appointmentRepo)
	appointmentController := controllers.NewAppointmentController(appointmentService)

	authRepo := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepo)
	authController := controllers.NewAuthController(authService)

	// Vet API
	vetRepo := repositories.NewVetRepository(db)
	vetService := services.NewVetService(vetRepo)
	vetController := controllers.NewVetController(vetService)

	// Slots API
	slotRepo := repositories.NewSlotRepository(db)
	slotService := services.NewSlotService(slotRepo)
	slotController := controllers.NewSlotController(slotService)

	// ตั้งค่า Gin
	r := gin.Default()

	r.StaticFile("/", "./public/index.html")
	r.StaticFile("/index.html", "./public/index.html")

	// ==========================================
	// Routes (API)
	// ==========================================
	r.POST("/api/auth/login", authController.Login)

	r.POST("/api/auth/logout",
		middlewares.RequireAuth,
		authController.Logout,
	)

	r.GET("/api/vets",
		middlewares.RequireAuth,
		vetController.GetAllVets,
	)

	r.GET("/api/vets/:id/slots",
		middlewares.RequireAuth,
		slotController.GetAvailableSlots,
	)

	r.POST("/api/vets/:id/slots",
		middlewares.RequireAuth,
		slotController.AddSlot,
	)

	r.GET("/api/slots",
		middlewares.RequireAuth,
		slotController.GetAllAvailableSlots,
	)

	r.POST("/api/appointments",
		middlewares.RequireAuth,
		appointmentController.CreateAppointment,
	)

	r.PUT("/api/appointments/:id",
		middlewares.RequireAuth,
		appointmentController.UpdateAppointment,
	)

	r.DELETE("/api/appointments/:id",
		middlewares.RequireAuth,
		appointmentController.CancelAppointment,
	)

	r.PATCH("/api/appointments/:id/status",
		middlewares.RequireAuth,
		appointmentController.UpdateStatus,
	)

	// รันเซิร์ฟเวอร์
	log.Println("Server is running on http://localhost:8080")
	r.Run(":8080")
}