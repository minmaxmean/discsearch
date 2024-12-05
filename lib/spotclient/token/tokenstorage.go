package token

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/m-nny/discsearch/lib/utils"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type TokenStorage interface {
	GetSpotifyToken(ctx context.Context, username string) (*oauth2.Token, error)
	StoreSpotifyToken(ctx context.Context, username string, spotifyToken *oauth2.Token) error
}

func GetToken(ctx context.Context, auth *spotifyauth.Authenticator, serverAddress string, tokenStorage TokenStorage, username string) (*oauth2.Token, error) {
	storedToken, err := tokenStorage.GetSpotifyToken(ctx, username)
	{
		var tokenExpiry time.Time
		if storedToken != nil {
			tokenExpiry = storedToken.Expiry
		}
		slog.DebugContext(ctx, "token.GetToken():", "storedToken.Expiriry", tokenExpiry, "storedToken.Valid", storedToken.Valid(), "err", err)
	}
	if err == nil && storedToken.Valid() {
		return storedToken, nil
	}
	freshToken, err := GetFreshToken(ctx, auth, serverAddress)
	if err != nil {
		return nil, err
	}
	if err := tokenStorage.StoreSpotifyToken(ctx, username, freshToken); err != nil {
		return nil, err
	}
	return freshToken, nil
}

func GetFreshToken(ctx context.Context, auth *spotifyauth.Authenticator, serverAddress string) (*oauth2.Token, error) {
	state := getState()
	url := auth.AuthURL(state)
	fmt.Fprintf(os.Stderr, "Login using following url: %s\n", url)
	if err := utils.OpenBrowser(url); err != nil {
		slog.ErrorContext(ctx, "Could not automatically open browser: %v", "err", err)
	}
	tokenCh := make(chan *oauth2.Token)
	errCh := make(chan error)
	callbackHandler := func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusNotFound)
			errCh <- err
			return
		}
		tokenCh <- token
	}
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/callback", callbackHandler)
	server := &http.Server{Addr: serverAddress, Handler: serverMux}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Could not start server", "err", err)
			errCh <- err
		}
	}()
	select {
	case token := <-tokenCh:
		if err := server.Shutdown(ctx); err != nil {
			return nil, fmt.Errorf("could not stop server: %v", err)
		}
		return token, nil
	case err := <-errCh:
		return nil, err
	}
}

func getState() string {
	return "42"
}

var (
	tsType      = flag.String("spotify_tstype", "json", "")
	tsSingleton TokenStorage
)

func GetTokenStorage() (TokenStorage, error) {
	if tsSingleton != nil {
		return tsSingleton, nil
	}
	if *tsType == "inmemory" {
		tsSingleton = NewInMemoryTokenStorage()
	} else if *tsType == "json" {
		if ts, err := NewJsonTokenStorage(); err != nil {
			return nil, err
		} else {
			tsSingleton = ts
		}
	} else {
		return nil, fmt.Errorf("unknown token storage type: %q", *tsType)
	}
	return tsSingleton, nil
}
