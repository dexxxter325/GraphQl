package service

import (
	"GRAPHQL/graph/model"
	"context"
	"fmt"
	"net/http"
	"time"
)

type CookieResponseWriter struct { //to use responsewriter in auth func`s
	http.ResponseWriter
}

func (p *postgresService) Register(ctx context.Context, input model.RegisterInput) (*model.User, error) {
	var user model.User
	input.Password = GenerateHashPassword(input.Password)
	createQuery := `insert into users(name,username,password) values ($1,$2,$3) RETURNING id,name,username,password`
	query := p.db.QueryRow(context.Background(), createQuery, input.Name, input.Username, input.Password)
	if err := query.Scan(&user.ID, &user.Name, &user.Username, &user.Password); err != nil {
		return &user, fmt.Errorf("failed scan in registeerr:%s", err)
	}
	//убираем куки прошлого юзера из контекста
	cookieResponse := ctx.Value("cookie").(*CookieResponseWriter)
	cookie := &http.Cookie{
		Name:     "session_id",
		HttpOnly: true,
	}
	http.SetCookie(cookieResponse.ResponseWriter, cookie)
	return &user, nil
}

func (p *postgresService) Login(ctx context.Context, input model.LoginInput) (*model.Login, error) {
	var login model.Login
	user, err := p.ValidateCredentials(input.Username, input.Password)
	if err != nil {
		return &login, fmt.Errorf("validateCredentails in login failed:%s", err)
	}
	existingSessionID, err := p.getSessionByUserId(user.ID)
	if existingSessionID != "" {
		// Если сессия уже существует, возвращаем ее и не создаем новую
		login.SessionID = existingSessionID
		return &login, nil
	}
	sessionId, err := generateSessionID()
	if err != nil {
		return &login, fmt.Errorf("can`t generate sessionID in login:%s", err)
	}
	err = p.saveSessionToDB(sessionId, user.ID)
	if err != nil {
		return &login, fmt.Errorf("fail to saveSessionTODB:%s", err)
	}
	cookieResponse := ctx.Value("cookie").(*CookieResponseWriter)
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Path:     "/",
	}
	http.SetCookie(cookieResponse.ResponseWriter, cookie)
	//fmt.Printf("cookie in login:%s\n", cookie)
	login.SessionID = sessionId
	return &login, nil
}

func (p *postgresService) Logout(ctx context.Context, id string) (bool, error) {
	//убираем куки прошлого юзера из контекста
	cookieResponse := ctx.Value("cookie").(*CookieResponseWriter)
	http.SetCookie(cookieResponse.ResponseWriter, &http.Cookie{
		Name:     "session_id",
		HttpOnly: true,
	})
	sessionsQuery := `delete from sessions where userid=$1`
	res, err := p.db.Exec(context.Background(), sessionsQuery, id)
	if err != nil {
		return false, fmt.Errorf("can`t exec in logout:%s", err)
	}
	if res.RowsAffected() == 0 {
		return false, fmt.Errorf("this object already deleted or doesn't exist")
	}
	return true, nil
}
