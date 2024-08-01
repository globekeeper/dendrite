package routing

import (
	"net/http"

	"github.com/matrix-org/dendrite/clientapi/httputil"
	"github.com/matrix-org/dendrite/setup/config"
	"github.com/matrix-org/dendrite/userapi/api"
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
	deviceAPI *api.Device,
	userAPI api.ClientUserAPI,
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

	// TODO: Fetch dms, operators and teams under the provided space.
	// WITH room_ids AS (
	// 	SELECT DISTINCT
	// 		(REGEXP_MATCHES(event_json, '"room_id":"([^"]+)"'))[1]::text AS room_id
	// 	FROM roomserver_event_json
	// 	WHERE event_json LIKE '%"state_key":"$1"%'
	// 	AND event_json LIKE '%"type":"m.space.parent"%'
	// ),
	// dm_rooms AS (
	// 	SELECT
	// 		ARRAY_AGG(DISTINCT r.room_id) AS dm_array
	// 	FROM roomserver_event_json e
	// 	CROSS JOIN LATERAL (
	// 		SELECT (REGEXP_MATCHES(e.event_json, '"room_id":"([^"]+)"'))[1]::text AS room_id
	// 	) AS r
	// 	WHERE e.event_json LIKE '%"is_direct":true%'
	// 	AND r.room_id = ANY (
	// 		SELECT room_id FROM room_ids
	// 	)
	// ),
	// operation_rooms AS (
	// 	SELECT
	// 		ARRAY_AGG(DISTINCT r.room_id) AS operation_array
	// 	FROM roomserver_event_json e
	// 	CROSS JOIN LATERAL (
	// 		SELECT (REGEXP_MATCHES(e.event_json, '"room_id":"([^"]+)"'))[1]::text AS room_id
	// 	) AS r
	// 	WHERE e.event_json LIKE '%"type":"connect.operation"%'
	// 	AND r.room_id = ANY (
	// 		SELECT room_id FROM room_ids
	// 	)
	// ),
	// team_rooms AS (
	// 	SELECT
	// 		ARRAY_AGG(DISTINCT r.room_id) AS team_array
	// 	FROM roomserver_event_json e
	// 	CROSS JOIN LATERAL (
	// 		SELECT (REGEXP_MATCHES(e.event_json, '"room_id":"([^"]+)"'))[1]::text AS room_id
	// 	) AS r
	// 	WHERE r.room_id = ANY (
	// 		SELECT room_id FROM room_ids
	// 	)
	// 	AND r.room_id NOT IN (
	// 		SELECT UNNEST(operation_rooms.operation_array) FROM operation_rooms
	// 	)
	// 	AND r.room_id NOT IN (
	// 		SELECT UNNEST(dm_rooms.dm_array) FROM dm_rooms
	// 	)
	// )
	// SELECT
	// 	dm_rooms.dm_array,
	// 	operation_rooms.operation_array,
	// 	team_rooms.team_array
	// FROM
	// 	dm_rooms,
	// 	operation_rooms,
	// 	team_rooms;

	if body.DataRetentions.Teams {
		// TODO: Iterate and purge stale data from teams
	}

	if body.DataRetentions.Operations {
		// TODO: Iterate and purge stale data from operations
	}

	if body.DataRetentions.Dms {
		// TODO: Iterate and purge stale data from dms
	}

	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: struct{}{},
	}
}
