package handler

import (
	"strings"
	"testing"

	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func TestIsTicketListCmd(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{name: "exact", text: "โพย", want: true},
		{name: "leading and trailing spaces", text: "  โพย  ", want: true},
		{name: "internal ascii space", text: "โ พย", want: true},
		{name: "internal tab", text: "โ\tพย", want: true},
		{name: "internal non breaking space", text: "โ\u00A0พย", want: true},
		{name: "internal zero width space", text: "โ\u200Bพย", want: true},
		{name: "prefix text", text: "ขอโพย", want: false},
		{name: "suffix text", text: "โพยครับ", want: false},
		{name: "different word", text: "โพ", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isTicketListCmd(tc.text)
			if got != tc.want {
				t.Fatalf("isTicketListCmd(%q) = %v, want %v", tc.text, got, tc.want)
			}
		})
	}
}

func TestBuildWelcomeMessage(t *testing.T) {
	t.Run("Thai with display name", func(t *testing.T) {
		msg := buildWelcomeMessage("Ploy", "th", true)
		if !strings.Contains(msg, "สวัสดีคุณ Ploy") {
			t.Fatalf("expected personalized greeting, got %q", msg)
		}
		if !strings.Contains(msg, "ยินดีต้อนรับสู่ Lotto Journal") {
			t.Fatalf("expected welcome body, got %q", msg)
		}
		if !strings.Contains(msg, "เปลี่ยนเป็นภาษาอังกฤษ") {
			t.Fatalf("expected toggle hint, got %q", msg)
		}
	})

	t.Run("Thai without display name", func(t *testing.T) {
		msg := buildWelcomeMessage("", "th", true)
		if !strings.Contains(msg, "สวัสดี!") {
			t.Fatalf("expected generic greeting, got %q", msg)
		}
		if !strings.Contains(msg, "ยินดีต้อนรับสู่ Lotto Journal") {
			t.Fatalf("expected welcome body, got %q", msg)
		}
	})

	t.Run("Thai returning", func(t *testing.T) {
		msg := buildWelcomeMessage("Ploy", "th", false)
		if !strings.Contains(msg, "สวัสดีคุณ Ploy") {
			t.Fatalf("expected personalized greeting, got %q", msg)
		}
		if !strings.Contains(msg, "ยินดีต้อนรับกลับสู่ Lotto Journal") {
			t.Fatalf("expected returning welcome body, got %q", msg)
		}
		if strings.Contains(msg, "เปลี่ยนเป็นภาษาอังกฤษ") {
			t.Fatalf("did not expect toggle hint, got %q", msg)
		}
	})

	t.Run("English with display name", func(t *testing.T) {
		msg := buildWelcomeMessage("John", "en", true)
		if !strings.Contains(msg, "Hello John!") {
			t.Fatalf("expected personalized greeting, got %q", msg)
		}
		if !strings.Contains(msg, "Welcome to Lotto Journal") {
			t.Fatalf("expected welcome body, got %q", msg)
		}
		if !strings.Contains(msg, "switch to Thai language") {
			t.Fatalf("expected toggle hint, got %q", msg)
		}
	})

	t.Run("English returning", func(t *testing.T) {
		msg := buildWelcomeMessage("John", "en", false)
		if !strings.Contains(msg, "Hello John!") {
			t.Fatalf("expected personalized greeting, got %q", msg)
		}
		if !strings.Contains(msg, "Welcome back to Lotto Journal") {
			t.Fatalf("expected returning welcome body, got %q", msg)
		}
		if strings.Contains(msg, "switch to Thai language") {
			t.Fatalf("did not expect toggle hint, got %q", msg)
		}
	})
}

func TestEventReplyToken(t *testing.T) {
	tests := []struct {
		name  string
		event webhook.EventInterface
		want  string
	}{
		{
			name:  "follow event",
			event: webhook.FollowEvent{ReplyToken: "follow-token"},
			want:  "follow-token",
		},
		{
			name:  "message event",
			event: webhook.MessageEvent{ReplyToken: "message-token"},
			want:  "message-token",
		},
		{
			name:  "unfollow event",
			event: webhook.UnfollowEvent{},
			want:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := eventReplyToken(tc.event)
			if got != tc.want {
				t.Fatalf("eventReplyToken(%T) = %q, want %q", tc.event, got, tc.want)
			}
		})
	}
}
