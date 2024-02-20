package service

import (
	"GRAPHQL/graph/model"
	"context"
	"fmt"
)

func (p *postgresService) AddCar(ctx context.Context, input model.CarInput) (*model.Car, error) {
	tx, err := p.db.Begin(ctx)
	var car model.Car
	createQuery := `INSERT INTO cars (userId,brand, model, year, price, mileage, description)
											VALUES ($1, $2, $3, $4, $5, $6,$7)
											RETURNING id, brand, model, year, price, mileage, description`
	query := tx.QueryRow(context.Background(), createQuery, input.UserID, input.Brand, input.Model, input.Year, input.Price, input.Mileage, input.Description)
	if err := query.Scan(&car.ID, &car.Brand, &car.Model, &car.Year, &car.Price, &car.Mileage, &car.Description); err != nil {
		tx.Rollback(ctx)
		return &car, fmt.Errorf("scan failed in create:%s", err)
	}
	user, err := p.GetUserByID(ctx, input.UserID)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	// Устанавливаем пользователя для добавленной машины
	car.User = user

	for _, observer := range carPublishedChannel { //перебираем наши каналы
		observer <- &car //передаем созданную машину во все каналы
	}

	if err = tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	return &car, nil
}

// UpdateCar is the resolver for the updateCar field.
func (p *postgresService) UpdateCar(ctx context.Context, id string, input model.CarInput) (*model.Car, error) {
	var car model.Car
	updateQuery := `UPDATE cars SET brand=$1, model=$2, year=$3, price=$4, mileage=$5, description=$6 WHERE id=$7 RETURNING id, brand, model, year, price, mileage, description`
	query := p.db.QueryRow(context.Background(), updateQuery, input.Brand, input.Model, input.Year, input.Price, input.Mileage, input.Description, id)
	if err := query.Scan(&car.ID, &car.Brand, &car.Model, &car.Year, &car.Price, &car.Mileage, &car.Description); err != nil {
		return &car, fmt.Errorf("scan failed in update:%s", err)
	}
	return &car, nil
}

// DeleteCar is the resolver for the deleteCar field.
func (p *postgresService) DeleteCar(ctx context.Context, id string) (*bool, error) {
	fail := false
	success := true
	deleteQuery := `DELETE FROM cars where id=$1`
	res, err := p.db.Exec(context.Background(), deleteQuery, id)
	if err != nil {
		return nil, fmt.Errorf("exec failed in delete:%s", err)
	}
	if res.RowsAffected() <= 0 {
		return &fail, fmt.Errorf("this object already deleted or doesn't exist")
	}
	return &success, nil
}

// GetAllCars is the resolver for the getAllCars field.
func (p *postgresService) GetAllCars(ctx context.Context) ([]*model.Car, error) {
	var cars []*model.Car
	getQuery := `SELECT cars.id, cars.brand, cars.model, cars.year, cars.price, cars.mileage, cars.description,
			   users.id, users.name, users.username,users.password FROM cars
				INNER JOIN users ON cars.userid = users.id ORDER BY cars.id ASC`
	query, err := p.db.Query(context.Background(), getQuery)
	if err != nil {
		return nil, fmt.Errorf("failed query in getall:%s", err)
	}
	defer query.Close()
	for query.Next() { //читает каждую строку из результата sql query ,чтобы затем добавить в срез
		var car model.Car
		var user model.User
		if err := query.Scan(&car.ID, &car.Brand, &car.Model, &car.Year, &car.Price, &car.Mileage, &car.Description,
			&user.ID, &user.Name, &user.Username, &user.Password); err != nil {
			return nil, fmt.Errorf("scan failed in getall:%s", err)
		}
		car.User = &user
		cars = append(cars, &car) //с каждой успешно прочитанной строкой,добавляем эл-т в срез
	}
	return cars, nil
}

// GetCarByID is the resolver for the getCarById field.
func (p *postgresService) GetCarByID(ctx context.Context, id string) (*model.Car, error) {
	var car model.Car
	var user model.User
	getByIdQuery := `SELECT cars.id, cars.brand, cars.model, cars.year, cars.price, cars.mileage, cars.description,
			   users.id, users.name, users.username,users.password FROM cars
				INNER JOIN users ON cars.userid = users.id where users.id=$1`
	query := p.db.QueryRow(context.Background(), getByIdQuery, id)
	if err := query.Scan(&car.ID, &car.Brand, &car.Model, &car.Year, &car.Price, &car.Mileage, &car.Description,
		&user.ID, &user.Name, &user.Username, &user.Password); err != nil {
		return nil, fmt.Errorf("scan failed in getbyID:%s", err)
	}
	car.User = &user
	return &car, nil
}

// GetUserByID is the resolver for the getUsers field.
func (p *postgresService) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	getUserQuery := `SELECT id, name, username,password FROM users WHERE id = $1`
	err := p.db.QueryRow(ctx, getUserQuery, userID).Scan(&user.ID, &user.Name, &user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
