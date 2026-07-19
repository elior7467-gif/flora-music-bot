package platforms

import (
	"context"
	"errors"
	"fmt"

	"github.com/Laky-64/gologging"
	"github.com/amarnathcjd/gogram/telegram"

	"main/internal/config"
	state "main/internal/core/models"
	"main/internal/utils"
)

const PlatformFallenApi state.PlatformName = "FallenApi"

type FallenApiPlatform struct {
	name state.PlatformName
}

func init() {
	Register(80, &FallenApiPlatform{
		name: PlatformFallenApi,
	})
}

func (f *FallenApiPlatform) Name() state.PlatformName {
	return f.name
}

func (f *FallenApiPlatform) CanGetTracks(query string) bool {
	return false
}

func (f *FallenApiPlatform) GetTracks(
	_ context.Context,
	_ string,
	_ bool,
) ([]*state.Track, error) {
	return nil, errors.New("fallenapi is a download-only platform")
}

func (f *FallenApiPlatform) CanDownload(
	source state.PlatformName,
) bool {
	if config.CustomAPIURL == "" {
		return false
	}
	return source == PlatformYouTube
}

func (f *FallenApiPlatform) Download(
	ctx context.Context,
	track *state.Track,
	statusMsg *telegram.NewMessage,
) (string, error) {
	if f := findFile(track); f != "" {
		gologging.Debug("FallenApi: Download -> Cached File -> " + f)
		return f, nil
	}

	var pm *telegram.ProgressManager
	if statusMsg != nil {
		pm = utils.GetProgress(statusMsg)
	}
	_ = pm // resty save-to-file doesn't report progress; kept for future use

	ext := ".mp3"
	if track.Video {
		ext = ".mp4"
	}
	path := getPath(track, ext)

	// Modified to use the requested workers API domain directly
	apiReqURL := fmt.Sprintf("https://testshit-yt.kustbotsweb.workers.dev/down?url=%s", track.URL)

	resp, err := rc.R().
		SetContext(ctx).
		SetResponseSaveFileName(path).
		Get(apiReqURL)
	if err != nil {
		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) {
			return "", err
		}
		return "", fmt.Errorf("failed to download %s: %w", track.URL, err)
	}

	if resp.IsStatusFailure() {
		return "", fmt.Errorf(
			"api request failed for %s with status: %d",
			track.URL,
			resp.StatusCode(),
		)
	}

	if !fileExists(path) {
		return "", errors.New("empty file returned by API")
	}

	return path, nil
}
