package internal

import (
	"context"
	"crypto/ed25519"

	"github.com/getsentry/sentry-go"
	"github.com/matrix-org/gomatrixserverlib"
	"github.com/matrix-org/gomatrixserverlib/fclient"
	"github.com/matrix-org/gomatrixserverlib/spec"
	"github.com/matrix-org/util"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"

	asAPI "github.com/matrix-org/dendrite/appservice/api"
	fsAPI "github.com/matrix-org/dendrite/federationapi/api"
	"github.com/matrix-org/dendrite/internal/caching"
	"github.com/matrix-org/dendrite/roomserver/acls"
	"github.com/matrix-org/dendrite/roomserver/api"
	"github.com/matrix-org/dendrite/roomserver/internal/input"
	"github.com/matrix-org/dendrite/roomserver/internal/perform"
	"github.com/matrix-org/dendrite/roomserver/internal/query"
	"github.com/matrix-org/dendrite/roomserver/producers"
	"github.com/matrix-org/dendrite/roomserver/storage"
	"github.com/matrix-org/dendrite/roomserver/types"
	"github.com/matrix-org/dendrite/setup/config"
	"github.com/matrix-org/dendrite/setup/jetstream"
	"github.com/matrix-org/dendrite/setup/process"
	userapi "github.com/matrix-org/dendrite/userapi/api"
)

// RoomserverInternalAPI is an implementation of api.RoomserverInternalAPI
type RoomserverInternalAPI struct {
	*input.Inputer
	*query.Queryer
	*perform.Inviter
	*perform.Joiner
	*perform.Peeker
	*perform.InboundPeeker
	*perform.Unpeeker
	*perform.Leaver
	*perform.Publisher
	*perform.Backfiller
	*perform.Forgetter
	*perform.Upgrader
	*perform.Admin
	*perform.Creator
	*perform.DataRetention
	ProcessContext         *process.ProcessContext
	DB                     storage.Database
	Cfg                    *config.Dendrite
	Cache                  caching.RoomServerCaches
	ServerName             spec.ServerName
	KeyRing                gomatrixserverlib.JSONVerifier
	ServerACLs             *acls.ServerACLs
	fsAPI                  fsAPI.RoomserverFederationAPI
	asAPI                  asAPI.AppServiceInternalAPI
	NATSClient             *nats.Conn
	JetStream              nats.JetStreamContext
	Durable                string
	InputRoomEventTopic    string // JetStream topic for new input room events
	OutputProducer         *producers.RoomEventProducer
	PerspectiveServerNames []spec.ServerName
	enableMetrics          bool
	defaultRoomVersion     gomatrixserverlib.RoomVersion
}

// CurrentStateEvent implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).CurrentStateEvent of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) CurrentStateEvent(ctx context.Context, roomID spec.RoomID, eventType string, stateKey string) (gomatrixserverlib.PDU, error) {
	panic("unimplemented")
}

// InputRoomEvents implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Inputer).InputRoomEvents of RoomserverInternalAPI.Inputer.
func (r *RoomserverInternalAPI) InputRoomEvents(ctx context.Context, req *api.InputRoomEventsRequest, res *api.InputRoomEventsResponse) {
	panic("unimplemented")
}

// InvitePending implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).InvitePending of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) InvitePending(ctx context.Context, roomID spec.RoomID, senderID spec.SenderID) (bool, error) {
	panic("unimplemented")
}

// JoinedUserCount implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).JoinedUserCount of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) JoinedUserCount(ctx context.Context, roomID string) (int, error) {
	panic("unimplemented")
}

// LocallyJoinedUsers implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).LocallyJoinedUsers of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) LocallyJoinedUsers(ctx context.Context, roomVersion gomatrixserverlib.RoomVersion, roomNID types.RoomNID) ([]gomatrixserverlib.PDU, error) {
	panic("unimplemented")
}

// PerformAdminDownloadState implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Admin).PerformAdminDownloadState of RoomserverInternalAPI.Admin.
func (r *RoomserverInternalAPI) PerformAdminDownloadState(ctx context.Context, roomID string, userID string, serverName spec.ServerName) error {
	panic("unimplemented")
}

