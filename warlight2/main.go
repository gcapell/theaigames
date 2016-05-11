package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

var (
	timebank, time_per_move, max_rounds, starting_armies, starting_pick_amount int

	your_bot, opponent_bot string

	starting_regions []int
)

func settings(p []string) string {
	switch p[0] {
	case "timebank":
		timebank = parseInt(p[1])
	case "time_per_move":
		time_per_move = parseInt(p[1])
	case "max_rounds":
		max_rounds = parseInt(p[1])
	case "your_bot":
		your_bot = p[1]
	case "opponent_bot":
		opponent_bot = p[1]
	case "starting_armies":
		starting_armies = parseInt(p[1])
	case "starting_regions":
		starting_regions = parse_ints(p[1:])
	case "starting_pick_amount  ":
		starting_pick_amount = parseInt(p[1])
	}
	return ""
}

func setup_map(p []string) {
	switch p[0] {
	case "super_regions":
		super_regions = parse_map(p[1:])
	case "regions":
		regions = parse_map(p[1:])
	case "neighbours":
		neighbours = parse_neighbours(p[1:])
	case "wastelands":
		wastelands = parse_ints(p[1:])
	case "opponent_starting_regions ":
		opponent_regions = parse_ints(p[1:])
	}
}

func pick_starting_region(p []string) string {
	return "whatever"
}

func update_map(p []string) {
}

func place_armies() string {
}

func attack_transfer() string {
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		f := strings.Fields(scanner.Text())
		reply := ""
		switch f[0] {
		case "settings":
			settings(f[1:])
		case "setup_map":
			setup_map(f[1:])
		case "update_map":
			update_map(f[1:])
		case "opponent_moves":
			opponent_moves(f[1:])
		case "pick_starting_region":
			reply = pick_starting_region(f[1:])
		case "go":
			switch f[1] {
			case "place_armies":
				reply = place_armies()
			case "attack/transfer":
				reply = attack_transfer()
			}
		}
		if reply == "" {
			continue
		}
		n, err := io.WriteString(os.Stdout, reply+"\n")
		if err != nil {
			log.Fatalf("%s writing reply %s", err, reply)
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("%s from Scan", err)
	}
}
