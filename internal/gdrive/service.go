package gdrive

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/d4ve-p/clonis/internal/database"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Service struct {
	store *database.Store
	config *oauth2.Config
}

var service *Service = nil

func Get(store *database.Store) *Service {
	if service != nil {
		return service
	}
	
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
	
	service = &Service {
		store: store,
		config: config,
	}
	
	return service
}

func (s *Service) GetAuthURL() string {
	return s.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (s *Service) HandleCallback(ctx context.Context, code string) error {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to exchange token: %w", err)
	}
	
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return err
	}
	
	return s.store.UpdateSetting("gdrive_token", string(tokenJSON))
}

func (s* Service) GetClient(ctx context.Context) (*drive.Service, error) {
	settings, err := s.store.GetSettings()
	if err != nil {
		return nil, err
	}
	
	tokenStr, ok := settings["gdrive_token"]
	if !ok || tokenStr == "" {
		return nil, errors.New("No google drive token found")
	}
	
	var token oauth2.Token
	if err := json.Unmarshal([]byte(tokenStr), &token); err != nil {
		return nil, fmt.Errorf("Invalid token format: %w", err)
	}
	
	tokenSource := s.config.TokenSource(ctx, &token)
	
	srv, err := drive.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, err
	}
	
	return srv, nil
}

func (s *Service) IsConnected() bool {
	settings, err := s.store.GetSettings()
	if err != nil {
		return false
	}
	
	_, ok := settings["gdrive_token"]
	return ok
}

