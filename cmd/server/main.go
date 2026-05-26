package main

import (
	"cryptox/cmd/seed"
	centralroutes "cryptox/packages/central_routes"
	"cryptox/packages/cloudinary"
	"cryptox/packages/config"
	"cryptox/packages/database"
	redisClient "cryptox/packages/redis"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	//Load config
	cfg := config.LeadConfig()

	//cloudinary
	cloudinary.InitCloudinary()

	//Connect postgres
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatal("Postgres connection failed:", err)
	}

	//add admin
	seed.Addadmin(db)

	// Connect Redis
	rdb, err := redisClient.NewRedisClient(cfg)
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	//fiber engine
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, //set the limit to 20mb for image uploading
	})

	//for frontend connect
	app.Use(cors.New(cors.Config{
    AllowOrigins: "http://localhost:5173, https://crypto-project-frontend-gold.vercel.app, https://crypto-project-frontend-git-main-tibin-peters-projects.vercel.app, https://crypto-project-frontend-ckw39wv16-tibin-peters-projects.vercel.app",
    AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
    AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    AllowCredentials: true,
}))

	//setup routes
	centralroutes.SetUp(app, db, rdb, cfg.JWTSecret, cfg.RazorpayKey, cfg.RazorpaySecret)

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)

	log.Fatal(app.Listen(":" + port))
}
