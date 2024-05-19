package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	model "genesis_test_task/internal/app/model"
	repo "genesis_test_task/internal/app/repo"
	service "genesis_test_task/internal/app/service"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/robfig/cron"
)

func Run() {
	// load configs and env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// load smtp user for notification service
	smptUser := model.Email(os.Getenv("SMTP_USER"))
	err = smptUser.Validate()
	if err != nil {
		log.Fatal(err)
	}

	// Load postgres
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/migrations/",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	_ = m
	log.Println("Migrations started")
	err = m.Up()
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	log.Println("Migrations ended")

	// create repo
	subRepo := repo.NewPostgresSubscriptionRepo(db,
		log.New(os.Stdout, "postgres: ", log.LstdFlags))

	exrRepo := repo.NewExchangeRateRepo(os.Getenv("EXCHANGE_RATE_API_KEY"),
		log.New(os.Stdout, "exchange_rate: ", log.LstdFlags))
	notServ := repo.NewGmailNotificationRepo(
		smptUser,
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		log.New(os.Stdout, "notification: ", log.LstdFlags))
	exRate, err := exrRepo.GetExchangeRate(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// create service
	subServ := service.NewSubscriptionService(subRepo, notServ, exrRepo)
	fmt.Println("Current exchange rate: ", exRate)
	// create controller handlers
	// exchange rate handler
	exRateHandler := func(w http.ResponseWriter, r *http.Request) {
		// Get the current exchange rate
		// TODO create service method for exchange rate for provision instead of repo
		rate, err := exrRepo.GetExchangeRate(context.Background())
		if err != nil {
			http.Error(w, "Failed to get exchange rate", http.StatusInternalServerError)
			return
		}
		// Write the exchange rate to the response
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", fmt.Sprint(rate.Value))
	}
	subscribeHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		rawEmail := r.Form.Get("email")
		email := model.Email(rawEmail)
		err = email.Validate()
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid email: %v", err), http.StatusBadRequest)
			return
		}
		found, err := subServ.CheckSubscriberExists(context.Background(), email)
		if err != nil {
			http.Error(w, "Failed to check subscriber", http.StatusInternalServerError)
			return
		}
		if found {
			http.Error(w, "Повертати, якщо e-mail вже є в базі даних", http.StatusConflict)
			return
		}

		err = subServ.Subscribe(context.Background(), email)
		if err != nil {
			//TODO fix error handling for subscriber
			log.Println(err)
			http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
			return
		}
		fmt.Println("E-mail додано")
	}

	// handlers
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/rate", exRateHandler)

	// CRON
	c := cron.New()
	// TODO Cron config (every day at 12:00:00 PM)
	c.AddFunc("0 0 12 * * *", func() {
		// c.AddFunc("* * * * * *", func() {
		// fmt.Println("Cron job started")
		err = subServ.NotifyAll(context.Background())
		if err != nil {
			log.Println(err)
		}
	})
	c.Start()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
