/*
This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org>
*/

// File contains various helpers and shortcuts for all code base

package libs

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/user"

	"golang.org/x/net/proxy"
	"gopkg.in/yaml.v2"
)

// GetHomeDirectory return path to home directory for current os user.
func GetHomeDirectory() (homeDir string) {
	if user, reason := user.Current(); reason != nil {
		log.Fatal("[-]", reason)
	} else {
		homeDir = user.HomeDir
	}
	return
}

// HTTPProxyClient set proxy for http client
func HTTPProxyClient(proxyAddr string) (client http.Client, ip string) {
	ip = ""

	if sock5, reason := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct); reason != nil {
		fmt.Println("[-] Error create proxy. Check settings")
		os.Exit(1)
	} else {
		//Check external IP
		client = http.Client{Transport: &http.Transport{Dial: sock5.Dial}}
		res, reason := client.Get("http://myexternalip.com/raw")
		if reason != nil {
			fmt.Println("[-] Could not get request to remote service:", reason)
			os.Exit(1)
		}

		if data, reason := ioutil.ReadAll(res.Body); reason != nil {
			fmt.Println("[-] Could not get external ip:", reason)
			os.Exit(1)
		} else {
			ip = (string(data))
		}

	}
	return
}

// Config represent configuration yml file
type Config struct {
	DataDir        string `yaml:"data-dir"`
	ProxyAddr      string `yaml:"proxy-addr"`
	RequestTimeout int    `yaml:"request-timeout"`
	SimBackend     []struct {
		Backend struct {
			APIKey  string `yaml:"api-key"`
			Name    string `yaml:"name"`
			Timeout int    `yaml:"timeout"`
		} `yaml:"backend"`
	} `yaml:"sim-backend"`
	TelegramAPI string   `yaml:"telegram-api"`
	TelegramID  int32    `yaml:"telegram-id"`
	Channels    []string `yaml:"channels"`
}

// LoadConfig represent configuration parse from yml file and convert to struct
func LoadConfig(path string) (c Config) {
	var _config Config

	if _, reason := os.Stat(path); os.IsNotExist(reason) {
		fmt.Println("[-] Could not found configuration file by path:", path)
		os.Exit(1)
	}

	if data, reason := ioutil.ReadFile(path); reason != nil {
		fmt.Println("[-] Could not read config file:", reason)
		os.Exit(1)
	} else if reason := yaml.Unmarshal(data, &_config); reason != nil {
		fmt.Println("[-] Could not parse config file:", reason)
		os.Exit(1)
	} else {
		c = _config
	}

	return
}

var (
	firstName = [...]string{"Richard", "Andrew", "Charles", "Milton", "Tupeni", "Meher", "Aaron", "Abagnale",
		"Abbey", "Abel", "Abelson", "Abourezk", "Abrams", "Ace", "Adams", "Gaston", "Roselyne", "Francis",
		"David", "Robert", "Alain", "Michael", "Joan", "Walter", "Philip James", "Bruce", "Frank", "Jack",
		"Russell", "Julio", "Bill", "Ann", "Douglas", "Leonard Jr.", "Margaret", "Stephen", "Noel", "Abraham",
		"William", "Jim", "Dinah", "Frank", "Stephen", "Seymour", "Francis", "Barry", "Robert X.", "Quentin",
		"Davy", "Oliver", "Aleister", "David", "Robert", "Allan", "Ward", "Marie", "Richard J.", "Carson", "Daniel",
		"Johnny", "Rodney", "Samuel", "Anthony", "Clarence", "Charles", "Dale", "Bette", "Miles", "Richard", "Carl R.",
		"Moshe", "Howard", "James", "John", "Edward de", "Eugene V.", "Jack", "Edgar", "Ellen", "Eugène", "Zack de la",
		"Tom", "Brad", "Anthony de", "Jacques de", "Deng", "John", "Chauncey", "Theo de", "Jacques", "René", "Jackie", "Paul",
		"David", "Carl", "Ron", "Caroline", "Alberto", "Edward", "Khalil", "William", "André", "George", "Eric",
		"John", "Haim", "Jean", "Rudy", "William Ewart", "Jean-Luc", "Mike", "Hermann", "Rudyard", "Henry", "Zalmay", "Stephanie",
	}
	lastName = [...]string{"Hank", "Frank", "Edward", "Reuben", "Hal", "James", "Creighton", "Jane", "Abigail", "Baba",
		"Babbage", "Babbitt", "Bacevich", "Bach", "Bachelard", "Bachelot", "Bacon", "Baddiel", "Baden-Powell", "Badiou", "Badnarik",
		"Baez", "Bagehot", "Bailey", "Baillie", "Bainimarama", "Baker", "Cortázar", "Cosby", "Coulter", "Coupland",
		"Courtney", "Covington", "Coward", "Cowley", "Cowper", "Crace", "Craik", "Crane", "Crane", "Cra", "Crick", "Crimmins", "Cringely",
		"Crisp", "Crockett", "Cromwell", "Cross", "Crowley", "Crumb", "Cunningham", "Curie", "Daley", "Daly",
		"Damon", "Dangerfield", "Daniel", "Daniels", "Darrow", "Darwin", "Dauten", "Davisson", "Dawkins", "Dayan",
		"Dean", "Dea", "Bono", "Debs", "Dee", "Degas", "Degeneres", "Delacroix", "Rocha", "Delay", "Delson", "Mello",
		"Molay", "Xiaoping", "Denham", "Depew", "Raadt", "Derrida", "Descartes", "DeShannon", "Desmond", "Deutsch",
		"Devine", "DeWolfe", "Dhavernas", "Giacometti", "Gibbon", "Gibran", "Gibson",
		"Gide", "Gilder", "Gill", "Gilmore", "Ginott", "Giraudoux", "Giuliani", "Gladstone", "Godard",
		"Godwin", "Goering", "Kipling", "Kissinger", "Khalilzad", "Zacharek",
	}
)

// GenerateUsername generate names for Telegram account
func GenerateUsername() (first string, last string) {
	first, last = firstName[rand.Intn(len(firstName))], lastName[rand.Intn(len(lastName))]
	return
}
