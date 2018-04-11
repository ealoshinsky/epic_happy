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

// File describe wrapper around API for http://smsko.ru (detail: http://smsko.ru/api.php)

package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ealoshinsky/epic_happy/libs"
)

const (
	smskoRuEndpointAPI = "http://smsko.ru/stubs/handler_api.php?"
)

func errorParseBody(body string) {
	switch string(body) {
	case "BAD_KEY":
		{
			fmt.Println("[-] Error parse response:BAD_KEY")
			os.Exit(1)
		}
	case "BAD_SERVICE":
		{
			fmt.Println("[-] Wrong service.")
			os.Exit(1)
		}
	case "BAD_ACTION":
		{
			fmt.Println("[-] Wrong action.")
			os.Exit(1)
		}
	case "NO_BALANCE":
		{
			fmt.Println("[-] Your balance to low. Check it")
			os.Exit(1)
		}
	case "NO_NUMBERS":
		{
			fmt.Println("[-] No phone numbers for rent.")
		}
	}
}

// ExecuteSmskoRu implement io loop for get virtual number, get sms from telegram (register new account)
// and store session by path from Config.DataDir
func ExecuteSmskoRu(countNumbers int, c *libs.Config) {
	var (
		APIKey string
	)
	for backend := range c.SimBackend {
		if strings.ToLower(c.SimBackend[backend].Backend.Name) == "smsko_ru" ||
			strings.ToLower(c.SimBackend[backend].Backend.Name) == "smsko.ru" {
			APIKey = c.SimBackend[backend].Backend.APIKey
		}
	}

	fmt.Println("[*] Prepares http client for request.")
	client, extIP := libs.HTTPProxyClient(c.ProxyAddr)
	fmt.Println("[+] Used connection ip address is ", extIP)

	countPhonesNumber := SmskoGetNumberStatus(APIKey, client)
	if countPhonesNumber == 0 {
		fmt.Println("[-] There are no phone numbers available for rent. Try later.")
		os.Exit(1)
	} else if countNumbers > countPhonesNumber && countPhonesNumber > 0 {
		fmt.Println("[!] Hey we try order only ", countPhonesNumber, "phones number. Not all what you need. Sorry...")
		countNumbers = countPhonesNumber
	} else {
		fmt.Println("[+] Available phone numbers is:", countPhonesNumber)
	}

	balance := SmskoGetBalance(APIKey, client)
	if balance == 0.00 {
		fmt.Println("Your balance to low")
		os.Exit(1)
	} else {
		fmt.Println("[+] Your balance:", balance)
	}

}

// SmskoGetNumberStatus return count available phones number for order
func SmskoGetNumberStatus(apiKey string, client http.Client) (countPhoneNumber int) {
	var URL = fmt.Sprintf(smskoRuEndpointAPI+"api_key=%s&action=getNumbersStatus", apiKey)
	if response, reason := client.Get(URL); reason != nil {
		fmt.Println("[-] Error on sent request:", reason)
		os.Exit(1)
	} else if body, reason := ioutil.ReadAll(response.Body); reason != nil {
		fmt.Println("[-] Error parse data after request:", reason)
		os.Exit(1)
	} else {
		// parse data
		data := make(map[string]int)
		if reason := json.Unmarshal(body, &data); reason != nil {
			errorParseBody(string(body))
		}
		countPhoneNumber = data["tg_0"]
	}
	return
}

// SmskoGetBalance return how much money you have.
func SmskoGetBalance(apiKey string, client http.Client) (balance float64) {
	var URL = fmt.Sprintf(smskoRuEndpointAPI+"api_key=%s&action=getBalance", apiKey)
	if response, reason := client.Get(URL); reason != nil {
		fmt.Println("[-] Error on sent request:", reason)
		os.Exit(1)
	} else if body, reason := ioutil.ReadAll(response.Body); reason != nil {
		fmt.Println("[-] Error parse data after request:", reason)
		os.Exit(1)
	} else {
		data := string(body)
		errorParseBody(data)
		_balance := strings.SplitAfter(string(body), ":")[1]
		if balance, reason = strconv.ParseFloat(_balance, 64); reason != nil {
			fmt.Println("[-] Could not convert balance.")
			os.Exit(1)
		}
	}
	return
}