// PerformAdminEvacuateRoom implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Admin).PerformAdminEvacuateRoom of RoomserverInternalAPI.Admin.
func (r *RoomserverInternalAPI) PerformAdminEvacuateRoom(ctx context.Context, roomID string) (affected []string, err error) {
	panic("unimplemented")
}

// PerformAdminEvacuateUser implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Admin).PerformAdminEvacuateUser of RoomserverInternalAPI.Admin.
func (r *RoomserverInternalAPI) PerformAdminEvacuateUser(ctx context.Context, userID string) (affected []string, err error) {
	panic("unimplemented")
}

// PerformAdminPurgeRoom implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Admin).PerformAdminPurgeRoom of RoomserverInternalAPI.Admin.
func (r *RoomserverInternalAPI) PerformAdminPurgeRoom(ctx context.Context, roomID string) error {
	panic("unimplemented")
}

// PerformBackfill implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Backfiller).PerformBackfill of RoomserverInternalAPI.Backfiller.
func (r *RoomserverInternalAPI) PerformBackfill(ctx context.Context, req *api.PerformBackfillRequest, res *api.PerformBackfillResponse) error {
	panic("unimplemented")
}

// PerformDataRetention implements api.RoomserverInternalAPI.
func (r *RoomserverInternalAPI) PerformDataRetention(ctx context.Context, req *api.PerformDataRetentionRequest, res *api.PerformDataRetentionResponse) error {
	panic("unimplemented")
}

// PerformInboundPeek implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*InboundPeeker).PerformInboundPeek of RoomserverInternalAPI.InboundPeeker.
func (r *RoomserverInternalAPI) PerformInboundPeek(ctx context.Context, req *api.PerformInboundPeekRequest, res *api.PerformInboundPeekResponse) error {
	panic("unimplemented")
}

// PerformJoin implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Joiner).PerformJoin of RoomserverInternalAPI.Joiner.
func (r *RoomserverInternalAPI) PerformJoin(ctx context.Context, req *api.PerformJoinRequest) (roomID string, joinedVia spec.ServerName, err error) {
	panic("unimplemented")
}

// PerformPeek implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Peeker).PerformPeek of RoomserverInternalAPI.Peeker.
func (r *RoomserverInternalAPI) PerformPeek(ctx context.Context, req *api.PerformPeekRequest) (roomID string, err error) {
	panic("unimplemented")
}

// PerformPublish implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Publisher).PerformPublish of RoomserverInternalAPI.Publisher.
func (r *RoomserverInternalAPI) PerformPublish(ctx context.Context, req *api.PerformPublishRequest) error {
	panic("unimplemented")
}

// PerformRoomUpgrade implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Upgrader).PerformRoomUpgrade of RoomserverInternalAPI.Upgrader.
func (r *RoomserverInternalAPI) PerformRoomUpgrade(ctx context.Context, roomID string, userID spec.UserID, roomVersion gomatrixserverlib.RoomVersion) (newRoomID string, err error) {
	panic("unimplemented")
}

// PerformUnpeek implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Unpeeker).PerformUnpeek of RoomserverInternalAPI.Unpeeker.
func (r *RoomserverInternalAPI) PerformUnpeek(ctx context.Context, roomID string, userID string, deviceID string) error {
	panic("unimplemented")
}

// QueryAuthChain implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryAuthChain of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryAuthChain(ctx context.Context, req *api.QueryAuthChainRequest, res *api.QueryAuthChainResponse) error {
	panic("unimplemented")
}

// QueryBulkStateContent implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryBulkStateContent of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryBulkStateContent(ctx context.Context, req *api.QueryBulkStateContentRequest, res *api.QueryBulkStateContentResponse) error {
	panic("unimplemented")
}

// QueryCurrentState implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryCurrentState of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryCurrentState(ctx context.Context, req *api.QueryCurrentStateRequest, res *api.QueryCurrentStateResponse) error {
	panic("unimplemented")
}

// QueryEventsByID implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryEventsByID of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryEventsByID(ctx context.Context, req *api.QueryEventsByIDRequest, res *api.QueryEventsByIDResponse) error {
	panic("unimplemented")
}

// QueryKnownUsers implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryKnownUsers of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryKnownUsers(ctx context.Context, req *api.QueryKnownUsersRequest, res *api.QueryKnownUsersResponse) error {
	panic("unimplemented")
}

