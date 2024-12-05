package spotclient

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/m-nny/discsearch/lib/spotclient/token"
)

var (
	clientID     = flag.String("spotify_clientid", "", "")
	clientSecret = flag.String("spotify_clientsecret", "", "")
	authAdress   = flag.String("spotify_authaddress", "localhost:3000", "")
)

type SpotifyClient struct {
	client   *spotify.Client
	username string
}

func validateFlags() error {
	if *clientID == "" {
		return fmt.Errorf("spotify_clientid is not set")
	}
	if *clientSecret == "" {
		return fmt.Errorf("spotify_clientsecret is not set")
	}
	if *authAdress == "" {
		return fmt.Errorf("spotify_authaddress is not set")
	}
	return nil
}

func New(ctx context.Context, username string, tokenStorage token.TokenStorage) (*SpotifyClient, error) {
	if err := validateFlags(); err != nil {
		return nil, err
	}
	auth := spotifyauth.New(
		spotifyauth.WithClientID(*clientID),
		spotifyauth.WithClientSecret(*clientSecret),
		spotifyauth.WithRedirectURL(fmt.Sprintf("http://%s/callback", *authAdress)),
		spotifyauth.WithScopes(spotifyauth.ScopeUserLibraryRead),
	)
	token, err := token.GetToken(ctx, auth, *authAdress, tokenStorage, username)
	if err != nil {
		return nil, err
	}
	client := spotify.New(auth.Client(ctx, token))
	return &SpotifyClient{client: client, username: username}, nil
}

func (s *SpotifyClient) Ping(ctx context.Context) error {
	me, err := s.client.CurrentUser(ctx)
	if err != nil {
		return err
	}
	slog.Debug("[spotifyClient.Ping]\n", "me", me)
	return nil
}
