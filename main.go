package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"io"
	"fmt"
	"strings"
)

type cliCommand struct {
	name string
	description string
	callback func(conf *config) error
}

type config struct {
	next *string
	previous *string
}

type LocationAreaResponse struct {
    Count    int     `json:"count"`
    Next     *string `json:"next"`
    Previous *string `json:"previous"`
    Results  []struct {
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"results"`
}

//test this
func createCommandMap() map[string]cliCommand {
	commandMap := map[string]cliCommand {
		"exit": {
			name: "exit",
			description: "exit the pokedex",
			callback: commandExit,
		},
		"help": {
			name: "help",
			description: "displays a help message",
			callback: commandHelp,
		},
		"map": {
			name: "map",
			description: "shows the next 20 locations",
			callback: commandMap,
		},
		"mapb": {
			name: "mapb",
			description: "shows previous 20 locations",
			callback: commandMapb,
		},
	}
	return commandMap
}

func cleanInput(text string) []string {
	clean := strings.Fields(text) 
	for i := range clean {
		clean[i] = strings.ToLower(clean[i])	
	}
	return clean
}

func commandExit(conf *config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	err := fmt.Errorf("could not close program")
	return err
}

func commandHelp(conf *config) error {
	commandMap := createCommandMap()
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for i := range commandMap {
		fmt.Printf("%s: %s\n", i, commandMap[i].description)	
	}
	return nil
}

func commandMap(conf *config) error {
	if conf.next == nil {
		fmt.Printf("End of map\n")
		return nil
	}
	res, err := http.Get(*conf.next)
	if err != nil {
		fmt.Printf("Could not retrieve map data\n")
		return err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		fmt.Printf("could not connect to the server\n")
		return nil
	}
	if err != nil {
		fmt.Printf("could not retrieve map data\n")
		return err
	}
	var location LocationAreaResponse
	err = json.Unmarshal(body, &location)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, result := range location.Results {
		fmt.Println(result.Name)
	}
	conf.previous = conf.next
	conf.next = location.Next
	return nil
}

func commandMapb(conf *config) error {
	if conf.previous == nil {
		fmt.Println("Start of map")
		return nil
	}
	res, err := http.Get(*conf.previous)
	if err != nil {
		fmt.Println(err)
		return err
	}
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		err = fmt.Errorf("Could not connect to server\n")
		fmt.Println(err)
		return err
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	var location LocationAreaResponse
	err = json.Unmarshal(body, &location)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, result := range location.Results {
		fmt.Println(result.Name)
	}
	conf.next = conf.previous
	conf.previous = location.Previous
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var conf *config
	mapConf := config{}
	mapURL := "https://pokeapi.co/api/v2/location-area"
	mapConf.next = &mapURL
	conf = &mapConf
	for  {
		commandMap := createCommandMap()
		fmt.Print("Pokedex > ")
		scanner.Scan()
		line := scanner.Text()
		clean := cleanInput(line)
		if len(clean) == 0 {
			continue
		}
		val, ok := commandMap[clean[0]]
		if ok {
			val.callback(conf)
		} else {
			fmt.Printf("Unknown command\n")
		}
	}
	return
}