// QueryLatestEventsAndState implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryLatestEventsAndState of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryLatestEventsAndState(ctx context.Context, req *api.QueryLatestEventsAndStateRequest, res *api.QueryLatestEventsAndStateResponse) error {
	panic("unimplemented")
}

// QueryLeftUsers implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryLeftUsers of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryLeftUsers(ctx context.Context, req *api.QueryLeftUsersRequest, res *api.QueryLeftUsersResponse) error {
	panic("unimplemented")
}

// QueryMembershipAtEvent implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryMembershipAtEvent of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryMembershipAtEvent(ctx context.Context, roomID spec.RoomID, eventIDs []string, senderID spec.SenderID) (map[string]*types.HeaderedEvent, error) {
	panic("unimplemented")
}

// QueryMembershipForSenderID implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryMembershipForSenderID of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryMembershipForSenderID(ctx context.Context, roomID spec.RoomID, senderID spec.SenderID, res *api.QueryMembershipForUserResponse) error {
	panic("unimplemented")
}

// QueryMembershipForUser implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryMembershipForUser of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryMembershipForUser(ctx context.Context, req *api.QueryMembershipForUserRequest, res *api.QueryMembershipForUserResponse) error {
	panic("unimplemented")
}

// QueryMembershipsForRoom implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryMembershipsForRoom of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryMembershipsForRoom(ctx context.Context, req *api.QueryMembershipsForRoomRequest, res *api.QueryMembershipsForRoomResponse) error {
	panic("unimplemented")
}

// QueryMissingEvents implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryMissingEvents of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryMissingEvents(ctx context.Context, req *api.QueryMissingEventsRequest, res *api.QueryMissingEventsResponse) error {
	panic("unimplemented")
}

// QueryNextRoomHierarchyPage implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryNextRoomHierarchyPage of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryNextRoomHierarchyPage(ctx context.Context, walker api.RoomHierarchyWalker, limit int) ([]fclient.RoomHierarchyRoom, *api.RoomHierarchyWalker, error) {
	panic("unimplemented")
}

// QueryPublishedRooms implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryPublishedRooms of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryPublishedRooms(ctx context.Context, req *api.QueryPublishedRoomsRequest, res *api.QueryPublishedRoomsResponse) error {
	panic("unimplemented")
}

// QueryRestrictedJoinAllowed implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryRestrictedJoinAllowed of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryRestrictedJoinAllowed(ctx context.Context, roomID spec.RoomID, senderID spec.SenderID) (string, error) {
	panic("unimplemented")
}

// QueryRoomInfo implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryRoomInfo of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryRoomInfo(ctx context.Context, roomID spec.RoomID) (*types.RoomInfo, error) {
	panic("unimplemented")
}

// QueryRoomVersionForRoom implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryRoomVersionForRoom of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryRoomVersionForRoom(ctx context.Context, roomID string) (gomatrixserverlib.RoomVersion, error) {
	panic("unimplemented")
}

// QueryRoomsForUser implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryRoomsForUser of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryRoomsForUser(ctx context.Context, userID spec.UserID, desiredMembership string) ([]spec.RoomID, error) {
	panic("unimplemented")
}

// QuerySenderIDForUser implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QuerySenderIDForUser of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QuerySenderIDForUser(ctx context.Context, roomID spec.RoomID, userID spec.UserID) (*spec.SenderID, error) {
	panic("unimplemented")
}

// QueryServerAllowedToSeeEvent implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryServerAllowedToSeeEvent of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryServerAllowedToSeeEvent(ctx context.Context, serverName spec.ServerName, eventID string, roomID string) (allowed bool, err error) {
	panic("unimplemented")
}

// QueryServerBannedFromRoom implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryServerBannedFromRoom of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryServerBannedFromRoom(ctx context.Context, req *api.QueryServerBannedFromRoomRequest, res *api.QueryServerBannedFromRoomResponse) error {
	panic("unimplemented")
}

// QueryServerJoinedToRoom implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryServerJoinedToRoom of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryServerJoinedToRoom(ctx context.Context, req *api.QueryServerJoinedToRoomRequest, res *api.QueryServerJoinedToRoomResponse) error {
	panic("unimplemented")
}

