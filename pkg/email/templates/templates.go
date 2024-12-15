package templates

import "fmt"

func StartServiceMsg(email string) string {
	return fmt.Sprintf("To: %s\r\n"+
		"Subject: service started\r\n"+
		"\r\n"+
		"\r\n", email)
}

func VerificationCodeMsg(email string, code int) string {
	return fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: verification code\r\n"+
			"\r\n"+
			"%d\r\n", email, code)
}

func ChangeIpAllertMsg(email string) string {
	return fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: allert!\r\n"+
			"\r\n"+
			"Log in from a new device.\r\n", email)
}
