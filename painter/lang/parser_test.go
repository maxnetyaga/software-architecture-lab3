package lang

import (
	"strings"
	"testing"
	"reflect"

	"github.com/maxnetyaga/software-architecture-lab3/painter"
)

func TestParser_Parse(t *testing.T) {
	var p Parser

	tests := []struct {
		name string
		input string
		expected []painter.Operation
		expectError bool
	}{
		{
			name: "empty input",
			input: "",
			expected: []painter.Operation{},
			expectError: false,
		},
		{
			name: "whitespace only",
			input: "   \n \t ",
			expected: []painter.Operation{},
			expectError: false,
		},
		{
			name: "valid white command",
			input: "white",
			expected: []painter.Operation{painter.WhiteOp{}},
			expectError: false,
		},
		{
			name: "valid green command",
			input: "green",
			expected: []painter.Operation{painter.GreenOp{}},
			expectError: false,
		},
		{
			name: "valid update command",
			input: "update",
			expected: []painter.Operation{painter.UpdateOp},
			expectError: false,
		},
		{
			name: "valid bgrect command",
			input: "bgrect 0.1 0.2 0.8 0.9",
			expected: []painter.Operation{painter.BgRectOp{X1: 0.1, Y1: 0.2, X2: 0.8, Y2: 0.9}},
			expectError: false,
		},
		{
			name: "valid figure command",
			input: "figure 0.5 0.5",
			expected: []painter.Operation{painter.FigureOp{X: 0.5, Y: 0.5}},
			expectError: false,
		},
		{
			name: "valid move command",
			input: "move 0.01 0.02",
			expected: []painter.Operation{painter.MoveOp{X: 0.01, Y: 0.02}},
			expectError: false,
		},
		{
			name: "valid reset command",
			input: "reset",
			expected: []painter.Operation{painter.ResetOp{}},
			expectError: false,
		},
		{
			name: "multiple valid commands",
			input: "white\nfigure 0.5 0.5\nupdate\ngreen\nreset",
			expected: []painter.Operation{
				painter.WhiteOp{},
				painter.FigureOp{X: 0.5, Y: 0.5},
				painter.UpdateOp,
				painter.GreenOp{},
				painter.ResetOp{},
			},
			expectError: false,
		},
		{
			name: "unknown command",
			input: "fill red",
			expected: nil,
			expectError: true,
		},
		{
			name: "bgrect missing arguments",
			input: "bgrect 0.1 0.2",
			expected: nil,
			expectError: true,
		},
		{
			name: "figure extra arguments",
			input: "figure 0.5 0.5 0.6",
			expected: nil,
			expectError: true,
		},
		{
			name: "move invalid argument type",
			input: "move abc 0.1",
			expected: nil,
			expectError: true,
		},
		{
			name: "update with argument",
			input: "update 1",
			expected: nil,
			expectError: true,
		},
		{
			name: "command with leading/trailing spaces",
			input: "  white  ",
			expected: []painter.Operation{painter.WhiteOp{}},
			expectError: false,
		},
		{
			name: "command with multiple spaces between parts",
			input: "figure   0.5   0.5",
			expected: []painter.Operation{painter.FigureOp{X: 0.5, Y: 0.5}},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			ops, err := p.Parse(reader)

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, Got error: %v", tt.expectError, err)
				return
			}

			if err == nil {
				if len(ops) != len(tt.expected) {
					t.Fatalf("Expected %d operations, got %d", len(tt.expected), len(ops))
				}

				for i := range ops {
					switch receivedOp := ops[i].(type) {
					case painter.WhiteOp:
						if _, ok := tt.expected[i].(painter.WhiteOp); !ok {
							t.Errorf("Operation type mismatch at index %d. Expected type: painter.WhiteOp, Got received type: %T", i, receivedOp)
						}
					case painter.GreenOp:
						if _, ok := tt.expected[i].(painter.GreenOp); !ok {
							t.Errorf("Operation type mismatch at index %d. Expected type: painter.GreenOp, Got received type: %T", i, receivedOp)
						}
					case painter.BgRectOp:
						expectedOp, ok := tt.expected[i].(painter.BgRectOp)
						if !ok {
							t.Errorf("Operation type mismatch at index %d. Expected type: painter.BgRectOp, Got received type: %T", i, receivedOp)
						} else if receivedOp != expectedOp {
							t.Errorf("BgRectOp mismatch at index %d. Expected: %v, Got: %v", i, expectedOp, receivedOp)
						}
					case painter.FigureOp:
						expectedOp, ok := tt.expected[i].(painter.FigureOp)
						if !ok {
							t.Errorf("Operation type mismatch at index %d. Expected type: painter.FigureOp, Got received type: %T", i, receivedOp)
						} else if receivedOp != expectedOp {
							t.Errorf("FigureOp mismatch at index %d. Expected: %v, Got: %v", i, expectedOp, receivedOp)
						}
					case painter.MoveOp:
						expectedOp, ok := tt.expected[i].(painter.MoveOp)
						if !ok {
							t.Errorf("Operation type mismatch at index %d. Expected type: painter.MoveOp, Got received type: %T", i, receivedOp)
						} else if receivedOp != expectedOp {
							t.Errorf("MoveOp mismatch at index %d. Expected: %v, Got: %v", i, expectedOp, receivedOp)
						}
					case painter.ResetOp:
						if _, ok := tt.expected[i].(painter.ResetOp); !ok {
							t.Errorf("Operation type mismatch at index %d. Expected type: painter.ResetOp, Got received type: %T", i, receivedOp)
						}
					default:
						if !reflect.DeepEqual(receivedOp, painter.UpdateOp) {
							 t.Errorf("Unexpected operation type or value at index %d. Expected: %v, Got: %v (Type: %T)", i, tt.expected[i], receivedOp, receivedOp)
						} else {
							if tt.expected[i] != painter.UpdateOp {
								t.Errorf("Operation mismatch at index %d. Expected: %v, Got: %v (both are UpdateOp, but expected was not)", i, tt.expected[i], receivedOp)
							}
						}
					}
				}
			}
		})
	}
}