// QuerySharedUsers implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QuerySharedUsers of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QuerySharedUsers(ctx context.Context, req *api.QuerySharedUsersRequest, res *api.QuerySharedUsersResponse) error {
	panic("unimplemented")
}

// QueryStateAfterEvents implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryStateAfterEvents of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryStateAfterEvents(ctx context.Context, req *api.QueryStateAfterEventsRequest, res *api.QueryStateAfterEventsResponse) error {
	panic("unimplemented")
}

// QueryStateAndAuthChain implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryStateAndAuthChain of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryStateAndAuthChain(ctx context.Context, req *api.QueryStateAndAuthChainRequest, res *api.QueryStateAndAuthChainResponse) error {
	panic("unimplemented")
}

// QueryUserIDForSender implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).QueryUserIDForSender of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) QueryUserIDForSender(ctx context.Context, roomID spec.RoomID, senderID spec.SenderID) (*spec.UserID, error) {
	panic("unimplemented")
}

// RestrictedRoomJoinInfo implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).RestrictedRoomJoinInfo of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) RestrictedRoomJoinInfo(ctx context.Context, roomID spec.RoomID, senderID spec.SenderID, localServerName spec.ServerName) (*gomatrixserverlib.RestrictedRoomJoinInfo, error) {
	panic("unimplemented")
}

// UserJoinedToRoom implements api.RoomserverInternalAPI.
// Subtle: this method shadows the method (*Queryer).UserJoinedToRoom of RoomserverInternalAPI.Queryer.
func (r *RoomserverInternalAPI) UserJoinedToRoom(ctx context.Context, roomID types.RoomNID, senderID spec.SenderID) (bool, error) {
	panic("unimplemented")
}

func NewRoomserverAPI(
	processContext *process.ProcessContext, dendriteCfg *config.Dendrite, roomserverDB storage.Database,
	js nats.JetStreamContext, nc *nats.Conn, caches caching.RoomServerCaches, enableMetrics bool,
) *RoomserverInternalAPI {
	var perspectiveServerNames []spec.ServerName
	for _, kp := range dendriteCfg.FederationAPI.KeyPerspectives {
		perspectiveServerNames = append(perspectiveServerNames, kp.ServerName)
	}

	serverACLs := acls.NewServerACLs(roomserverDB)
	producer := &producers.RoomEventProducer{
		Topic:     string(dendriteCfg.Global.JetStream.Prefixed(jetstream.OutputRoomEvent)),
		JetStream: js,
		ACLs:      serverACLs,
	}
	a := &RoomserverInternalAPI{
		ProcessContext:         processContext,
		DB:                     roomserverDB,
		Cfg:                    dendriteCfg,
		Cache:                  caches,
		ServerName:             dendriteCfg.Global.ServerName,
		PerspectiveServerNames: perspectiveServerNames,
		InputRoomEventTopic:    dendriteCfg.Global.JetStream.Prefixed(jetstream.InputRoomEvent),
		OutputProducer:         producer,
		JetStream:              js,
		NATSClient:             nc,
		Durable:                dendriteCfg.Global.JetStream.Durable("RoomserverInputConsumer"),
		ServerACLs:             serverACLs,
		enableMetrics:          enableMetrics,
		defaultRoomVersion:     dendriteCfg.RoomServer.DefaultRoomVersion,
		// perform-er structs + queryer struct get initialised when we have a federation sender to use
	}
	return a
}

