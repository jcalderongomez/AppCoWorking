package config

import (
	"backend/models"
	"backend/seeders"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Error al conectar a la base de datos:", err)
	}

	fmt.Println("✅ Conexión exitosa a PostgreSQL")

	// 🛡️ SOLO borra tablas si DB_DROP=true
	if os.Getenv("DB_DROP") == "true" {
		err := db.Migrator().DropTable(
			&models.Role{},
			&models.User{},
			&models.Space{},
			&models.Booking{},
			&models.Payment{},
			&models.ActivityLog{},
		)
		if err != nil {
			log.Fatal("❌ Error al eliminar tablas:", err)
		}
		fmt.Println("⚠️ Tablas eliminadas por configuración DB_DROP=true")
	}

	// 🔄 Migrar modelos
	err = db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Space{},
		&models.Booking{},
		&models.Payment{},
		&models.ActivityLog{},
	)
	if err != nil {
		log.Fatal("❌ Error en migraciones:", err)
	}
	fmt.Println("🚀 Migraciones aplicadas correctamente")

	// Guardar la instancia global
	DB = db

	// Ejecutar seeders (con la instancia db)
	seeders.RunSeeders(db)
}
