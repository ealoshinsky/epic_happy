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
	"sync"
	"time"

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
			os.Exit(1)
		}
	case "NO_ACTIVATION":
		{
			fmt.Println("[-] No activation found")
		}
	}
}

// ExecuteSmskoRu implement io loop for get virtual number, get sms from telegram (register new account)
// and store session by path from Config.DataDir
func ExecuteSmskoRu(countNumbers int, c *libs.Config) {
	var APIKey string
	// Warning
	if countNumbers > 45 {
		fmt.Println(strings.ToUpper("[WARNING]\nThe number of leased phone numbers is more than 45\n" +
			"considering the features of the operating system, not all rented phone \n" +
			"numbers will be embodied in the telegram accounts."))
	}

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
		fmt.Println("[+] Your balance:", balance, "rub")
	}

	fmt.Println("[!] Checks work environment.")
	TempSessionPath := os.TempDir() + "/epdata/"
	if reason := os.MkdirAll(TempSessionPath, 0777); reason != nil {
		fmt.Println("Could not create directory:", reason)
	}

	lock := &sync.WaitGroup{}
	for range make([]struct{}, countNumbers) {
		lock.Add(1)
		go func(l *sync.WaitGroup) {
			defer l.Done()
			time.Sleep(250 * time.Microsecond)
			phoneNumber, orderID := SmskoGetNumber(APIKey, client)
			if phoneNumber != "" && orderID != "" {
				fmt.Println("[+] Success order phone number:", phoneNumber, "with orderID:", orderID)
			} else {
				lock.Done()
			}
			sessionPath := TempSessionPath + phoneNumber
			session := libs.NewSession(c.TelegramID, c.TelegramAPI, sessionPath)
			fmt.Println(sessionPath)
			if session == nil {
				os.Remove(sessionPath)
				lock.Done()
			}
			if reason := session.ConnectToServer(); reason != nil {
				fmt.Println("[-] Error:", reason)
				os.Remove(sessionPath)
				lock.Done()
			}

			authCode, reason := session.AuthSendCode(phoneNumber)
			if reason != nil {
				fmt.Println("[-] Error sent authenticate code:", reason)
				os.Remove(sessionPath)
				lock.Done()
			} else {
				fmt.Println(SmskoSetStatus(APIKey, orderID, ready, client))
			}

			if !authCode.Phone_registered {
				fmt.Println("Phone number", phoneNumber, "isn't registered")
			} else {
				fmt.Println("Phone number is registered")
				SmskoSetStatus(APIKey, orderID, ban, client)
				os.Remove(sessionPath)
				lock.Done()
			}

			deadline := time.After(18 * time.Minute)
			heartbeat := time.Tick(30 * time.Second)
			for {
				select {
				case <-deadline:
					{
						os.Remove(sessionPath)
						SmskoSetStatus(APIKey, orderID, cancel, client)
						fmt.Println("End of time generate account. Clear all")
						lock.Done()
					}
				case <-heartbeat:
					{
						fmt.Println("Processing and wait sms code")
						status := SmskoGetStatus(APIKey, orderID, client)

						if status == "STATUS_WAIT_CODE" {
							fmt.Println("Wait sms code")
						} else if status == "NO_ACTIVATION" {
							fmt.Println("No activation found.")
							lock.Done()
						} else if status == "BAD_STATUS" {
							fmt.Println("[-] Unknown error. Wrong status.")
							lock.Done()
						} else if strings.Contains(status, "STATUS_OK") {
							sms := strings.Split(status, ":")
							reason := session.RegisterNewAccount(phoneNumber, sms[1], authCode.Phone_code_hash)
							if reason != nil {
								fmt.Println("[-] Error on registration new user:", reason)
								lock.Done()
							}

							SmskoSetStatus(APIKey, orderID, cancel, client)
							session.DisconnectFromServer()
						} else {
							fmt.Println(status)
						}
					}
				}
			}
		}(lock) //
	}
	lock.Wait()
	//Copy valid session to dataDIr
	sessionFiles, reason := ioutil.ReadDir(TempSessionPath)
	if reason != nil {
		fmt.Println("[-] It is not possible to read the directory with session files")
		os.Exit(1)
	}
	//TODO: copy files from tem directory to dataDIR
	for sessioFileID := range sessionFiles {
		if reason := os.Rename(TempSessionPath+sessionFiles[sessioFileID].Name(),
			c.DataDir+"/"+sessionFiles[sessioFileID].Name()); reason != nil {
			fmt.Println("[-] Could not copy session files")
		}
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

// SmskoGetStatus return status of order by id
func SmskoGetStatus(apiKey string, orderID string, client http.Client) (status string) {
	var URL = fmt.Sprintf(smskoRuEndpointAPI+"api_key=%s&action=getStatus&id=%s", apiKey, orderID)
	if response, reason := client.Get(URL); reason != nil {
		fmt.Println("[-] Error on sent request:", reason)
		os.Exit(1)
	} else if body, reason := ioutil.ReadAll(response.Body); reason != nil {
		fmt.Println("[-] Error parse data after request:", reason)
		os.Exit(1)
	} else {
		status = string(body)
	}
	return
}

const (
	cancel = -1
	ready  = 1
	end    = 6
	ban    = 8
)

// SmskoSetStatus change rent status by order id
func SmskoSetStatus(apiKey string, orderID string, status int, client http.Client) (resp string) {
	var URL = fmt.Sprintf(smskoRuEndpointAPI+"api_key=%s&action=setStatus&status=%d&id=%s",
		apiKey, status, orderID)
	if response, reason := client.Get(URL); reason != nil {
		fmt.Println("[-] Error on sent request:", reason)
		os.Exit(1)
	} else if body, reason := ioutil.ReadAll(response.Body); reason != nil {
		fmt.Println("[-] Error parse data after request:", reason)
		os.Exit(1)
	} else {
		resp = string(body)
	}
	return
}

// SmskoGetNumber phone number rent
func SmskoGetNumber(apiKey string, client http.Client) (phoneNumber string, orderID string) {
	var URL = fmt.Sprintf(smskoRuEndpointAPI+"api_key=%s&action=getNumber&service=tg", apiKey)
	if response, reason := client.Get(URL); reason != nil {
		fmt.Println("[-] Error on sent request:", reason)
		os.Exit(1)
	} else if body, reason := ioutil.ReadAll(response.Body); reason != nil {
		fmt.Println("[-] Error parse data after request:", reason)
		os.Exit(1)
	} else {
		errorParseBody(string(body))
		data := strings.Split(string(body), ":")
		fmt.Println(data)
		orderID = data[1]
		phoneNumber = data[2]
	}
	return
}
