package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Kzynyi/ass/internal/err"
)

type User struct {
	Id       uint32
	Name     string
	Password string
}

var Username string
var Password string

func getConn(ctx context.Context) *pgx.Conn {
	conn, err := pgx.Connect(ctx, "postgres://"+Username+":"+Password+"@localhost:5432/go_ang")
	if err != nil {
		panic(err)
	}
	return conn
}

func CreateUser(ctx context.Context, user *User) (int64, error) {

	conn := getConn(ctx)
	defer conn.Close(ctx)
	exists := checkExistingInfo(conn, ctx, user.Name)
	if exists {
		e := new(err.UserExists)
		e.Name = user.Name
		return 0, e
	}
	tag, err := conn.Exec(ctx, "insert into users values ($1, $2, $3)", uuid.New().ID(), encodePassword(user.Password), user.Name)
	if err != nil {
		panic(err)
	}
	if r := tag.RowsAffected(); r > 0 {
		return r, nil
	}
	return 0, nil
}

func SignInUser(ctx context.Context, authInfo map[string]string) (string, error) {
	conn := getConn(ctx)
	defer conn.Close(ctx)
	count, password := checkUser(conn, ctx, authInfo["Name"])
	fmt.Println(count, password)
	if count < 1 {
		e := new(err.UserNotFound)
		e.Name = authInfo["Name"]
		return "", e
	}
	status := checkPassword(authInfo["Password"], password)
	if !status {
		e := new(err.PasswordMismatch)
		return "", e
	}
	token := generateToken()
	return token, nil
}

func generateToken() string {
	key := []byte("secretkey")
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	token, err := t.SignedString(key)
	if err != nil {
		panic(err)
	}
	return token
}

func checkUser(conn *pgx.Conn, ctx context.Context, name string) (count int, password string) {
	row := conn.QueryRow(ctx, "select count(id), users.password from users where uname = $1 group by users.password", name)
	row.Scan(&count, &password)
	return
}

func checkPassword(normalPw, encodedPw string) bool {
	if res := encodePassword(normalPw); res == encodedPw {
		return true
	} else {
		return false
	}
}

func checkExistingInfo(conn *pgx.Conn, ctx context.Context, name string) bool {
	var count int = 0
	row := conn.QueryRow(ctx, "select count(id) from users where uname like '%' || $1 || '%'", name)
	row.Scan(&count)
	fmt.Println(count)
	return count > 0
}

func encodePassword(password string) string {

	text := []byte(password)
	hash := sha256.Sum256(text)
	return hex.EncodeToString(hash[:])

}
