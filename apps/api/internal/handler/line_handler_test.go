package handler

import (
	"strings"
	"testing"
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
	t.Run("with display name", func(t *testing.T) {
		msg := buildWelcomeMessage("Ploy")
		if !strings.Contains(msg, "สวัสดีคุณ Ploy") {
			t.Fatalf("expected personalized greeting, got %q", msg)
		}
		if !strings.Contains(msg, "ยินดีต้อนรับสู่ Lotto Journal") {
			t.Fatalf("expected welcome body, got %q", msg)
		}
	})

	t.Run("without display name", func(t *testing.T) {
		msg := buildWelcomeMessage("")
		if !strings.Contains(msg, "สวัสดี!") {
			t.Fatalf("expected generic greeting, got %q", msg)
		}
		if !strings.Contains(msg, "ยินดีต้อนรับสู่ Lotto Journal") {
			t.Fatalf("expected welcome body, got %q", msg)
		}
	})
}
