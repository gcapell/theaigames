package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Settings
var (
	timebank, time_per_move, max_rounds, starting_armies, starting_pick_amount int

	me, them string

	starting_regions []int
)

func settings(p []string) {
	switch p[0] {
	case "timebank":
		timebank = parseInt(p[1])
	case "time_per_move":
		time_per_move = parseInt(p[1])
	case "max_rounds":
		max_rounds = parseInt(p[1])
	case "your_bot":
		me = p[1]
	case "opponent_bot":
		them = p[1]
	case "starting_armies":
		starting_armies = parseInt(p[1])
	case "starting_regions":
		starting_regions = parseInts(p[1:])
	case "starting_pick_amount":
		starting_pick_amount = parseInt(p[1])
	default:
		log.Fatal("settings unrecognised: ", p)
	}
}

type continent struct {
	id      int
	bonus   int
	regions []*region
}

type region struct {
	owner            string
	armies           int
	id               int
	c                *continent
	near             []*region
	wasteland        bool
	opponentStarting bool
}

func (r region) String() string {
	return fmt.Sprintf("r:%d,o:%s,c:%d,a:%d",
		r.id, r.owner, r.c.id, r.armies)
}

func contains(as []int, b int) bool {
	for _, a := range as {
		if a == b {
			return true
		}
	}
	return false
}

func (r region) attack(attacked []int) (int, int, bool) {
	var weak *region
	for _, o := range r.near {
		if o.owner == me || contains(attacked, o.id) {
			continue
		}
		if weak == nil || weak.armies > o.armies {
			weak = o
		}
	}
	if weak == nil {
		return 0, 0, false
	}
	if float64(r.armies-1)*0.6 > float64(weak.armies)+0.5 {
		return weak.id, r.armies - 1, true
	}
	return 0, 0, false
}

var (
	continents map[int]*continent
	regions    map[int]*region
)

func showState() {
	log.Println("State:")
	log.Println("continents")
	for _, c := range continents {
		log.Println(*c)
	}
	log.Println("regions")
	for _, r := range regions {
		log.Printf("%s", *r)
	}
}

func addNeighbour(a, b *region) {
	a.near = append(a.near, b)
	b.near = append(b.near, a)
}

func setup_map(p []string) {
	switch p[0] {
	case "super_regions":
		continents = make(map[int]*continent)
		ns := parseInts(p[1:])
		for len(ns) > 0 {
			id, b := ns[0], ns[1]
			ns = ns[2:]
			continents[id] = &continent{id: id, bonus: b}
		}
	case "regions":
		regions = make(map[int]*region)
		ns := parseInts(p[1:])
		for len(ns) > 0 {
			id, c := ns[0], ns[1]
			ns = ns[2:]
			r := &region{id: id, c: continents[c]}
			regions[id] = r
			continents[c].regions = append(continents[c].regions, r)
		}
	case "neighbors":
		p = p[1:]
		for len(p) > 0 {
			r, ns := p[0], p[1]
			p = p[2:]
			for _, n := range commaInts(ns) {
				addNeighbour(regions[parseInt(r)], regions[n])
			}
		}
	case "wastelands":
		for _, r := range parseInts(p[1:]) {
			regions[r].wasteland = true
		}
	case "opponent_starting_regions":
		for _, r := range parseInts(p[1:]) {
			regions[r].opponentStarting = true
		}
	default:
		log.Fatal("unrecognised setup_map", p)
	}
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
