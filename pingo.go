package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"bufio"

	"github.com/spf13/pflag"
)

//go:embed osinfo.json
var osinfoJSON []byte

type OSInfo struct {
	DeviceOS string `json:"Device / OS"`
	Version  string `json:"Version"`
	Protocol string `json:"Protocol"`
	TTL      string `json:"TTL"`
}

type Guess struct {
	OSInfo
	Difference int
}

type ByDifference []Guess

func (a ByDifference) Len() int           { return len(a) }
func (a ByDifference) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDifference) Less(i, j int) bool { return a[i].Difference < a[j].Difference }

func getTTLFromPing(target string) (int, string, error) {
	out, err := exec.Command("ping", "-c", "1", target).Output()
	if err != nil {
		return 0, "", err
	}

	// regex to match ttl from ping output
	ttlRegex := regexp.MustCompile(`ttl=(\d+)`)
	matches := ttlRegex.FindStringSubmatch(string(out))
	if len(matches) != 2 {
		return 0, "", fmt.Errorf("could not find TTL in ping output")
	}

	ttl, err := strconv.Atoi(matches[1])
	return ttl, string(out), err
}

func processTarget(target string, tries, tolerance int, verbose bool) {
	// Load OSInfo data
	var osInfo []OSInfo
	err := json.Unmarshal(osinfoJSON, &osInfo)
	if err != nil {
		fmt.Println("Error parsing osinfo.json:", err)
		os.Exit(1)
	}

	// Get average TTL value and guess OS for each try
	guessedTTLs := make([]int, tries)
	for i := 0; i < tries; i++ {
		// Get number of hops using traceroute
		out, err := exec.Command("traceroute", "-w", "5", "-n", "-m", "64", target).Output()
		if err != nil {
			fmt.Println("Error running traceroute:", err)
			os.Exit(1)
		}

		if verbose {
			fmt.Println(string(out)) // print traceroute output
		}

		hops := len(strings.Split(string(out), "\n")) - 1

		// Get TTL value from pinging the target
		ttl, pingOutput, err := getTTLFromPing(target)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return
		}

		if verbose {
			fmt.Println(pingOutput) // print ping output
		}

		// Guess OS based on TTL value
		guessedTTL := ttl + hops
		guessedTTLs[i] = guessedTTL
	}

	// Calculate average TTL
	sum := 0
	for _, ttl := range guessedTTLs {
		sum += ttl
	}
	avgTTL := sum / len(guessedTTLs)

	// Calculate standard deviation
	var sumOfSquares int
	for _, ttl := range guessedTTLs {
		sumOfSquares += (ttl - avgTTL) * (ttl - avgTTL)
	}
	stdDev := math.Sqrt(float64(sumOfSquares) / float64(len(guessedTTLs)))

	// Guess OS based on average TTL with tolerance
	var guesses []Guess
	for _, info := range osInfo {
		ttl, _ := strconv.Atoi(info.TTL)
		difference := abs(ttl - avgTTL)
		if difference <= tolerance {
			guesses = append(guesses, Guess{info, difference})
		}
	}

	// Sort by difference and print
	sort.Sort(ByDifference(guesses))
	for _, guess := range guesses {
		fmt.Println(guess.DeviceOS, guess.Version, "diff:", guess.Difference)
	}

	if verbose {
		fmt.Println("\nAverage TTL:", avgTTL, "Standard Deviation:", stdDev)
	}
}

func main() {
	tries := pflag.Int("retries", 1, "Number of tries")
	tolerance := pflag.Int("tolerance", 0, "TTL tolerance")
	verbose := pflag.Bool("verbose", false, "Print verbose output")
	pflag.Parse()

	targets := pflag.Args()
	if len(targets) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			targets = append(targets, scanner.Text())
		}
	}

	if len(targets) == 0 {
		fmt.Println("Please specify a target")
		os.Exit(1)
	}

	for _, target := range targets {
		processTarget(target, *tries, *tolerance, *verbose)
	}
}

// abs returns the absolute value of an integer
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
