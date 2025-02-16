// backend/cmd/server/main.go
package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"yopta-template/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	authMiddleware "yopta-template/internal/middleware"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки переменных окружения: %v", err)
	}

	dsn := os.Getenv("DSN")
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiry := os.Getenv("JWT_EXPIRY")
	if dsn == "" || jwtSecret == "" || jwtExpiry == "" {
		log.Fatal(
			"Не все необходимые переменные окружения установлены (DSN, JWT_SECRET, JWT_EXPIRY)",
		)
	}

	// Подключение к базе данных
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Проверка подключения к базе данных
	if err := db.Ping(); err != nil {
		log.Fatalf("Ошибка проверки соединения с базой данных: %v", err)
	}

	// Создание нового роутера
	r := chi.NewRouter()

	// Подключение middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders: []string{"Link"},
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Обслуживание статических файлов (фронтенд)
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/*", fs)

	// Публичные маршруты: регистрация и логин
	r.Route("/api", func(r chi.Router) {
		r.Post("/register", handlers.Register(db, jwtSecret, jwtExpiry))
		r.Post("/login", handlers.Login(db, jwtSecret, jwtExpiry))
		r.Post("/refresh-token", handlers.RefreshToken(jwtSecret, jwtExpiry))
	})

	// Защищённые маршруты: требуется JWT-аутентификация
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.JWTMiddleware(jwtSecret))
		r.Get("/api/profile", handlers.GetProfile(db))
		r.Post("/api/change-password", handlers.ChangePassword(db))
	})

	// Запуск сервера с graceful shutdown
	srv := &http.Server{
		Addr:    ":6033",
		Handler: r,
	}

	go func() {
		log.Printf("Сервер запущен на %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Ожидание сигнала прерывания для корректного завершения работы
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Сервер выключается...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении работы сервера: %v", err)
	}
	log.Println("Сервер успешно завершил работу")
}
