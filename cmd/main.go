package main

import (
	"context"
	"fmt"
	"log"
	httpNet "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/throindev/payments/cmd/config"
	"github.com/throindev/payments/cmd/http"
	"github.com/throindev/payments/internal/infra/mercadopago"
	"github.com/throindev/payments/internal/infra/mongodb"
	"github.com/throindev/payments/internal/usecases"
)

func main() {
	config.Load()
	mongoClient := mongodb.NewMongoClient(config.AppConfig.MongoURI, config.AppConfig.DbName)

	r := gin.Default()

	// Middlewares
	r.Use(gin.Recovery())
	// r.Use(middleware.AuthMiddleware())

	// payment providers
	mercadopagoProvider := mercadopago.NewClient()

	//repositories mongo
	paymentMongoRepository := mongodb.NewPaymentMongoRepository(mongoClient)
	subscriptionMongoRepository := mongodb.NewSubscriptionMongoRepository(mongoClient)
	planMongoRepository := mongodb.NewPlanMongoRepository(mongoClient)

	//usecases
	planUsecase := usecases.NewPlanUsecases(planMongoRepository)
	subscriptionUsecases := usecases.NewSubscriptionUsecases(subscriptionMongoRepository, planUsecase)
	paymentUsecases := usecases.NewPaymentUsecases(paymentMongoRepository, mercadopagoProvider, subscriptionUsecases, planUsecase)

	//controllers
	paymentController := http.NewPaymentController(paymentUsecases)
	planController := http.NewPlanController(planUsecase)

	api := r.Group("/api")
	{
		api.POST("/payment/mercadopago/callback", paymentController.CallbackfromMercadoPago)
		api.POST("/payment", paymentController.CreatePayment)
		api.POST("/plan", planController.CreatePlan)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "ok"})
	})

	srv := &httpNet.Server{
		Addr:    fmt.Sprintf(":%s", config.AppConfig.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != httpNet.ErrServerClosed {
			log.Fatalf("Error to set up server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Error to shut down server:", err)
	}

	log.Println("Server has been finished")
}
