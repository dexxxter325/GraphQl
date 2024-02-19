package graph

import (
	"GRAPHQL/graph/model"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"time"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *pgxpool.Pool
}

var carPublishedChannel map[string]chan *model.Car

func init() { //инициализируем пустую мапу для дальнейшего использования
	carPublishedChannel = make(map[string]chan *model.Car)
}

func NewResolver(DB *pgxpool.Pool) *Resolver {
	return &Resolver{DB: DB}
}

func generateHashPassword(password string) string {
	hash := sha1.New()                                                                      // Создаем новый хеш(кодировка) SHA1
	hash.Write([]byte(password))                                                            // Записываем в него байтовое представление пароля
	return fmt.Sprintf("%x", hash.Sum([]byte("8934gfhjgj389gwjuf9pFS0f89asujf903ghm39gk"))) // добавляем соль в хэш-пароль и возвращаем его
}

func generateSessionID() (string, error) {
	randArray := make([]byte, 16)  //срез из 16 эл-тов
	_, err := rand.Read(randArray) //заполняем срез рандомными данными
	if err != nil {
		return "", fmt.Errorf("can`t fill array in generateSession:%s", err)
	}
	sessionID := hex.EncodeToString(randArray) //преобразование среза байтов в строку шестнадцатеричного представления
	return sessionID, nil
}

func (r *Resolver) ValidateSession(sessionID string) (bool, error) {
	var expiresAt time.Time
	query := `select expires_at from sessions where session_id=$1`
	row := r.DB.QueryRow(context.Background(), query, sessionID)
	if err := row.Scan(&expiresAt); err != nil {
		return false, fmt.Errorf("cookie doesn't exist,u must login!:%s", err)
	}
	return true, nil
}

func (r *Resolver) saveSessionToDB(sessionID, userID string) error {
	userIDint, _ := strconv.Atoi(userID)
	currentTime := time.Now()
	expiresAt := currentTime.Add(time.Hour * 24 * 30)
	query := `insert into sessions (session_id,userid,created_at,expires_at) values ($1,$2,$3,$4)`
	_, err := r.DB.Exec(context.Background(), query, sessionID, userIDint, currentTime, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save session:%s", err)
	}
	return nil
}

func (r *Resolver) ValidateCredentials(username, password string) (*model.User, error) {
	var user model.User
	password = generateHashPassword(password)
	query := `select * from users where username=$1 and password=$2`
	ValidateQuery := r.DB.QueryRow(context.Background(), query, username, password)
	if err := ValidateQuery.Scan(&user.ID, &user.Name, &user.Username, &user.Password); err != nil {
		return nil, fmt.Errorf("validate failed:%s", err)
	}
	return &user, nil
}

/*func (r *Resolver) getUserIDFromSessionID(sessionID string) (string, error) {
	var userId string
	query := `SELECT userid FROM sessions WHERE session_id = $1`
	err := r.DB.QueryRow(context.Background(), query, sessionID).Scan(&userId)
	if err != nil {
		return "", fmt.Errorf("failed to get user ID from session ID: %v", err)
	}
	return userId, nil
}*/

func (r *Resolver) getSessionByUserId(id string) (string, error) {
	var session string
	query := `select session_id from sessions where userid=$1`
	getQuery := r.DB.QueryRow(context.Background(), query, id)
	if err := getQuery.Scan(&session); err != nil {
		return "", fmt.Errorf("can't scan in getSessionByUserId:%s", err)
	}
	return session, nil
}
