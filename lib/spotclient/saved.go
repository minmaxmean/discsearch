package spotclient

import (
	"context"
	"log/slog"

	"github.com/m-nny/discsearch/lib/utils"
	"github.com/zmb3/spotify/v2"
)

func (s *SpotifyClient) SavedTracks(ctx context.Context) ([]spotify.SavedTrack, error) {
	return utils.CachedExec("spotify/saved_tracks/"+s.username, func() ([]spotify.SavedTrack, error) {
		var tracks []spotify.SavedTrack
		resp, err := s.client.CurrentUsersTracks(ctx, spotify.Limit(50))
		for ; err == nil; err = s.client.NextPage(ctx, resp) {
			slog.Debug("spotify.GetUserTracks():", "len(resp.Tracks)=", len(resp.Tracks), "offset", resp.Offset, "total", resp.Total)
			tracks = append(tracks, resp.Tracks...)
		}
		return tracks, nil
	})
}
