package pubsub

import (
	"testing"

	"github.com/diamondburned/facechat/backend/facechat"
)

func TestTypeName(t *testing.T) {
	if name := typeName(facechat.Account{}); name != "Account" {
		t.Fatal("Unexpected type name for facechat.Account:", name)
	}
}
