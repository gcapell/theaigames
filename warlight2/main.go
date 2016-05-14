package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
)

func opponent_moves(p []string) {

}

func pick_starting_region(p []string) string {
	p = p[1:] // skip over time limit
	return p[rand.Intn(len(p))]
}

func update_map(p []string) {
	for len(p) > 0 {
		id, player, count := p[0], p[1], p[2]
		p = p[3:]
		r := regions[parseInt(id)]
		r.owner = player
		r.armies = parseInt(count)
	}
}

func place_armies() string {
	// choose a random region we own, dump all armies on it
	var mine []*region
	for _, r := range regions {
		if r.owner == me {
			mine = append(mine, r)
		}
	}
	r := mine[rand.Intn(len(mine))]
	r.armies += starting_armies
	return fmt.Sprintf("%s place_armies %d %d", me, r.id, starting_armies)
}

func attack_transfer() string {
	// for every region we own, if we can successfully attack something, do
	var attacks []string
	var attacked []int
	for _, r := range regions {
		if r.owner != me {
			continue
		}
		if dst, n, ok := r.attack(attacked); ok {
			attacks = append(attacks, fmt.Sprintf("%s attack/transfer %d %d %d", me, r.id, dst, n))
			attacked = append(attacked, dst)
		}
	}
	if len(attacks) == 0 {
		return "No moves"
	}
	return strings.Join(attacks, ", ")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		f := strings.Fields(scanner.Text())
		reply := ""
		if len(f) < 2 {
			continue
		}
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
		os.Stdout.Sync()
		if err != nil {
			log.Fatalf("%s writing reply %s", err, reply)
		}

	}
	log.Printf("Scan loop exit. EOF?")
	if err := scanner.Err(); err != nil {
		log.Fatalf("%s from Scan", err)
	}
	showState()
}