// SetFederationInputAPI passes in a federation input API reference so that we can
// avoid the chicken-and-egg problem of both the roomserver input API and the
// federation input API being interdependent.
func (r *RoomserverInternalAPI) SetFederationAPI(fsAPI fsAPI.RoomserverFederationAPI, keyRing *gomatrixserverlib.KeyRing) {
	r.fsAPI = fsAPI
	r.KeyRing = keyRing

	r.Queryer = &query.Queryer{
		DB:                r.DB,
		Cache:             r.Cache,
		IsLocalServerName: r.Cfg.Global.IsLocalServerName,
		ServerACLs:        r.ServerACLs,
		Cfg:               r.Cfg,
		FSAPI:             fsAPI,
	}

	r.Inputer = &input.Inputer{
		Cfg:                 &r.Cfg.RoomServer,
		ProcessContext:      r.ProcessContext,
		DB:                  r.DB,
		InputRoomEventTopic: r.InputRoomEventTopic,
		OutputProducer:      r.OutputProducer,
		JetStream:           r.JetStream,
		NATSClient:          r.NATSClient,
		Durable:             nats.Durable(r.Durable),
		ServerName:          r.ServerName,
		SigningIdentity:     r.SigningIdentityFor,
		FSAPI:               fsAPI,
		RSAPI:               r,
		KeyRing:             keyRing,
		ACLs:                r.ServerACLs,
		Queryer:             r.Queryer,
		EnableMetrics:       r.enableMetrics,
	}
	r.Inviter = &perform.Inviter{
		DB:      r.DB,
		Cfg:     &r.Cfg.RoomServer,
		FSAPI:   r.fsAPI,
		RSAPI:   r,
		Inputer: r.Inputer,
	}
	r.Joiner = &perform.Joiner{
		Cfg:     &r.Cfg.RoomServer,
		DB:      r.DB,
		FSAPI:   r.fsAPI,
		RSAPI:   r,
		Inputer: r.Inputer,
		Queryer: r.Queryer,
	}
	r.Peeker = &perform.Peeker{
		ServerName: r.ServerName,
		Cfg:        &r.Cfg.RoomServer,
		DB:         r.DB,
		FSAPI:      r.fsAPI,
		Inputer:    r.Inputer,
	}
	r.InboundPeeker = &perform.InboundPeeker{
		DB:      r.DB,
		Inputer: r.Inputer,
	}
	r.Unpeeker = &perform.Unpeeker{
		ServerName: r.ServerName,
		Cfg:        &r.Cfg.RoomServer,
		FSAPI:      r.fsAPI,
		Inputer:    r.Inputer,
	}
	r.Leaver = &perform.Leaver{
		Cfg:     &r.Cfg.RoomServer,
		DB:      r.DB,
		FSAPI:   r.fsAPI,
		RSAPI:   r,
		Inputer: r.Inputer,
	}
	r.Publisher = &perform.Publisher{
		DB: r.DB,
	}
	r.Backfiller = &perform.Backfiller{
		IsLocalServerName: r.Cfg.Global.IsLocalServerName,
		DB:                r.DB,
		FSAPI:             r.fsAPI,
		Querier:           r.Queryer,
		KeyRing:           r.KeyRing,
		// Perspective servers are trusted to not lie about server keys, so we will also
		// prefer these servers when backfilling (assuming they are in the room) rather
		// than trying random servers
		PreferServers: r.PerspectiveServerNames,
	}
	r.Forgetter = &perform.Forgetter{
		DB: r.DB,
	}
	r.Upgrader = &perform.Upgrader{
		Cfg:    &r.Cfg.RoomServer,
		URSAPI: r,
	}
	r.Admin = &perform.Admin{
		DB:      r.DB,
		Cfg:     &r.Cfg.RoomServer,
		Inputer: r.Inputer,
		Queryer: r.Queryer,
		Leaver:  r.Leaver,
	}
	r.Creator = &perform.Creator{
		DB:    r.DB,
		Cfg:   &r.Cfg.RoomServer,
		RSAPI: r,
	}

	if err := r.Inputer.Start(); err != nil {
		logrus.WithError(err).Panic("failed to start roomserver input API")
	}
}

func (r *RoomserverInternalAPI) SetUserAPI(userAPI userapi.RoomserverUserAPI) {
	r.Leaver.UserAPI = userAPI
	r.Inputer.UserAPI = userAPI
}

func (r *RoomserverInternalAPI) SetAppserviceAPI(asAPI asAPI.AppServiceInternalAPI) {
	r.asAPI = asAPI
}

func (r *RoomserverInternalAPI) DefaultRoomVersion() gomatrixserverlib.RoomVersion {
	return r.defaultRoomVersion
}

func (r *RoomserverInternalAPI) IsKnownRoom(ctx context.Context, roomID spec.RoomID) (bool, error) {
	return r.Inviter.IsKnownRoom(ctx, roomID)
}

func (r *RoomserverInternalAPI) StateQuerier() gomatrixserverlib.StateQuerier {
	return r.Inviter.StateQuerier()
}

