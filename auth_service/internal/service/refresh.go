package service

import (
	"auth_service/internal/repository/psql"
	"auth_service/internal/tokens"
	"errors"
	"net/http"
	"time"
)

var ErrInvalidRefresh = errors.New("invalid refresh")

// Принимает строку refresh, возвращает (access, newRefresh)
func RotateRefresh(refresh string, userAgent, ip string) (string, string, error) {
	hash := tokens.HashRefreshToken(refresh)
	row, err := psql.GetRefreshToken(hash)
	if err != nil {
		return "", "", ErrInvalidRefresh
	}
	if row.RevokedAt.Valid || time.Now().After(row.ExpiresAt) {
		return "", "", ErrInvalidRefresh
	}

	// Ротируем: помечаем старый, создаём новый
	_ = psql.RevokeRefreshToken(hash)

	newRefresh, err := tokens.GenerateRefreshOpaque(32)
	if err != nil {
		return "", "", err
	}
	newHash := tokens.HashRefreshToken(newRefresh)
	exp := time.Now().Add(tokens.RefreshTTL())
	if err := psql.SaveRefreshToken(row.UserID, newHash, userAgent, ip, exp); err != nil {
		return "", "", err
	}

	// Новый access
	access, err := tokens.GenerateJWT(row.UserID, "")
	if err != nil {
		return "", "", err
	}
	return access, newRefresh, nil
}

// Утилита для установки cookie
func SetRefreshCookie(w http.ResponseWriter, token string, ttl time.Duration) {
	c := &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // в проде true за HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(ttl.Seconds()),
	}
	http.SetCookie(w, c)
}
