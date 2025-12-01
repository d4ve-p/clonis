package gdrive

import (
	"fmt"
	"os"

	"github.com/d4ve-p/clonis/internal/database"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

type Service struct {
	store *database.Store
	config *oauth2.Config
}

func New(store *database.Store) *Service {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	appURL := os.Getenv("APP_URL")
	
	if appURL == "" {
		appURL = fmt.Sprintf("http://localhost:%s", os.Getenv("PORT"))
	}
	
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  appURL + "/drive/callback",
		Scopes: []string{
			drive.DriveFileScope,
			drive.DriveMetadataScope,
		},
		Endpoint: google.Endpoint,
	}
	
	return &Service {
		store: store,
		config: config,
	}
}