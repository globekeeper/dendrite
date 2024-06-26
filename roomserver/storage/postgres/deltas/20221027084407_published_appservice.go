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

package deltas

import (
	"context"
	"database/sql"
	"fmt"
)

func UpPulishedAppservice(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `ALTER TABLE roomserver_published ADD COLUMN IF NOT EXISTS appservice_id TEXT NOT NULL DEFAULT '';`)
	if err != nil {
		return fmt.Errorf("failed to execute upgrade: %w", err)
	}
	_, err = tx.ExecContext(ctx, `ALTER TABLE roomserver_published ADD COLUMN IF NOT EXISTS network_id TEXT NOT NULL DEFAULT '';`)
	if err != nil {
		return fmt.Errorf("failed to execute upgrade: %w", err)
	}
	_, err = tx.ExecContext(ctx, `
		ALTER TABLE roomserver_published DROP CONSTRAINT IF EXISTS roomserver_published_pkey;
		ALTER TABLE roomserver_published ADD PRIMARY KEY (room_id, appservice_id, network_id);
	`)
	if err != nil {
		return fmt.Errorf("failed to execute upgrade: %w", err)
	}
	return nil
}

func DownPublishedAppservice(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `ALTER TABLE roomserver_published DROP COLUMN IF EXISTS appservice_id;`)
	if err != nil {
		return fmt.Errorf("failed to execute downgrade: %w", err)
	}
	_, err = tx.ExecContext(ctx, `ALTER TABLE roomserver_published DROP COLUMN IF EXISTS network_id;`)
	if err != nil {
		return fmt.Errorf("failed to execute downgrade: %w", err)
	}
	return nil
}
