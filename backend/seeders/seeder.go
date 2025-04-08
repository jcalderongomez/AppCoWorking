package seeders

import (
	"backend/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunSeeders(db *gorm.DB) {
	var count int64
	db.Model(&models.Role{}).Count(&count)

	if count > 0 {
		log.Println("⚠️ Seeders ya fueron ejecutados anteriormente. No se insertarán datos.")
		return
	}

	log.Println("🚀 Ejecutando seeders...")

	// Crear roles
	roles := []models.Role{
		{Nombre: "Administrador"},
		{Nombre: "Usuario"},
		{Nombre: "Recepcionista"},
	}
	db.Create(&roles)

	// Hashear contraseñas
	adminPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	userPass, _ := bcrypt.GenerateFromPassword([]byte("usuario123"), bcrypt.DefaultCost)

	// Crear usuarios
	users := []models.User{
		{
			Nombre:        "Admin",
			Email:         "admin@example.com",
			Password:      string(adminPass),
			RolID:         1,
			FechaRegistro: time.Now(),
		},
		{
			Nombre:        "Usuario1",
			Email:         "user1@example.com",
			Password:      string(userPass),
			RolID:         2,
			FechaRegistro: time.Now(),
		},
	}
	db.Create(&users)

	// Crear espacios
	spaces := []models.Space{
		{Nombre: "Sala de Reuniones A", Capacidad: 10, Ubicacion: "Piso 1"},
		{Nombre: "Sala de Conferencias B", Capacidad: 50, Ubicacion: "Piso 2"},
	}
	db.Create(&spaces)

	// Crear reservas
	bookings := []models.Booking{
		{
			UsuarioID:   1,
			EspacioID:   1,
			FechaInicio: time.Now().AddDate(0, 0, 1),
			FechaFin:    time.Now().AddDate(0, 0, 1).Add(2 * time.Hour),
			Estado:      "confirmada",
		},
		{
			UsuarioID:   2,
			EspacioID:   2,
			FechaInicio: time.Now().AddDate(0, 0, 2),
			FechaFin:    time.Now().AddDate(0, 0, 2).Add(3 * time.Hour),
			Estado:      "pendiente",
		},
	}
	db.Create(&bookings)

	// Crear pagos
	var reservas []models.Booking
	db.Find(&reservas)
	if len(reservas) >= 2 {
		payments := []models.Payment{
			{
				ReservationID: reservas[0].ID,
				Amount:        100.50,
				PaymentMethod: "Tarjeta",
				Status:        "Pagado",
				CreatedAt:     time.Now(),
			},
			{
				ReservationID: reservas[1].ID,
				Amount:        50.00,
				PaymentMethod: "Efectivo",
				Status:        "Pendiente",
				CreatedAt:     time.Now(),
			},
		}
		db.Create(&payments)
	}

	// Crear logs de actividad
	logs := []models.ActivityLog{
		{
			UserID:  1,
			Accion:  "Inicio de sesión",
			Detalles: "El usuario Admin inició sesión",
		},
		{
			UserID:  2,
			Accion:  "Reserva creada",
			Detalles: "Usuario1 reservó la Sala de Conferencias B",
		},
	}
	db.Create(&logs)

	log.Println("✅ Seeders ejecutados correctamente.")
}
