package routing

import (
	"net/http"

	"github.com/matrix-org/dendrite/clientapi/auth"
	"github.com/matrix-org/dendrite/clientapi/httputil"
	"github.com/matrix-org/dendrite/userapi/api"
	"github.com/matrix-org/util"
)

type DataRetentionRequest struct {
	DataRetentions DataRetention `json:"data_retentions"`
}

type DataRetention struct {
	SpaceID   string `json:"space_id"`
	Timeframe int64  `json:"timeframe"`
}

// Triggred by an application service job.
// Purges stale data according to data retention policy.
// For large spaces with many rooms this operation may take a considerable amount of time.
func AdminDataRetention(
	req *http.Request,
	userInteractiveAuth *auth.UserInteractive,
	accountAPI api.ClientUserAPI,
	deviceAPI *api.Device,
) util.JSONResponse {
	var body DataRetentionRequest
	if reqErr := httputil.UnmarshalJSONRequest(req, &body); reqErr != nil {
		return *reqErr
	}

	// TODO: Iterate over spaces, fetch associated rooms with SQL:
	// SELECT DISTINCT
	// REGEXP_MATCHES(event_json, '"room_id":"([^"]+)"') AS room_id
	// FROM roomserver_event_json
	// WHERE event_json LIKE '%"state_key":"<space_id>"%'
	// AND event_json LIKE '%"type":"m.space.parent"%';

	// TODO: Purge stale data from rooms
	// TODO: Go over /purge PR to check which other components of Dendrite need to be purged.

	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: struct{}{},
	}
}
