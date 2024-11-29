package main

import (
	"log/slog"

	"github.com/hackandpray/MediaCurator/utils"
)

func main() {
	mailer, err := utils.NewEmailSender("testing@hackandpray.com")
	if err != nil {
		slog.Error("Error creating email sender", "error", err)
		return
	}

	body := "Hello, world!"
	err = mailer.SendEmail("jwhenry28@gmail.com", "Testing", body)
	if err != nil {
		slog.Error("Error sending email", "error", err)
	}
}
