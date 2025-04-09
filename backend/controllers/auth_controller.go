package controllers

import (
	"backend/config"
	"backend/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Estructura para recibir credenciales de inicio de sesión
type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Generar JWT
func GenerateToken(user models.User) (string, error) {
	secretKey := []byte(config.JWTSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(secretKey)
}

// Registrar usuario
func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar la contraseña"})
		return
	}

	user.Password = string(hashedPassword)
	config.DB.Create(&user)

	// Crear log de actividad
	CreateActivityLog(user.ID, "registro", "Usuario registrado")

	c.JSON(http.StatusOK, gin.H{"message": "Usuario registrado con éxito"})
}

// Iniciar sesión
func Login(c *gin.Context) {
	var input LoginInput
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar usuario por email
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// Comparar contraseñas
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// Generar token JWT
	token, err := GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
		return
	}

	// Crear log de actividad
	CreateActivityLog(user.ID, "inicio de sesión", "Usuario inició sesión")

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Middleware para validar token
// Middleware para validar token y extraer datos del usuario
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
			c.Abort()
			return
		}

		// Validar que el token tenga formato Bearer <token>
		var tokenString string
		fmt.Sscanf(authHeader, "Bearer %s", &tokenString)

		// Parsear el token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Extraer claims (info del token) y poner user_id en el contexto
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, ok := claims["user_id"]; ok {
				c.Set("user_id", userID)
			}
		}

		c.Next()
	}

}
func GetMyProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, user)
}


// Actualizar perfil del usuario logueado
func UpdateMyProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input struct {
		Nombre        string `json:"nombre"`
		Password      string `json:"password"`
		NuevaPassword string `json:"nueva_password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Validar la contraseña actual si quiere cambiarla
	if input.NuevaPassword != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Contraseña actual incorrecta"})
			return
		}
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(input.NuevaPassword), bcrypt.DefaultCost)
		user.Password = string(hashedPass)
	}

	// Actualizar el nombre si se proporciona
	if input.Nombre != "" {
		user.Nombre = input.Nombre
	}

	// Guardar cambios
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el perfil"})
		return
	}

	// Log de actividad
	CreateActivityLog(user.ID, "actualización de perfil", "Usuario actualizó su perfil")

	c.JSON(http.StatusOK, gin.H{"message": "Perfil actualizado con éxito"})
}
