package graph

import (
	"GRAPHQL/response"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type CookieResponseWriter struct { //to use responsewriter in func Login
	http.ResponseWriter
}

func (r *Resolver) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookieResponse := CookieResponseWriter{
			ResponseWriter: w,
		}
		ctx := context.WithValue(req.Context(), "cookie", &cookieResponse)
		cookie, err := req.Cookie("session_id")
		if err != nil {
			response.ErrHandler(w, fmt.Errorf("session with this id doesn't exist:%s", err.Error()), http.StatusUnauthorized)
			return
		}
		fmt.Printf("cookie in middleware:%s\nexpiresat:%v\n", cookie.Value, cookie.MaxAge)

		operationName, err := extractOperationName(req)
		if err != nil {
			response.ErrHandler(w, fmt.Errorf("extractOperationName failed in middleware:%s", err.Error()), http.StatusUnauthorized)
			return
		}
		if requiresAuthentication(operationName) {
			// Получаем идентификатор сессии из куки.
			sessionID := cookie.Value
			ok, err := r.ValidateSession(sessionID) //истекла ли сессия
			if !ok {
				response.ErrHandler(w, fmt.Errorf("ValidateSession failed:%s", err.Error()), http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(&cookieResponse, req.WithContext(ctx))
	})
}

// Функция для извлечения имени операции из запроса GraphQL.
func extractOperationName(req *http.Request) (string, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", fmt.Errorf("can't read request body: %s", err)
	}
	defer req.Body.Close()
	req.Body = io.NopCloser(bytes.NewReader(body))

	var query struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(body, &query); err != nil {
		return "", fmt.Errorf("can't decode JSON: %s", err)
	}
	re := regexp.MustCompile(`\b(\w+)\b`)
	matches := re.FindAllString(query.Query, -1)

	for _, match := range matches { //ищем первое слово после mutation
		if strings.EqualFold(match, "mutation") {
			continue
		}
		return match, nil
	}

	return "", err
}

// Функция, определяющая, требует ли данная операция авторизации.
func requiresAuthentication(operationName string) bool {
	switch operationName {
	case "register", "login", "logout":
		return false // Для этих операций авторизация не требуется.
	default:
		return true // Для всех остальных операций авторизация требуется.
	}
}
