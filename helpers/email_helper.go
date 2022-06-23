package helpers

import (
	"doctorrank_go/configs"
	mail "github.com/xhit/go-simple-mail/v2"
	"strconv"
	"sync"
)

var instantiated *mail.SMTPClient
var once sync.Once

func SmtpClient() (*mail.SMTPClient, error) {
	once.Do(func() {
		server := mail.NewSMTPClient()
		server.Host = configs.MAIL_SERVER_DOMAIN
		port, _ := strconv.ParseInt(configs.MAIL_SERVER_PORT, 10, 64)
		server.Port = int(port)
		server.Username = configs.MAIL_SERVER_USERNAME
		server.Password = configs.MAIL_SERVER_PASSWORD
		server.Encryption = mail.EncryptionTLS

		smtpClient, err := server.Connect()
		if err == nil {
			instantiated = smtpClient
		}
	})
	return instantiated, nil
}

func SendConfirmationMail(name, emailAddress, link string) error {
	smtpClient, err := SmtpClient()
	// Create email
	email := mail.NewMSG()
	email.SetFrom("Doctorrank <" + configs.MAIL_SERVER_EMAIL_FROM + ">")
	email.AddTo(emailAddress)
	email.SetSubject("Email Confirmation")
	email.SetBody(mail.TextHTML, getHtml(name, link))

	// Send email
	err = email.Send(smtpClient)
	return err
}

func SendPasswordResetEmail(emailAddress, link string) error {
	smtpClient, err := SmtpClient()
	// Create email
	email := mail.NewMSG()
	email.SetFrom("Doctorrank <" + configs.MAIL_SERVER_EMAIL_FROM + ">")
	email.AddTo(emailAddress)
	email.SetSubject("Reset Password")
	email.SetBody(mail.TextHTML, getPasswordResetHtml(link))

	// Send email
	err = email.Send(smtpClient)
	return err
}

func getHtml(name, link string) string {
	htmlBody := `
		<!DOCTYPE html>
		<html>
		<head>
		    <meta http-equiv="Content-type" content="text/html" charset="UTF-8">
		    <title>Document</title>
			<style>
		        .btn-activate {
		            background-color: red;
		            padding: 10px 20px;
		            color: white;
		            border-radius: 4px;
		            font-weight: 600;
		        }
		    </style>
		</head>
		<body>
			<h1>Welcome ` + name + ` !</h1>
			<p>Please follow the link below to activate your account.</p>
			<a href="` + link + `" class="btn-activate" target="_blank">ACTIVATE</a>
		    <p>The link will expire within 24 hours.</p>
		    <br>
		    <p>Regards,<br>Doctorrank team</p>
		</body>
		</html>	
	`
	return htmlBody
}

func getPasswordResetHtml(link string) string {
	htmlBody := `
		<!DOCTYPE html>
		<html>
		<head>
		    <meta http-equiv="Content-type" content="text/html" charset="UTF-8">
		    <title>Document</title>
			<style>
		        .btn-activate {
		            background-color: red;
		            padding: 10px 20px;
		            color: white;
		            border-radius: 4px;
		            font-weight: 600;
		        }
		    </style>
		</head>
		<body>
			<h1>Reset password!</h1>
			<p>Please follow the link below to reset your password.</p>
			<a href="` + link + `" class="btn-activate" target="_blank">RESET PASSWORD</a>
		    <p>The link will expire within 24 hours.</p>
		    <br>
		    <p>Regards,<br>Doctorrank team</p>
		</body>
		</html>	
	`
	return htmlBody
}
