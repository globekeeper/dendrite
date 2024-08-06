// Copyright 2017-2018 New Vector Ltd
// Copyright 2019-2020 The Matrix.org Foundation C.I.C.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgres

import (
	"context"
	"database/sql"

	"github.com/matrix-org/dendrite/internal"
	"github.com/matrix-org/dendrite/internal/sqlutil"
	"github.com/matrix-org/dendrite/roomserver/storage/tables"
	"github.com/matrix-org/dendrite/roomserver/types"
)

const eventJSONSchema = `
-- Stores the JSON for each event. This kept separate from the main events
-- table to keep the rows in the main events table small.
CREATE TABLE IF NOT EXISTS roomserver_event_json (
    -- Local numeric ID for the event.
    event_nid BIGINT NOT NULL PRIMARY KEY,
    -- The JSON for the event.
    -- Stored as TEXT because this should be valid UTF-8.
    -- Not stored as a JSONB because we always just pull the entire event
    -- so there is no point in postgres parsing it.
    -- Not stored as JSON because we already validate the JSON in the server
    -- so there is no point in postgres validating it.
    -- TODO: Should we be compressing the events with Snappy or DEFLATE?
    event_json TEXT NOT NULL
);
`

const insertEventJSONSQL = "" +
	"INSERT INTO roomserver_event_json (event_nid, event_json) VALUES ($1, $2)" +
	" ON CONFLICT (event_nid) DO UPDATE SET event_json=$2"

// Bulk event JSON lookup by numeric event ID.
// Sort by the numeric event ID.
// This means that we can use binary search to lookup by numeric event ID.
const bulkSelectEventJSONSQL = "" +
	"SELECT event_nid, event_json FROM roomserver_event_json" +
	" WHERE event_nid = ANY($1)" +
	" ORDER BY event_nid ASC"

const selectRoomsUnderSpaceSQL = "" +
	"WITH room_ids AS (" +
	"	SELECT DISTINCT" +
	"		(REGEXP_MATCHES(event_json, '\"room_id\":\"([^\"]+)\"'))[1]::text AS room_id" +
	"	FROM roomserver_event_json" +
	"	WHERE event_json LIKE '%\"state_key\":\"$1\"%'" +
	"	AND event_json LIKE '%\"type\":\"m.space.parent\"%'" +
	")," +
	"dm_rooms AS (" +
	"	SELECT" +
	"		ARRAY_AGG(DISTINCT r.room_id) AS dm_array" +
	"	FROM roomserver_event_json e" +
	"	CROSS JOIN LATERAL (" +
	"		SELECT (REGEXP_MATCHES(e.event_json, '\"room_id\":\"([^\"]+)\"'))[1]::text AS room_id" +
	"	) AS r" +
	"	WHERE e.event_json LIKE '%\"is_direct\":true%'" +
	"	AND r.room_id = ANY (" +
	"		SELECT room_id FROM room_ids" +
	"	)" +
	")," +
	"operation_rooms AS (" +
	"	SELECT" +
	"		ARRAY_AGG(DISTINCT r.room_id) AS operation_array" +
	"	FROM roomserver_event_json e" +
	"	CROSS JOIN LATERAL (" +
	"		SELECT (REGEXP_MATCHES(e.event_json, '\"room_id\":\"([^\"]+)\"'))[1]::text AS room_id" +
	"	) AS r" +
	"	WHERE e.event_json LIKE '%\"type\":\"connect.operation\"%'" +
	"	AND r.room_id = ANY (" +
	"		SELECT room_id FROM room_ids" +
	"	)" +
	")," +
	"team_rooms AS (" +
	"	SELECT" +
	"		ARRAY_AGG(DISTINCT r.room_id) AS team_array" +
	"	FROM roomserver_event_json e" +
	"	CROSS JOIN LATERAL (" +
	"		SELECT (REGEXP_MATCHES(e.event_json, '\"room_id\":\"([^\"]+)\"'))[1]::text AS room_id" +
	"	) AS r" +
	"	WHERE r.room_id = ANY (" +
	"		SELECT room_id FROM room_ids" +
	"	)" +
	"	AND r.room_id NOT IN (" +
	"		SELECT UNNEST(operation_rooms.operation_array) FROM operation_rooms" +
	"	)" +
	"	AND r.room_id NOT IN (" +
	"		SELECT UNNEST(dm_rooms.dm_array) FROM dm_rooms" +
	"	)" +
	")" +
	"SELECT" +
	"	dm_rooms.dm_array," +
	"	operation_rooms.operation_array," +
	"	team_rooms.team_array" +
	"FROM" +
	"	dm_rooms," +
	"	operation_rooms," +
	"	team_rooms;"

