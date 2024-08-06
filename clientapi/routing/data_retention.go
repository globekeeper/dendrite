package routing

import (
	"context"
	"net/http"

	"github.com/matrix-org/dendrite/clientapi/httputil"
	"github.com/matrix-org/dendrite/roomserver/api"
	roomserverAPI "github.com/matrix-org/dendrite/roomserver/api"
	"github.com/matrix-org/dendrite/setup/config"
	"github.com/matrix-org/gomatrixserverlib/spec"
	"github.com/matrix-org/util"
)

type DataRetentionRequest struct {
	DataRetentions DataRetention `json:"data_retentions"`
}

type DataRetention struct {
	SpaceID    string `json:"space_id"`
	Enabled    bool   `json:"enabled"`
	MaxAge     int32  `json:"max_age,required"`
	Teams      bool   `json:"teams"`
	Operations bool   `json:"operations"`
	Dms        bool   `json:"dms"`
}

// Triggred by an application service job.
// Purges stale data according to data retention policy provided in the request body.
// For large spaces with many rooms this operation may take a considerable amount of time.
func PostDataRetention(
	req *http.Request,
	cfg *config.ClientAPI,
	rsAPI roomserverAPI.ClientRoomserverAPI,
) util.JSONResponse {
	var body DataRetentionRequest
	if reqErr := httputil.UnmarshalJSONRequest(req, &body); reqErr != nil {
		return *reqErr
	}

	if body.DataRetentions.MaxAge <= 0 || body.DataRetentions.SpaceID == "" {
		return util.JSONResponse{
			Code: http.StatusBadRequest,
			JSON: spec.BadJSON("missing max_age or space_id"),
		}
	}

	// Validate the roomID
	validRoomID, err := spec.NewRoomID(body.DataRetentions.SpaceID)
	if err != nil {
		return util.JSONResponse{
			Code: http.StatusBadRequest,
			JSON: spec.InvalidParam("space_id is invalid"),
		}
	}

	queryReq := api.QueryRoomsUnderSpaceRequest{
		SpaceID: validRoomID.String(),
	}

	var queryRes api.QueryRoomsUnderSpaceResponse
	if queryErr := rsAPI.QueryRoomsUnderSpace(req.Context(), &queryReq, &queryRes); queryErr != nil {
		util.GetLogger(req.Context()).WithError(queryErr).Error("rsAPI.QueryRoomsUnderSpace failed")
		return util.JSONResponse{
			Code: http.StatusInternalServerError,
			JSON: spec.InternalServerError{},
		}
	}

	if body.DataRetentions.Teams {
		// TODO: Replace with PerformDataRetention once it's implemented
		for _, roomId := range queryRes.Teams {
			if err = rsAPI.PerformAdminPurgeRoom(context.Background(), roomId); err != nil {
				return util.ErrorResponse(err)
			}
		}
	}

	if body.DataRetentions.Operations {
		for _, roomId := range queryRes.Operations {
			// TODO: Replace with PerformDataRetention once it's implemented
			if err = rsAPI.PerformAdminPurgeRoom(context.Background(), roomId); err != nil {
				return util.ErrorResponse(err)
			}
		}
	}

	if body.DataRetentions.Dms {
		for _, roomId := range queryRes.DMs {
			// TODO: Replace with PerformDataRetention once it's implemented
			if err = rsAPI.PerformAdminPurgeRoom(context.Background(), roomId); err != nil {
				return util.ErrorResponse(err)
			}
		}
	}

	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: struct{}{},
	}
}
