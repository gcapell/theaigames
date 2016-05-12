package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// Settings
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
		starting_regions = parseInts(p[1:])
	case "starting_pick_amount":
		starting_pick_amount = parseInt(p[1])
	default:
		log.Fatal("settings unrecognised: ", p)
	}
	return ""
}

func parseInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return n
}

func parseInts(ss []string) []int {
	reply := make([]int, len(ss))
	for n, s := range ss {
		reply[n] = parseInt(s)
	}
	return reply
}

func commaInts(s string) []int {
	var reply []int
	for _, n := range strings.Split(s, ",") {
		reply = append(reply, parseInt(n))
	}
	return reply
}

type (
	region    int
	bonus     int
	continent int
)

var (
	continentBonus    map[continent]bonus
	regionToContinent map[region]continent
	adjacency         map[region][]region
	wasteland         map[region]bool
	opponentStarting  map[region]bool
)

func addNeighbour(a, b region) {
	adjacency[a] = append(adjacency[a], b)
	adjacency[b] = append(adjacency[b], a)
}

func setup_map(p []string) {
	switch p[0] {
	case "super_regions":
		continentBonus = make(map[region]bonus)
		ns := parseInts(p[1:])
		for len(ns) > 0 {
			r, b := ns[0], ns[1]
			ns = ns[2:]
			continentBonus[continent(r)] = bonus(b)
		}
	case "regions":
		regionToContinent = make(map[region]continent)
		ns := parseInts(p[1:])
		for len(ns) > 0 {
			r, s := ns[0], ns[1]
			ns = ns[2:]
			regionToContinent[region(r)] = continent(s)
		}
	case "neighbours":
		p = p[1:]
		for len(p) > 0 {
			r, ns := p[0], p[1]
			p = p[2:]
			for _, n := range commaInts(ns) {
				addNeighbour(region(parseInt(r)), region(n))
			}
		}
	case "wastelands":
		wasteland = make(map[region]bool)
		for _, r := range parseInts(p[1:]) {
			wasteland[region(r)] = true
		}
	case "opponent_starting_regions":
		opponentStarting = make(map[region]bool)
		for _, r := range parseInts(p[1:]) {
			opponentStarting[region(r)] = true
		}
	default:
		log.Fatal("unrecognised setup_map", p)
	}
}

func opponent_moves(p []string) {

}

func pick_starting_region(p []string) string {
	return "whatever"
}

func update_map(p []string) {
}

func place_armies() string {
	return ""
}

func attack_transfer() string {
	return ""
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
		default:
			log.Fatal("unrecognised ", f)
		}
		if reply == "" {
			continue
		}
		_, err := io.WriteString(os.Stdout, reply+"\n")
		if err != nil {
			log.Fatalf("%s writing reply %s", err, reply)
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("%s from Scan", err)
	}
}
