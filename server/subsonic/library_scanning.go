package subsonic

import (
	"net/http"
	"time"

	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model/request"
	"github.com/navidrome/navidrome/server/subsonic/responses"
	"github.com/navidrome/navidrome/utils"
)

func (api *Router) GetScanStatus(r *http.Request) (*responses.Subsonic, error) {
	// TODO handle multiple mediafolders
	ctx := r.Context()
	scanStatus := &responses.ScanStatus{
		Scanning:    false,
		Count:       0,
		FolderCount: 0,
		LastScan:    nil,
	}
	for _, mediaFolder := range conf.Server.MusicFolder {
		status, err := api.scanner.Status(mediaFolder)
		if err != nil {
			log.Error(ctx, "Error retrieving Scanner status", err)
			return nil, newError(responses.ErrorGeneric, "Internal Error")
		}
		if status.Scanning {
			scanStatus.Scanning = status.Scanning
		}
		scanStatus.Count += int64(status.Count)
		scanStatus.FolderCount += int64(status.FolderCount)
		if scanStatus.LastScan == nil || scanStatus.LastScan.Before(status.LastScan) {
			scanStatus.LastScan = &status.LastScan
		}
	}
	response := newResponse()
	response.ScanStatus = scanStatus
	return response, nil
}

func (api *Router) StartScan(r *http.Request) (*responses.Subsonic, error) {
	ctx := r.Context()
	loggedUser, ok := request.UserFrom(ctx)
	if !ok {
		return nil, newError(responses.ErrorGeneric, "Internal error")
	}

	if !loggedUser.IsAdmin {
		return nil, newError(responses.ErrorAuthorizationFail)
	}

	fullScan := utils.ParamBool(r, "fullScan", false)

	go func() {
		start := time.Now()
		log.Info(ctx, "Triggering manual scan", "fullScan", fullScan, "user", loggedUser.UserName)
		err := api.scanner.RescanAll(ctx, fullScan)
		if err != nil {
			log.Error(ctx, "Error scanning", err)
			return
		}
		log.Info(ctx, "Manual scan complete", "user", loggedUser.UserName, "elapsed", time.Since(start).Round(100*time.Millisecond))
	}()

	return api.GetScanStatus(r)
}
