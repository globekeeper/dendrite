// Copyright 2021 Dan Peleg <dan@globekeeper.com>
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
	"github.com/matrix-org/dendrite/setup/config"
	"github.com/matrix-org/dendrite/userapi/api"
	"github.com/matrix-org/gomatrixserverlib"
)

// Database represents a pusher database.
type Database struct {
	db      *sql.DB
	pushers pushersStatements
}

// NewDatabase creates a new puser database
func NewDatabase(dbProperties *config.DatabaseOptions, serverName gomatrixserverlib.ServerName) (*Database, error) {
	db, err := sqlutil.Open(dbProperties)
	if err != nil {
		return nil, err
	}
	d := pushersStatements{}

	// Create tables before executing migrations so we don't fail if the table is missing,
	// and THEN prepare statements so we don't fail due to referencing new columns
	if err = d.execSchema(db); err != nil {
		return nil, err
	}
	m := sqlutil.NewMigrations()
	if err = m.RunDeltas(db, dbProperties); err != nil {
		return nil, err
	}

	if err = d.prepare(db, serverName); err != nil {
		return nil, err
	}

	return &Database{db, d}, nil
}

// GetPushersByLocalpart returns the pusers matching the given localpart.
func (d *Database) GetPushersByLocalpart(
	ctx context.Context, localpart string,
) ([]api.Pusher, error) {
	return d.pushers.selectPushersByLocalpart(ctx, nil, localpart)
}

// GetPushersByPushkey returns the pusers matching the given localpart.
func (d *Database) GetPushersByPushkey(
	ctx context.Context, localpart, pushkey string,
) (*api.Pusher, error) {
	return d.pushers.selectPushersByPushkey(ctx, localpart, pushkey)
}
