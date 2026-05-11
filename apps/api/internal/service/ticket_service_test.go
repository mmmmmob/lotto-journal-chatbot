package service

import "testing"

func TestParseTicketInput_SpacesAroundX(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		number string
		qty    int
		typeID string
	}{
		{
			name:   "space before x",
			input:  "144333 x2",
			number: "144333",
			qty:    2,
			typeID: "L6",
		},
		{
			name:   "space before and after x",
			input:  "122222 x 3",
			number: "122222",
			qty:    3,
			typeID: "L6",
		},
		{
			name:   "space after x",
			input:  "333333x 9",
			number: "333333",
			qty:    9,
			typeID: "L6",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tickets, invalid := ParseTicketInput(tc.input)

			if len(invalid) != 0 {
				t.Fatalf("expected no invalid tokens, got %v", invalid)
			}
			if len(tickets) != 1 {
				t.Fatalf("expected 1 ticket, got %d (%v)", len(tickets), tickets)
			}

			got := tickets[0]
			if got.Number != tc.number || got.Quantity != tc.qty || got.Type != tc.typeID {
				t.Fatalf("unexpected parsed ticket: got=%+v want={Number:%s Quantity:%d Type:%s}", got, tc.number, tc.qty, tc.typeID)
			}
		})
	}
}

func TestParseTicketInput_UnicodeWhitespaceAndMultiplicationX(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		number string
		qty    int
		typeID string
	}{
		{
			name:   "multiplication sign without spaces",
			input:  "123456×2",
			number: "123456",
			qty:    2,
			typeID: "L6",
		},
		{
			name:   "non breaking spaces and multiplication sign",
			input:  "123456\u00A0×\u00A02",
			number: "123456",
			qty:    2,
			typeID: "L6",
		},
		{
			name:   "full width x",
			input:  "456ｘ4",
			number: "456",
			qty:    4,
			typeID: "N3",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tickets, invalid := ParseTicketInput(tc.input)

			if len(invalid) != 0 {
				t.Fatalf("expected no invalid tokens, got %v", invalid)
			}
			if len(tickets) != 1 {
				t.Fatalf("expected 1 ticket, got %d (%v)", len(tickets), tickets)
			}

			got := tickets[0]
			if got.Number != tc.number || got.Quantity != tc.qty || got.Type != tc.typeID {
				t.Fatalf("unexpected parsed ticket: got=%+v want={Number:%s Quantity:%d Type:%s}", got, tc.number, tc.qty, tc.typeID)
			}
		})
	}
}
