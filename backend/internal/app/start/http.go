package start

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/w0ikid/study-platform/internal/api/routes"
	"gitlab.com/w0ikid/study-platform/internal/app/config"
	"gitlab.com/w0ikid/study-platform/internal/domain/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"




	"github.com/swaggo/gin-swagger"          // gin-swagger middleware
    "github.com/swaggo/files"                // swagger embed files
    _ "gitlab.com/w0ikid/study-platform/docs"                // docs is generated by Swag CLI, you have to import it.
)
func HTTP(cfg *config.Config, userUseCase *usecase.UserUseCase, courseUseCase *usecase.CourseUseCase, lessonUseCase *usecase.LessonUseCase , enrollment *usecase.EnrollmentUseCase, lessonProgressUseCase *usecase.LessonProgressUseCase, certificateUseCase *usecase.CertificateUseCase)  {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"}, // адрес фронта
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	// Swagger UI доступен по /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupRoutes(router, userUseCase, courseUseCase,  lessonUseCase, enrollment, lessonProgressUseCase,certificateUseCase ,cfg)

	
	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:    ":" + cfg.HTTPServer.Port,
		Handler: router,
	}

	// Запускаем сервер в горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ждем сигнала для грациозного завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Даем 5 секунд на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}