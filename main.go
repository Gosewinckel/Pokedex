package main

import (
	"Pokedex/internal/pokecache"
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"io"
	"fmt"
	"strings"
	"time"
)

type cliCommand struct {
	name string
	description string
	callback func(conf *config, cache *pokecache.Cache, area *string) error
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
		"explore": {
			name: "explore",
			description: "explore an area",
			callback: commandExplore,
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

func commandExit(conf *config, cache *pokecache.Cache, area *string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	err := fmt.Errorf("could not close program")
	return err
}

func commandHelp(conf *config, cache *pokecache.Cache, area *string) error {
	commandMap := createCommandMap()
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for i := range commandMap {
		fmt.Printf("%s: %s\n", i, commandMap[i].description)	
	}
	return nil
}

func commandMap(conf *config, cache *pokecache.Cache, area *string) error {
	if conf.next == nil {
		fmt.Printf("End of map\n")
		return nil
	}
	var location LocationAreaResponse

	//in cache
	if data, ok := cache.Get(*conf.next); ok {
		cache.Add(*conf.next, data)
		err := json.Unmarshal(data, &location)
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else {
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
		err = json.Unmarshal(body, &location)
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, result := range location.Results {
			fmt.Println(result.Name)
		}
	}
	conf.previous = conf.next
	conf.next = location.Next
	return nil
}

func commandMapb(conf *config, cache *pokecache.Cache, area *string) error {
	if conf.previous == nil {
		fmt.Println("Start of map")
		return nil
	}
	var location LocationAreaResponse

	if data, ok := cache.Get(*conf.previous); ok {
		cache.Add(*conf.previous, data)
		err := json.Unmarshal(data, &location)
		if err != nil {
			fmt.Println(err)
			return err
		}
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

func commandExplore(conf *config, cache *pokecache.Cache, area *string) error {
	url := "https://pokeapi.co/api/v2/location-area" + *area
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var conf *config
	mapConf := config{}
	mapURL := "https://pokeapi.co/api/v2/location-area"
	mapConf.next = &mapURL
	conf = &mapConf
	commandMap := createCommandMap()
	cache := pokecache.NewCache(time.Second * 5)
	for  {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		line := scanner.Text()
		clean := cleanInput(line)
		if len(clean) == 0 {
			continue
		}
		val, ok := commandMap[clean[0]]
		if ok {
			val.callback(conf, cache)
		} else {
			fmt.Printf("Unknown command\n")
		}
	}
	return
}
