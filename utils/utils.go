package utils

import (
	"SmallZomato/models"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/sirupsen/logrus"
	"github.com/teris-io/shortid"
	"golang.org/x/crypto/bcrypt"
)

var generator *shortid.Shortid

const generatorSeed = 1000

type FieldError struct {
	Err validator.ValidationErrors
}

func (q FieldError) GetSingleError() string {
	errorString := ""
	for _, e := range q.Err {
		errorString = "Invalid " + e.Field()
	}
	return errorString
}

type RequestErr struct {
	// ID for the request
	// Example: 8YeCqPXmM
	ID string `json:"id"`

	// MessageToUser will contain error message
	// Example: Invalid Email
	MessageToUser string `json:"messageToUser"`

	// DeveloperInfo will contain additional developer info related with error
	// Example: Invalid email format
	DeveloperInfo string `json:"developerInfo"`

	// Err contains the error or exception message
	// Example: validation on email failed with error invalid email format
	Err string `json:"error"`

	// StatusCode will contain the status code for the error
	// Example: 500
	StatusCode int `json:"statusCode"`

	// IsClientError will be false if some internal server error occurred
	IsClientError bool `json:"isClientError"`
} // @name RequestErr
type clientError struct {
	ID            string `json:"id"`
	MessageToUser string `json:"messageToUser"`
	DeveloperInfo string `json:"developerInfo"`
	Err           string `json:"error"`
	StatusCode    int    `json:"statusCode"`
	IsClientError bool   `json:"isClientError"`
}

func init() {
	n, err := rand.Int(rand.Reader, big.NewInt(generatorSeed))
	if err != nil {
		logrus.Panicf("failed to initialize utilities with random seed, %+v", err)
		return
	}

	g, err := shortid.New(1, shortid.DefaultABC, n.Uint64())

	if err != nil {
		logrus.Panicf("Failed to initialize utils package with error: %+v", err)
	}

	generator = g
}

// ParseBody parses the values from io reader to a given interface
func ParseBody(body io.Reader, out interface{}) error {
	err := json.NewDecoder(body).Decode(out)
	if err != nil {
		return err
	}

	return nil
}

// EncodeJSONBody writes the JSON body to response writer
func EncodeJSONBody(resp http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(resp).Encode(data)
}

// RespondJSON sends the interface as a JSON
func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	if body != nil {
		if err := EncodeJSONBody(w, body); err != nil {
			logrus.Errorf("Failed to respond JSON with error: %+v", err)
		}
	}
}

// newClientError creates structured client error response message
func newClientError(err error, statusCode int, messageToUser string, additionalInfoForDevs ...string) *clientError {
	additionalInfoJoined := strings.Join(additionalInfoForDevs, "\n")
	if additionalInfoJoined == "" {
		additionalInfoJoined = messageToUser
	}

	errorID, _ := generator.Generate()
	var errString string
	if err != nil {
		errString = err.Error()
	}
	return &clientError{
		ID:            errorID,
		MessageToUser: messageToUser,
		DeveloperInfo: additionalInfoJoined,
		Err:           errString,
		StatusCode:    statusCode,
		IsClientError: true,
	}
}

// RespondError sends an error message to the API caller and logs the error
func RespondError(w http.ResponseWriter, statusCode int, err error, messageToUser string, additionalInfoForDevs ...string) {
	logrus.Errorf("status: %d, message: %s, err: %+v ", statusCode, messageToUser, err)
	clientError := newClientError(err, statusCode, messageToUser, additionalInfoForDevs...)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(clientError); err != nil {
		logrus.Errorf("Failed to send error to caller with error: %+v", err)
	}
}

// HashString generates SHA256 for a given string
func HashString(toHash string) string {
	sha := sha512.New()
	sha.Write([]byte(toHash))
	return hex.EncodeToString(sha.Sum(nil))
}

// HashPassword returns the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// CheckValidation returns the current validation status
func CheckValidation(i interface{}) validator.ValidationErrors {
	v := validator.New()
	err := v.Struct(i)
	if err == nil {
		return nil
	}
	return err.(validator.ValidationErrors)
}

// TrimAll removes a given rune form given string
func TrimAll(str string, remove rune) string {
	return strings.Map(func(r rune) rune {
		if r == remove {
			return -1
		}
		return r
	}, str)
}

// TrimStringAfter trims anything after given delimiter
func TrimStringAfter(s, delim string) string {
	if idx := strings.Index(s, delim); idx != -1 {
		return s[:idx]
	}
	return s
}
func GetDistanceFromLat(location models.Location, latitude string, longitude string) float64 {
	var R = float64(6371) // Radius of the earth in km
	l1, _ := strconv.ParseFloat(latitude, 8)
	lo1, _ := strconv.ParseFloat(longitude, 8)
	l2, _ := strconv.ParseFloat(location.Latitude, 8)
	lo2, _ := strconv.ParseFloat(location.Longitude, 8)
	var dLat = deg2rad(l2 - l1) // deg2rad below
	var dLon = deg2rad(lo2 - lo1)
	var a = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(l1))*math.Cos(deg2rad(l2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var d = R * c // Distance in km
	return d
}
func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
func GenerateJWT(email string, role models.Role) (string, error) {

	var mySigningKey = []byte("123")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