type eventJSONStatements struct {
	insertEventJSONStmt          *sql.Stmt
	bulkSelectEventJSONStmt      *sql.Stmt
	selectRoomsUnderSpaceSQLStmt *sql.Stmt
}

func CreateEventJSONTable(db *sql.DB) error {
	_, err := db.Exec(eventJSONSchema)
	return err
}

func PrepareEventJSONTable(db *sql.DB) (tables.EventJSON, error) {
	s := &eventJSONStatements{}

	return s, sqlutil.StatementList{
		{&s.insertEventJSONStmt, insertEventJSONSQL},
		{&s.bulkSelectEventJSONStmt, bulkSelectEventJSONSQL},
		{&s.selectRoomsUnderSpaceSQLStmt, selectRoomsUnderSpaceSQL},
	}.Prepare(db)
}

func (s *eventJSONStatements) InsertEventJSON(
	ctx context.Context, txn *sql.Tx, eventNID types.EventNID, eventJSON []byte,
) error {
	stmt := sqlutil.TxStmt(txn, s.insertEventJSONStmt)
	_, err := stmt.ExecContext(ctx, int64(eventNID), eventJSON)
	return err
}

func (s *eventJSONStatements) BulkSelectEventJSON(
	ctx context.Context, txn *sql.Tx, eventNIDs []types.EventNID,
) ([]tables.EventJSONPair, error) {
	stmt := sqlutil.TxStmt(txn, s.bulkSelectEventJSONStmt)
	rows, err := stmt.QueryContext(ctx, eventNIDsAsArray(eventNIDs))
	if err != nil {
		return nil, err
	}
	defer internal.CloseAndLogIfError(ctx, rows, "bulkSelectEventJSON: rows.close() failed")

	// We know that we will only get as many results as event NIDs
	// because of the unique constraint on event NIDs.
	// So we can allocate an array of the correct size now.
	// We might get fewer results than NIDs so we adjust the length of the slice before returning it.
	results := make([]tables.EventJSONPair, len(eventNIDs))
	i := 0
	var eventNID int64
	for ; rows.Next(); i++ {
		result := &results[i]
		if err := rows.Scan(&eventNID, &result.EventJSON); err != nil {
			return nil, err
		}
		result.EventNID = types.EventNID(eventNID)
	}
	return results[:i], rows.Err()
}

func (s *eventJSONStatements) SelectRoomsUnderSpace(
	ctx context.Context, txn *sql.Tx, spaceID string,
) ([]string, []string, []string, error) {
	stmt := sqlutil.TxStmt(txn, s.selectRoomsUnderSpaceSQLStmt)
	rows, err := stmt.QueryContext(ctx, spaceID)
	if err != nil {
		return nil, nil, nil, err
	}
	defer internal.CloseAndLogIfError(ctx, rows, "SelectRoomsUnderSpaceSQL: rows.close() failed")

	var (
		Dms        []string
		Operations []string
		Teams      []string
	)

	if err := rows.Scan(&Dms, &Operations, &Teams); err != nil {
		return nil, nil, nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, nil, nil, err
	}

	return Dms, Operations, Teams, nil
}
