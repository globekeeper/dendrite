package pushrules

// A RuleSet contains all the various push rules for an
// account. Listed in decreasing order of priority.
type RuleSet struct {
	Override  []*Rule `json:"override,omitempty"`
	Content   []*Rule `json:"content,omitempty"`
	Room      []*Rule `json:"room,omitempty"`
	Sender    []*Rule `json:"sender,omitempty"`
	Underride []*Rule `json:"underride,omitempty"`
}

type Rule struct {
	RuleID  string    `json:"rule_id"` // Required. For SenderKind, this is the MXID of the matching sender.
	Default bool      `json:"default"` // Required.
	Enabled bool      `json:"enabled"` // Required.
	Actions []*Action `json:"actions"` // Required.

	Conditions []*Condition `json:"conditions,omitempty"` // Only allowed for OverrideKind and UnderrideKind.
	Pattern    string       `json:"pattern,omitempty"`    // Required for ContentKind.

	// kind and scope are part of the push rules request/responses,
	// but not of the core data model.
}

type Kind string

const (
	UnknownKind   Kind = ""
	OverrideKind  Kind = "override"
	ContentKind   Kind = "content"
	RoomKind      Kind = "room"
	SenderKind    Kind = "sender"
	UnderrideKind Kind = "underride"
	// The server-default rules have the lowest priority.
)
