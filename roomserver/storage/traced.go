package storage

import (
	"context"

	"github.com/matrix-org/dendrite/roomserver/api"
	"github.com/matrix-org/dendrite/roomserver/storage/shared"
	"github.com/matrix-org/dendrite/roomserver/storage/tables"
	"github.com/matrix-org/dendrite/roomserver/types"
	"github.com/matrix-org/gomatrixserverlib"
	"github.com/opentracing/opentracing-go"
)

type TracedDatabase struct {
	Db Database
}

func (d *TracedDatabase) SupportsConcurrentRoomInputs() bool {
	return d.Db.SupportsConcurrentRoomInputs()
}
func (d *TracedDatabase) RoomInfo(ctx context.Context, roomID string) (*types.RoomInfo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RoomInfo")
	defer span.Finish()
	return d.Db.RoomInfo(ctx, roomID)
}
func (d *TracedDatabase) AddState(ctx context.Context, roomNID types.RoomNID, stateBlockNIDs []types.StateBlockNID, state []types.StateEntry) (types.StateSnapshotNID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AddState")
	defer span.Finish()
	return d.Db.AddState(ctx, roomNID, stateBlockNIDs, state)
}
func (d *TracedDatabase) StateAtEventIDs(ctx context.Context, eventIDs []string) ([]types.StateAtEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StateAtEventIDs")
	defer span.Finish()
	return d.Db.StateAtEventIDs(ctx, eventIDs)
}
func (d *TracedDatabase) EventTypeNIDs(ctx context.Context, eventTypes []string) (map[string]types.EventTypeNID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventTypeNIDs")
	defer span.Finish()
	return d.Db.EventTypeNIDs(ctx, eventTypes)
}
func (d *TracedDatabase) EventStateKeyNIDs(ctx context.Context, eventStateKeys []string) (map[string]types.EventStateKeyNID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventStateKeyNIDs")
	defer span.Finish()
	return d.Db.EventStateKeyNIDs(ctx, eventStateKeys)
}
func (d *TracedDatabase) StateBlockNIDs(ctx context.Context, stateNIDs []types.StateSnapshotNID) ([]types.StateBlockNIDList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StateBlockNIDs")
	defer span.Finish()
	return d.Db.StateBlockNIDs(ctx, stateNIDs)
}
func (d *TracedDatabase) StateEntries(ctx context.Context, stateBlockNIDs []types.StateBlockNID) ([]types.StateEntryList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StateEntries")
	defer span.Finish()
	return d.Db.StateEntries(ctx, stateBlockNIDs)
}
func (d *TracedDatabase) StateEntriesForTuples(ctx context.Context, stateBlockNIDs []types.StateBlockNID, stateKeyTuples []types.StateKeyTuple) ([]types.StateEntryList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StateEntriesForTuples")
	defer span.Finish()
	return d.Db.StateEntriesForTuples(ctx, stateBlockNIDs, stateKeyTuples)
}
func (d *TracedDatabase) Events(ctx context.Context, eventNIDs []types.EventNID) ([]types.Event, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Events")
	defer span.Finish()
	return d.Db.Events(ctx, eventNIDs)
}
func (d *TracedDatabase) SnapshotNIDFromEventID(ctx context.Context, eventID string) (types.StateSnapshotNID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SnapshotNIDFromEventID")
	defer span.Finish()
	return d.Db.SnapshotNIDFromEventID(ctx, eventID)
}
func (d *TracedDatabase) StoreEvent(ctx context.Context, event *gomatrixserverlib.Event, txnAndSessionID *api.TransactionID, authEventNIDs []types.EventNID, isRejected bool) (types.RoomNID, types.StateAtEvent, *gomatrixserverlib.Event, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StoreEvent")
	defer span.Finish()
	return d.Db.StoreEvent(ctx, event, txnAndSessionID, authEventNIDs, isRejected)
}
func (d *TracedDatabase) StateEntriesForEventIDs(ctx context.Context, eventIDs []string) ([]types.StateEntry, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StateEntriesForEventIDs")
	defer span.Finish()
	return d.Db.StateEntriesForEventIDs(ctx, eventIDs)
}
func (d *TracedDatabase) EventStateKeys(ctx context.Context, eventStateKeyNIDs []types.EventStateKeyNID) (map[types.EventStateKeyNID]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventStateKeys")
	defer span.Finish()
	return d.Db.EventStateKeys(ctx, eventStateKeyNIDs)
}
func (d *TracedDatabase) EventNIDs(ctx context.Context, eventIDs []string) (map[string]types.EventNID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventNIDs")
	defer span.Finish()
	return d.Db.EventNIDs(ctx, eventIDs)
}
func (d *TracedDatabase) SetState(ctx context.Context, eventNID types.EventNID, stateNID types.StateSnapshotNID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SetState")
	defer span.Finish()
	return d.Db.SetState(ctx, eventNID, stateNID)
}
func (d *TracedDatabase) EventIDs(ctx context.Context, eventNIDs []types.EventNID) (map[types.EventNID]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventIDs")
	defer span.Finish()
	return d.Db.EventIDs(ctx, eventNIDs)
}
func (d *TracedDatabase) GetLatestEventsForUpdate(ctx context.Context, roomInfo types.RoomInfo) (*shared.LatestEventsUpdater, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetLatestEventsForUpdate")
	defer span.Finish()
	return d.Db.GetLatestEventsForUpdate(ctx, roomInfo)
}
func (d *TracedDatabase) GetTransactionEventID(ctx context.Context, transactionID string, sessionID int64, userID string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetTransactionEventID")
	defer span.Finish()
	return d.Db.GetTransactionEventID(ctx, transactionID, sessionID, userID)
}
func (d *TracedDatabase) LatestEventIDs(ctx context.Context, roomNID types.RoomNID) ([]gomatrixserverlib.EventReference, types.StateSnapshotNID, int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LatestEventIDs")
	defer span.Finish()
	return d.Db.LatestEventIDs(ctx, roomNID)
}
func (d *TracedDatabase) GetInvitesForUser(ctx context.Context, roomNID types.RoomNID, targetUserNID types.EventStateKeyNID) (senderUserIDs []types.EventStateKeyNID, eventIDs []string, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetInvitesForUser")
	defer span.Finish()
	return d.Db.GetInvitesForUser(ctx, roomNID, targetUserNID)
}
func (d *TracedDatabase) SetRoomAlias(ctx context.Context, alias string, roomID string, creatorUserID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SetRoomAlias")
	defer span.Finish()
	return d.Db.SetRoomAlias(ctx, alias, roomID, creatorUserID)
}
func (d *TracedDatabase) GetRoomIDForAlias(ctx context.Context, alias string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetRoomIDForAlias")
	defer span.Finish()
	return d.Db.GetRoomIDForAlias(ctx, alias)
}
func (d *TracedDatabase) GetAliasesForRoomID(ctx context.Context, roomID string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetAliasesForRoomID")
	defer span.Finish()
	return d.Db.GetAliasesForRoomID(ctx, roomID)
}
func (d *TracedDatabase) GetCreatorIDForAlias(ctx context.Context, alias string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetCreatorIDForAlias")
	defer span.Finish()
	return d.Db.GetCreatorIDForAlias(ctx, alias)
}
func (d *TracedDatabase) RemoveRoomAlias(ctx context.Context, alias string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RemoveRoomAlias")
	defer span.Finish()
	return d.Db.RemoveRoomAlias(ctx, alias)
}
func (d *TracedDatabase) MembershipUpdater(ctx context.Context, roomID, targetUserID string, targetLocal bool, roomVersion gomatrixserverlib.RoomVersion) (*shared.MembershipUpdater, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MembershipUpdater")
	defer span.Finish()
	return d.Db.MembershipUpdater(ctx, roomID, targetUserID, targetLocal, roomVersion)
}
func (d *TracedDatabase) GetMembership(ctx context.Context, roomNID types.RoomNID, requestSenderUserID string) (membershipEventNID types.EventNID, stillInRoom, isRoomForgotten bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetMembership")
	defer span.Finish()
	return d.Db.GetMembership(ctx, roomNID, requestSenderUserID)
}
func (d *TracedDatabase) GetMembershipEventNIDsForRoom(ctx context.Context, roomNID types.RoomNID, joinOnly bool, localOnly bool) ([]types.EventNID, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetMembershipEventNIDsForRoom")
	defer span.Finish()
	return d.Db.GetMembershipEventNIDsForRoom(ctx, roomNID, joinOnly, localOnly)
}
func (d *TracedDatabase) EventsFromIDs(ctx context.Context, eventIDs []string) ([]types.Event, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventsFromIDs")
	defer span.Finish()
	return d.Db.EventsFromIDs(ctx, eventIDs)
}
func (d *TracedDatabase) PublishRoom(ctx context.Context, roomID string, publish bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PublishRoom")
	defer span.Finish()
	return d.Db.PublishRoom(ctx, roomID, publish)
}
func (d *TracedDatabase) GetPublishedRooms(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetPublishedRooms")
	defer span.Finish()
	return d.Db.GetPublishedRooms(ctx)
}
func (d *TracedDatabase) GetStateEvent(ctx context.Context, roomID, evType, stateKey string) (*gomatrixserverlib.HeaderedEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetStateEvent")
	defer span.Finish()
	return d.Db.GetStateEvent(ctx, roomID, evType, stateKey)
}
func (d *TracedDatabase) GetRoomsByMembership(ctx context.Context, userID, membership string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetRoomsByMembership")
	defer span.Finish()
	return d.Db.GetRoomsByMembership(ctx, userID, membership)
}
func (d *TracedDatabase) GetBulkStateContent(ctx context.Context, roomIDs []string, tuples []gomatrixserverlib.StateKeyTuple, allowWildcards bool) ([]tables.StrippedEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetBulkStateContent")
	defer span.Finish()
	return d.Db.GetBulkStateContent(ctx, roomIDs, tuples, allowWildcards)
}
func (d *TracedDatabase) JoinedUsersSetInRooms(ctx context.Context, roomIDs []string) (map[string]int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JoinedUsersSetInRooms")
	defer span.Finish()
	return d.Db.JoinedUsersSetInRooms(ctx, roomIDs)
}
func (d *TracedDatabase) GetLocalServerInRoom(ctx context.Context, roomNID types.RoomNID) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetLocalServerInRoom")
	defer span.Finish()
	return d.Db.GetLocalServerInRoom(ctx, roomNID)
}
func (d *TracedDatabase) GetServerInRoom(ctx context.Context, roomNID types.RoomNID, serverName gomatrixserverlib.ServerName) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetServerInRoom")
	defer span.Finish()
	return d.Db.GetServerInRoom(ctx, roomNID, serverName)
}
func (d *TracedDatabase) GetKnownUsers(ctx context.Context, userID, searchString string, limit int) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetKnownUsers")
	defer span.Finish()
	return d.Db.GetKnownUsers(ctx, userID, searchString, limit)
}
func (d *TracedDatabase) GetKnownRooms(ctx context.Context) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetKnownRooms")
	defer span.Finish()
	return d.Db.GetKnownRooms(ctx)
}
func (d *TracedDatabase) ForgetRoom(ctx context.Context, userID, roomID string, forget bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ForgetRoom")
	defer span.Finish()
	return d.Db.ForgetRoom(ctx, userID, roomID, forget)
}
