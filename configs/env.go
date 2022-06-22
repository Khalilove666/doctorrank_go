package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var SECRET_KEY = Env("SECRET_KEY")
var PORT = Env("PORT")
var MONGOURI = Env("MONGOURI")
var CLIENT = Env("CLIENT")
var DOMAIN = Env("DOMAIN")
var TOKEN_MINUTES = Env("TOKEN_MINUTES")
var REFRESH_TOKEN_MINUTES = Env("REFRESH_TOKEN_MINUTES")
var FILESYSTEM_PATH = Env("FILESYSTEM_PATH")
var MAIL_SERVER_DOMAIN = Env("MAIL_SERVER_DOMAIN")
var MAIL_SERVER_PORT = Env("MAIL_SERVER_PORT")
var MAIL_SERVER_USERNAME = Env("MAIL_SERVER_USERNAME")
var MAIL_SERVER_PASSWORD = Env("MAIL_SERVER_PASSWORD")
var MAIL_SERVER_EMAIL_FROM = Env("MAIL_SERVER_EMAIL_FROM")

func Env(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(key)
}
