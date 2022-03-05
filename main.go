package main

import (
	"project-intern/config"
	"project-intern/controller"
	"project-intern/middleware"
	"project-intern/repository"
	"project-intern/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db              *gorm.DB                          = config.SetupDatabaseConnection()
	userRepository  repository.UserRepository         = repository.NewUserRepository(db)
	verifRepository repository.VerificationRepository = repository.NewVerifRepository(db)
	jwtService      service.JWTService                = service.NewJWTService()
	userService     service.UserService               = service.NewUserService(userRepository)
	authService     service.AuthService               = service.NewAuthService(userRepository)
	verifService    service.VerificationService       = service.NewVerifService(verifRepository)
	authController  controller.AuthController         = controller.NewAuthController(authService, jwtService, verifService)
	userController  controller.UserController         = controller.NewUserController(userService, jwtService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	router := gin.Default()
	router.Use(cors.Default())

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/verification", authController.Verification)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("api/user", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("/profile", userController.Profile)
		userRoutes.PUT("/profile", userController.Update)
	}

	r.Run()

}
