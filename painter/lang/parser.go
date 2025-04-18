package lang

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/maxnetyaga/software-architecture-lab3/painter"
)

type Parser struct {
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()
		op, err := parse(commandLine)
		if err != nil {
			return nil, fmt.Errorf("failed to parse command '%s': %w", commandLine, err)
		}
		if op != nil {
			res = append(res, op)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return res, nil
}

func parse(commandLine string) (painter.Operation, error) {
	fields := strings.Fields(commandLine)
	if len(fields) == 0 {
		return nil, nil
	}

	instruction := fields[0]
	args := fields[1:]

	switch instruction {
	case "white":
		if len(args) != 0 {
			return nil, fmt.Errorf("unexpected arguments for white command")
		}
		return painter.WhiteOp{}, nil
	case "green":
		if len(args) != 0 {
			return nil, fmt.Errorf("unexpected arguments for green command")
		}
		return painter.GreenOp{}, nil
	case "update":
		if len(args) != 0 {
			return nil, fmt.Errorf("unexpected arguments for update command")
		}
		return painter.UpdateOp, nil
	case "bgrect":
		if len(args) != 4 {
			return nil, fmt.Errorf("bgrect command requires 4 arguments")
		}
		x1, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid argument for bgrect: %w", err)
		}
		y1, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid argument for bgrect: %w", err)
		}
		x2, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid argument for bgrect: %w", err)
		}
		y2, err := strconv.ParseFloat(args[3], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid argument for bgrect: %w", err)
		}
		return painter.BgRectOp{X1: x1, Y1: y1, X2: x2, Y2: y2}, nil
	case "figure":
		if len(args) != 2 {
			return nil, fmt.Errorf("figure command requires 2 arguments")
		}
		x, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid argument for figure: %w", err)
		}
		y, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid argument for figure: %w", err)
		}
		return painter.FigureOp{X: x, Y: y}, nil
	case "move":
		if len(args) != 2 {
			return nil, fmt.Errorf("move command requires 2 arguments")
		}
		x, err := strconv.ParseFloat(args[0], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid argument for move: %w", err)
		}
		y, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid argument for move: %w", err)
		}
		return painter.MoveOp{X: x, Y: y}, nil
	case "reset":
		if len(args) != 0 {
			return nil, fmt.Errorf("unexpected arguments for reset command")
		}
		return painter.ResetOp{}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", instruction)
	}
}