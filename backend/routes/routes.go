package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas del backend
func SetupRoutes(r *gin.Engine) {
	// 🔐 Autenticación
	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", controllers.Register)
		authRoutes.POST("/login", controllers.Login)
	}

	// 🧾 Logs de actividad
	logRoutes := r.Group("/api/logs")
	{
		logRoutes.GET("/", controllers.GetLogs)
		// logRoutes.GET("/:id", controllers.GetActivityLogByID)
		logRoutes.POST("/", controllers.CreateLog) // solo si permites logs manuales
		// logRoutes.DELETE("/:id", controllers.DeleteActivityLog)
	}

	// 🧪 Ruta de prueba protegida
	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Servidor activo 🔥"})
	})

	// ✅ Rutas protegidas con autenticación
	protected := r.Group("/api")
	protected.Use(controllers.AuthMiddleware())
	{
		// 👤 Usuarios
		userRoutes := protected.Group("/users")
		{
			userRoutes.GET("/", controllers.GetUsers)
			userRoutes.GET("/:id", controllers.GetUserByID)
			userRoutes.POST("/", controllers.CreateUser)
			userRoutes.PUT("/:id", controllers.UpdateUser)
			userRoutes.DELETE("/:id", controllers.DeleteUser)
		}

		// 🏢 Espacios
		spaceRoutes := protected.Group("/spaces")
		{
			spaceRoutes.GET("/", controllers.GetSpaces)
			spaceRoutes.GET("/:id", controllers.GetSpaceByID)
			spaceRoutes.POST("/", controllers.CreateSpace)
			spaceRoutes.PUT("/:id", controllers.UpdateSpace)
			spaceRoutes.DELETE("/:id", controllers.DeleteSpace)
		}

		// 📅 Reservas
		bookingRoutes := protected.Group("/bookings")
		{
			bookingRoutes.GET("/", controllers.GetBookings)
			bookingRoutes.GET("/:id", controllers.GetBookingByID)
			bookingRoutes.POST("/", controllers.CreateBooking)
			bookingRoutes.PUT("/:id", controllers.UpdateBooking)
			bookingRoutes.DELETE("/:id", controllers.DeleteBooking)
		}

		// 💳 Pagos
		paymentRoutes := protected.Group("/payments")
		{
			paymentRoutes.GET("/", controllers.GetPayments)
			paymentRoutes.GET("/:id", controllers.GetPaymentByID)
			paymentRoutes.POST("/", controllers.CreatePayment)
			paymentRoutes.PUT("/:id", controllers.UpdatePayment)
			paymentRoutes.DELETE("/:id", controllers.DeletePayment)
		}
	}
}
