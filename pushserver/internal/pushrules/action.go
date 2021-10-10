package pushrules

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Action struct {
	Kind  ActionKind  `json:"-"` // Custom encoding in JSON.
	Tweak TweakKey    `json:"-"` // Custom encoding in JSON.
	Value interface{} `json:"value,omitempty"`
}

func (a *Action) MarshalJSON() ([]byte, error) {
	if a.Value == nil {
		return json.Marshal(a.Kind)
	}

	if a.Kind != SetTweakAction {
		return nil, fmt.Errorf("only set_tweak actions may have a value, but got kind %q", a.Kind)
	}

	return json.Marshal(map[string]interface{}{
		string(a.Kind): a.Tweak,
		"value":        a.Value,
	})
}

func (a *Action) UnmarshalJSON(bs []byte) error {
	if bytes.HasPrefix(bs, []byte("\"")) {
		return json.Unmarshal(bs, &a.Kind)
	}

	var raw struct {
		SetTweak TweakKey    `json:"set_tweak"`
		Value    interface{} `json:"value"`
	}
	if err := json.Unmarshal(bs, &raw); err != nil {
		return err
	}
	if raw.SetTweak == UnknownTweak {
		return fmt.Errorf("got unknown action JSON: %s", string(bs))
	}
	a.Kind = SetTweakAction
	a.Tweak = raw.SetTweak
	a.Value = raw.Value

	return nil
}

type ActionKind string

const (
	UnknownAction    ActionKind = ""
	NotifyAction     ActionKind = "notify"
	DontNotifyAction ActionKind = "dont_notify"
	CoalesceAction   ActionKind = "coalesce"
	SetTweakAction   ActionKind = "set_tweak"
)

type TweakKey string

const (
	UnknownTweak   TweakKey = ""
	SoundTweak     TweakKey = "sound"
	HighlightTweak TweakKey = "highlight"
)
