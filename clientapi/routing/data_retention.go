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
	"github.com/sirupsen/logrus"
)

type DataRetentionRequest struct {
	DataRetention api.PerformDataRetentionRequest `json:"data_retention"`
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
	dr := body.DataRetention

	if dr.MaxAge <= 0 || dr.SpaceID == "" {
		return util.JSONResponse{
			Code: http.StatusBadRequest,
			JSON: spec.BadJSON("missing max_age or space_id"),
		}
	}

	// Validate the roomID
	validRoomID, err := spec.NewRoomID(dr.SpaceID)
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

	if dr.Teams {
		logrus.Infof("Performing data retention on teams in space %s", dr.SpaceID)
		for _, roomId := range queryRes.Teams {
			if err = rsAPI.PerformDataRetention(context.Background(), &dr, roomId); err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"space_id": dr.SpaceID,
				}).Errorf("Failed to perform data retention on team with id %s", roomId)
				return util.ErrorResponse(err)
			}
		}
	}

	if dr.Operations {
		logrus.Infof("Performing data retention on operations in space %s", dr.SpaceID)
		for _, roomId := range queryRes.Operations {
			if err = rsAPI.PerformDataRetention(context.Background(), &dr, roomId); err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"space_id": dr.SpaceID,
				}).Errorf("Failed to perform data retention on operation with id %s", roomId)
				return util.ErrorResponse(err)
			}
		}
	}

	if dr.Dms {
		logrus.Infof("Performing data retention on dms in space %s", dr.SpaceID)
		for _, roomId := range queryRes.DMs {
			if err = rsAPI.PerformDataRetention(context.Background(), &dr, roomId); err != nil {
				logrus.WithError(err).WithFields(logrus.Fields{
					"space_id": dr.SpaceID,
				}).Errorf("Failed to perform data retention on dm with id %s", roomId)
				return util.ErrorResponse(err)
			}
		}
	}

	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: struct{}{},
	}
}