func (r *RoomserverInternalAPI) HandleInvite(
	ctx context.Context, inviteEvent *types.HeaderedEvent,
) error {
	outputEvents, err := r.Inviter.ProcessInviteMembership(ctx, inviteEvent)
	if err != nil {
		return err
	}
	return r.OutputProducer.ProduceRoomEvents(inviteEvent.RoomID().String(), outputEvents)
}

func (r *RoomserverInternalAPI) PerformCreateRoom(
	ctx context.Context, userID spec.UserID, roomID spec.RoomID, createRequest *api.PerformCreateRoomRequest,
) (string, *util.JSONResponse) {
	return r.Creator.PerformCreateRoom(ctx, userID, roomID, createRequest)
}

func (r *RoomserverInternalAPI) PerformInvite(
	ctx context.Context,
	req *api.PerformInviteRequest,
) error {
	return r.Inviter.PerformInvite(ctx, req)
}

func (r *RoomserverInternalAPI) PerformLeave(
	ctx context.Context,
	req *api.PerformLeaveRequest,
	res *api.PerformLeaveResponse,
) error {
	outputEvents, err := r.Leaver.PerformLeave(ctx, req, res)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	if len(outputEvents) == 0 {
		return nil
	}
	return r.OutputProducer.ProduceRoomEvents(req.RoomID, outputEvents)
}

func (r *RoomserverInternalAPI) PerformForget(
	ctx context.Context,
	req *api.PerformForgetRequest,
	resp *api.PerformForgetResponse,
) error {
	return r.Forgetter.PerformForget(ctx, req, resp)
}

// GetOrCreateUserRoomPrivateKey gets the user room key for the specified user. If no key exists yet, a new one is created.
func (r *RoomserverInternalAPI) GetOrCreateUserRoomPrivateKey(ctx context.Context, userID spec.UserID, roomID spec.RoomID) (ed25519.PrivateKey, error) {
	key, err := r.DB.SelectUserRoomPrivateKey(ctx, userID, roomID)
	if err != nil {
		return nil, err
	}
	// no key found, create one
	if len(key) == 0 {
		_, key, err = ed25519.GenerateKey(nil)
		if err != nil {
			return nil, err
		}
		key, err = r.DB.InsertUserRoomPrivatePublicKey(ctx, userID, roomID, key)
		if err != nil {
			return nil, err
		}
	}
	return key, nil
}

func (r *RoomserverInternalAPI) StoreUserRoomPublicKey(ctx context.Context, senderID spec.SenderID, userID spec.UserID, roomID spec.RoomID) error {
	pubKeyBytes, err := senderID.RawBytes()
	if err != nil {
		return err
	}
	_, err = r.DB.InsertUserRoomPublicKey(ctx, userID, roomID, ed25519.PublicKey(pubKeyBytes))
	return err
}

func (r *RoomserverInternalAPI) SigningIdentityFor(ctx context.Context, roomID spec.RoomID, senderID spec.UserID) (fclient.SigningIdentity, error) {
	roomVersion, ok := r.Cache.GetRoomVersion(roomID.String())
	if !ok {
		roomInfo, err := r.DB.RoomInfo(ctx, roomID.String())
		if err != nil {
			return fclient.SigningIdentity{}, err
		}
		if roomInfo != nil {
			roomVersion = roomInfo.RoomVersion
		}
	}
	if roomVersion == gomatrixserverlib.RoomVersionPseudoIDs {
		privKey, err := r.GetOrCreateUserRoomPrivateKey(ctx, senderID, roomID)
		if err != nil {
			return fclient.SigningIdentity{}, err
		}
		return fclient.SigningIdentity{
			PrivateKey: privKey,
			KeyID:      "ed25519:1",
			ServerName: spec.ServerName(spec.SenderIDFromPseudoIDKey(privKey)),
		}, nil
	}
	identity, err := r.Cfg.Global.SigningIdentityFor(senderID.Domain())
	if err != nil {
		return fclient.SigningIdentity{}, err
	}
	return *identity, err
}

func (r *RoomserverInternalAPI) AssignRoomNID(ctx context.Context, roomID spec.RoomID, roomVersion gomatrixserverlib.RoomVersion) (roomNID types.RoomNID, err error) {
	return r.DB.AssignRoomNID(ctx, roomID, roomVersion)
}
