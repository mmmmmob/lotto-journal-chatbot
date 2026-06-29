package localization

import (
	"testing"

	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
)

func TestGetDictionary(t *testing.T) {
	// 1. Test supported languages
	t.Run("TH dictionary", func(t *testing.T) {
		dict := GetDictionary("th")
		if dict == nil {
			t.Fatal("expected TH dictionary, got nil")
		}
		if dict.LanguageSwitched != "เปลี่ยนภาษาเป็นภาษาไทยเรียบร้อย 🇹🇭" {
			t.Errorf("unexpected TH string, got %q", dict.LanguageSwitched)
		}
		if dict.WelcomeGreetingGeneric != "👋 สวัสดี!" {
			t.Errorf("unexpected TH generic greeting, got %q", dict.WelcomeGreetingGeneric)
		}
		if dict.WelcomeGreetingPersonal != "👋 สวัสดีคุณ %s!" {
			t.Errorf("unexpected TH personal greeting, got %q", dict.WelcomeGreetingPersonal)
		}
	})

	t.Run("EN dictionary", func(t *testing.T) {
		dict := GetDictionary("en")
		if dict == nil {
			t.Fatal("expected EN dictionary, got nil")
		}
		if dict.LanguageSwitched != "Language changed to English 🇺🇸" {
			t.Errorf("unexpected EN string, got %q", dict.LanguageSwitched)
		}
		if dict.WelcomeGreetingGeneric != "👋 Hello!" {
			t.Errorf("unexpected EN generic greeting, got %q", dict.WelcomeGreetingGeneric)
		}
		if dict.WelcomeGreetingPersonal != "👋 Hello %s!" {
			t.Errorf("unexpected EN personal greeting, got %q", dict.WelcomeGreetingPersonal)
		}
	})

	// 2. Test fallback
	t.Run("fallback to EN", func(t *testing.T) {
		dict := GetDictionary("invalid-lang")
		if dict == nil {
			t.Fatal("expected fallback dictionary, got nil")
		}
		if dict.LanguageSwitched != "Language changed to English 🇺🇸" {
			t.Errorf("expected EN fallback, got %q", dict.LanguageSwitched)
		}
	})

	t.Run("case insensitivity matches", func(t *testing.T) {
		dict := GetDictionary("TH")
		if dict.LanguageSwitched != "เปลี่ยนภาษาเป็นภาษาไทยเรียบร้อย 🇹🇭" {
			t.Errorf("expected case-insensitive matching for TH, got %q", dict.LanguageSwitched)
		}
	})
}

func TestGetQuickReplies(t *testing.T) {
	t.Run("Thai Quick Replies", func(t *testing.T) {
		qr := GetQuickReplies("th")
		if qr == nil || len(qr.Items) == 0 {
			t.Fatal("expected Thai quick replies, got nil or empty")
		}
		if len(qr.Items) != 4 {
			t.Errorf("expected 4 quick reply items, got %d", len(qr.Items))
		}
		// check first item
		firstAction, ok := qr.Items[0].Action.(*messaging_api.MessageAction)
		if !ok || firstAction.Label != "📝 โพย" || firstAction.Text != "โพย" {
			t.Errorf("unexpected first item: %+v", qr.Items[0])
		}
		// check last item (change language)
		lastAction, ok := qr.Items[3].Action.(*messaging_api.MessageAction)
		if !ok || lastAction.Label != "🇺🇸 Change Language" || lastAction.Text != "english" {
			t.Errorf("unexpected last item: %+v", qr.Items[3])
		}
	})

	t.Run("English Quick Replies", func(t *testing.T) {
		qr := GetQuickReplies("en")
		if qr == nil || len(qr.Items) == 0 {
			t.Fatal("expected English quick replies, got nil or empty")
		}
		if len(qr.Items) != 4 {
			t.Errorf("expected 4 quick reply items, got %d", len(qr.Items))
		}
		// check first item
		firstAction, ok := qr.Items[0].Action.(*messaging_api.MessageAction)
		if !ok || firstAction.Label != "📝 List" || firstAction.Text != "list" {
			t.Errorf("unexpected first item: %+v", qr.Items[0])
		}
		// check last item (change language)
		lastAction, ok := qr.Items[3].Action.(*messaging_api.MessageAction)
		if !ok || lastAction.Label != "🇹🇭 เปลี่ยนภาษา" || lastAction.Text != "ไทย" {
			t.Errorf("unexpected last item: %+v", qr.Items[3])
		}
	})
}
