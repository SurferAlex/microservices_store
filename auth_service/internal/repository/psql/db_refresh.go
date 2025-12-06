package psql

import (
	"database/sql"
	"time"
)

type RefreshRow struct {
	TokenHash string
	UserID    int
	ExpiresAt time.Time
	RevokedAt sql.NullTime
	UserAgent string
	IP        string
}

func SaveRefreshToken(userID int, tokenHash, userAgent, ip string, expiresAt time.Time) error {
	_, err := db.Exec(`
	    INSERT INTO refresh_tokens (token_hash, user_id, expires_at, user_agent, ip)
		VALUES ($1, $2, $3, $4, $5)
		`, tokenHash, userID, expiresAt, userAgent, ip)
	return err
}

func GetRefreshToken(tokenHash string) (*RefreshRow, error) {
	row := db.QueryRow(`
		SELECT token_hash, user_id, expires_at, revoked_at, COALESCE(user_agent,''), COALESCE(ip,'')
		FROM refresh_tokens
		WHERE token_hash = $1
	`, tokenHash)
	var r RefreshRow
	if err := row.Scan(&r.TokenHash, &r.UserID, &r.ExpiresAt, &r.RevokedAt, &r.UserAgent, &r.IP); err != nil {
		return nil, err
	}
	return &r, nil
}

func RevokeRefreshToken(tokenHash string) error {
	_, err := db.Exec(`UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = $1`, tokenHash)
	return err
}
