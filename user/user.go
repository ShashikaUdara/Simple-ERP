package user

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"upper.io/db.v3"
	"upper.io/db.v3/mysql"

	"erp.com/erp/erpdb"

	"golang.org/x/crypto/bcrypt"
)

// user
type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// session
type UserSession struct {
	UserID           string    `json:"user_id"`
	SessionToken     string    `json:"session_token"`
	IPAddress        string    `json:"ip_address"`
	UserAgent        string    `json:"user_agent"`
	CreationTime     time.Time `json:"creation_time"`
	LastActivityTime time.Time `json:"last_activity_time"`
}

const (
	sessionIDLength    = 32
	sessionTokenLength = 250
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

type ClientInfo struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

// GenerateSessionID generates a random session ID of specified length.
func GenerateSessionID() string {
	return generateRandomString(sessionIDLength)
}

// GenerateSessionToken generates a random session token of specified length.
func GenerateSessionToken() string {
	return generateRandomString(sessionTokenLength)
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	// Extract the IP address from the "ip:port" format.
	ip = strings.Split(ip, ":")[0]

	// Parse IP address to check if it's a valid IP.
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "" // Invalid IP address
	}

	return ip
}

func GetClientInfo(r *http.Request) ClientInfo {
	clientIP := GetClientIP(r)
	userAgent := r.Header.Get("User-Agent")

	return ClientInfo{
		IP:        clientIP,
		UserAgent: userAgent,
	}
}

func IsUserExists(email string) (bool, error) {
	session, err := mysql.Open(erpdb.Settings)
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
		return false, err
	}
	defer session.Close()

	table := session.Collection("user")

	cond := db.Cond{"email": email}

	count, err := table.Find(cond).Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func InsertUserData(user User) (int, error) {
	session, err := mysql.Open(erpdb.Settings)
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
		return -1, err
	}
	defer session.Close()

	hashed_password, err := HashPassword(user.Password)
	if err != nil {
		log.Fatalf("Error hasing the password: %v", err)
		return -2, err
	}

	table := session.Collection("user")

	data := erpdb.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashed_password,
	}

	fmt.Println(data)

	_, err = table.Insert(data)
	if err != nil {
		log.Fatalf("Error inserting data into the table: %v", err)
		return -3, err
	}

	count, err := table.Find().Count()
	if err != nil {
		return -4, err
	}

	return int(count), nil
}

func UpdateUserSession(token string) (bool, error) {
	session, err := mysql.Open(erpdb.Settings)
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
		return false, err
	}
	defer session.Close()

	table := session.Collection("user_session")

	updateData := db.Cond{
		"is_active":          false,
		"last_activity_time": time.Now(),
	}

	condition := db.Cond{"session_token": token}

	err = table.Find(condition).Update(updateData)
	if err != nil {
		log.Fatalf("Error updating data in the table: %v", err)
		return false, err
	}

	return true, nil
}

func HashPassword(password string) (string, error) {
	// Generate a salt with a cost factor of 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyUser(user User) (bool, error) {
	fmt.Println("VerifyUser")

	// retrive user data (verify user id (email))
	userData, err := GetUserData(user.Email)
	if err == db.ErrNoMoreRows {
		return false, err
	}

	if err != nil {
		log.Fatalf("Error getting user data: %v", err)
		return false, err
	}

	// verify user password
	isValid := VerifyPassword(userData.Password, user.Password)

	return isValid, nil
}

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GetUserData(user_id string) (*erpdb.User, error) {
	var user erpdb.User

	session, err := mysql.Open(erpdb.Settings)
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
		return nil, err
	}
	defer session.Close()

	err = session.Collection("user").Find(db.Cond{"email": user_id}).One(&user)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, err
		}
		log.Fatalf("Error retrieving user data: %v", err)
		return nil, err
	}

	return &user, nil
}

func CreateUserSession(user User, r *http.Request) (string, error) {
	fmt.Println("CreateUserSession")

	session, err := mysql.Open(erpdb.Settings)
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
		return "-1", err
	}
	defer session.Close()

	table := session.Collection("user_session")

	creation_time := time.Now()
	modified_time := time.Now()

	clientInfo := GetClientInfo(r)
	session_token := GenerateSessionToken()

	data := erpdb.UserSession{
		UserID:           user.Email,
		SessionToken:     session_token,
		IPAddress:        clientInfo.IP,
		UserAgent:        clientInfo.UserAgent,
		CreationTime:     creation_time,
		LastActivityTime: modified_time,
	}

	fmt.Println(data)

	_, err = table.Insert(data)
	if err != nil {
		log.Fatalf("Error inserting data into the table: %v", err)
		return "-2", err
	}

	count, err := table.Find().Count()
	if err != nil {
		return "-3", err
	}

	if int(count) > 0 {
		return session_token, nil
	} else {
		return "-3", err
	}
}
