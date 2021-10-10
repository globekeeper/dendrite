package pushrules

type Condition struct {
	Kind    ConditionKind `json:"kind"`              // Required.
	Key     string        `json:"key,omitempty"`     // Required for EventMatchCondition and SenderNotificationPermissionCondition.
	Pattern string        `json:"pattern,omitempty"` // Required for EventMatchCondition.
	Is      string        `json:"is,omitempty"`      // Required for RoomMemberCountCondition.
}

// ConditionKind represents a kind of condition.
//
// SPEC: Unrecognised conditions MUST NOT match any events,
// effectively making the push rule disabled.
type ConditionKind string

const (
	UnknownCondition                      ConditionKind = ""
	EventMatchCondition                   ConditionKind = "event_match"
	ContainsDisplayNameCondition          ConditionKind = "contains_display_name"
	RoomMemberCountCondition              ConditionKind = "room_member_count"
	SenderNotificationPermissionCondition ConditionKind = "sender_notification_permission"
)
