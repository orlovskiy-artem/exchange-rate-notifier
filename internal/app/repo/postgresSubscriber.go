package repo

// import (
// 	"context"
// 	"database/sql"
// 	"log"

// 	model "genesis_test_task/internal/app/model"
// )

// type PostgresSubscriberRepo struct {
// 	db  *sql.DB
// 	log *log.Logger
// 	// config *Config
// }

// func NewPostgresSubscriberRepo(
// 	db *sql.DB,
// 	log *log.Logger,
// ) *PostgresSubscriptionRepo {
// 	return &PostgresSubscriptionRepo{db, log}
// }

// func (r PostgresSubscriptionRepo) createSubscriberTable(ctx context.Context) error {
// 	query := `CREATE TABLE IF NOT EXISTS subscribers (
// 		id SERIAL PRIMARY KEY,
// 		email TEXT NOT NULL,
// 		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// 	)`
// 	_, err := r.db.ExecContext(ctx, query)
// 	if err != nil {
// 		log.Fatal(err)
// 		return err
// 	}
// 	return nil
// }

// // Save implements ISubscriberRepo.
// func (r *PostgresSubscriptionRepo) SaveSubscriber(ctx context.Context, s model.Subscription) error {
// 	query := `INSERT INTO subscribers (email) VALUES ($1)`
// 	_, err := r.db.ExecContext(ctx, query, s.Email)
// 	return err
// }

// func (r *PostgresSubscriptionRepo) FindSubscriber(ctx context.Context, e model.Email) (*model.Subscription, bool, error) {
// 	query := `SELECT id, email FROM subscribers WHERE email = $1`
// 	row := r.db.QueryRowContext(ctx, query, e)
// 	var id int
// 	var dbEmail string
// 	err := row.Scan(&id, &dbEmail)
// 	if err == sql.ErrNoRows {
// 		return nil, false, nil
// 	}
// 	if err != nil {
// 		return nil, false, err
// 	}
// 	return &model.Subscription{ID: id, Email: e}, true, nil
// }

// func (r *PostgresSubscriptionRepo) ForEachSubscription(ctx context.Context, hf model.SubscriptionHandle) (err error) {
// 	query := `SELECT id, email FROM subscribers`
// 	rows, err := r.db.QueryContext(ctx, query)
// 	if err != nil {
// 		return err
// 	}
// 	defer rows.Close()

// 	// var subscribers []model.Subscription
// 	for rows.Next() {
// 		var id int
// 		var email string
// 		err = rows.Scan(&id, &email)
// 		if err != nil {
// 			return err
// 		}
// 		subscription := model.Subscription{
// 			ID:    id,
// 			Email: model.Email(email),
// 		}
// 		err = hf(ctx, subscription)
// 		if err != nil {
// 			log.Print(err)
// 		}
// 	}
// 	if err = rows.Err(); err != nil {
// 		return err
// 	}
// 	return nil
// }
