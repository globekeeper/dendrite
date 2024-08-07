// Copyright 2022 The Matrix.org Foundation C.I.C.
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

	"github.com/matrix-org/dendrite/internal/sqlutil"
	"github.com/matrix-org/dendrite/roomserver/api"
	"github.com/matrix-org/dendrite/roomserver/types"
)

const dataRetentionEventJSONSQL = "" +
	"DELETE FROM roomserver_event_json WHERE event_nid = ANY(" +
	"	SELECT re.event_nid" +
	"	FROM roomserver_events AS re" +
	"	JOIN roomserver_event_json AS rej ON re.event_nid = rej.event_nid" +
	"	WHERE rej.event_json LIKE '%\"type\":\"m.room.encrypted\"%'" +
	"	AND CAST(rej.event_json::jsonb->>'origin_server_ts' AS BIGINT) >= EXTRACT(EPOCH FROM NOW()) * 1000 - $2" +
	"	AND re.room_nid = $1" +
	")"

const dataRetentionEventsSQL = "" +
	"DELETE FROM roomserver_events WHERE event_nid = ANY(" +
	"	SELECT re.event_nid" +
	"	FROM roomserver_events AS re" +
	"	JOIN roomserver_event_json AS rej ON re.event_nid = rej.event_nid" +
	"	WHERE rej.event_json LIKE '%\"type\":\"m.room.encrypted\"%'" +
	"	AND CAST(rej.event_json::jsonb->>'origin_server_ts' AS BIGINT) >= EXTRACT(EPOCH FROM NOW()) * 1000 - $2" +
	"	AND re.room_nid = $1" +
	")"

const dataRetentionPreviousEventsSQL = "" +
	"UPDATE roomserver_previous_events SET event_nids = array_remove(event_nids, subquery.event_nid) FROM (" +
	"	SELECT re.event_nid" +
	"	FROM roomserver_events AS re" +
	"	JOIN roomserver_event_json AS rej ON re.event_nid = rej.event_nid" +
	"	WHERE rej.event_json LIKE '%\"type\":\"m.room.encrypted\"%'" +
	"	AND CAST(rej.event_json::jsonb->>'origin_server_ts' AS BIGINT) >= EXTRACT(EPOCH FROM NOW()) * 1000 - $2" +
	"	AND re.room_nid = $1" +
	") AS subquery" +
	"WHERE event_nids @> ARRAY[subquery.event_nid]"

const dataRetentionStateBlockEntriesSQL = "" +
	"UPDATE roomserver_state_block SET event_nids = array_remove(event_nids, subquery.event_nid) FROM (" +
	"	SELECT re.event_nid" +
	"	FROM roomserver_events AS re" +
	"	JOIN roomserver_event_json AS rej ON re.event_nid = rej.event_nid" +
	"	WHERE rej.event_json LIKE '%\"type\":\"m.room.encrypted\"%'" +
	"	AND CAST(rej.event_json::jsonb->>'origin_server_ts' AS BIGINT) >= EXTRACT(EPOCH FROM NOW()) * 1000 - $2" +
	"	AND re.room_nid = $1" +
	") AS subquery" +
	"WHERE event_nids @> ARRAY[subquery.event_nid]"

type dataRetentionStatements struct {
	dataRetentionEventJSONStmt         *sql.Stmt
	dataRetentionEventsStmt            *sql.Stmt
	dataRetentionPreviousEventsStmt    *sql.Stmt
	dataRetentionStateBlockEntriesStmt *sql.Stmt
}

func PrepareDataRetentionStatements(db *sql.DB) (*dataRetentionStatements, error) {
	s := &dataRetentionStatements{}

	return s, sqlutil.StatementList{
		{&s.dataRetentionEventJSONStmt, dataRetentionEventJSONSQL},
		{&s.dataRetentionEventsStmt, dataRetentionEventsSQL},
		{&s.dataRetentionPreviousEventsStmt, dataRetentionPreviousEventsSQL},
		{&s.dataRetentionStateBlockEntriesStmt, dataRetentionStateBlockEntriesSQL},
	}.Prepare(db)
}

func (s *dataRetentionStatements) DataRetentionInRoom(
	ctx context.Context, txn *sql.Tx, dr *api.PerformDataRetentionRequest, roomNID types.RoomNID, roomID string,
) error {
	dataRetentionByRoomNID := []*sql.Stmt{
		s.dataRetentionEventJSONStmt,
		s.dataRetentionEventsStmt,
		s.dataRetentionPreviousEventsStmt,
		s.dataRetentionStateBlockEntriesStmt,
	}
	for _, stmt := range dataRetentionByRoomNID {
		_, err := sqlutil.TxStmt(txn, stmt).ExecContext(ctx, roomNID)
		if err != nil {
			return err
		}
	}
	return nil
}
