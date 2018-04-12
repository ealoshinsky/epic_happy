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

// The file describes the methods of working with the API telegram

package libs

import (
	"errors"
	"fmt"

	"github.com/shelomentsevd/mtproto"
)

// SessionTelegram integrate mtproto.MTProto
type SessionTelegram struct{ mtproto.MTProto }

// NewSession initialize new session telegram SessionTelegram file by path
func NewSession(TelegramID int32, TelegramAPI string, sessionPath string) *SessionTelegram {

	if proto, reason := mtproto.NewMTProto(TelegramID, TelegramAPI,
		mtproto.WithAuthFile(sessionPath, false)); reason != nil {
		fmt.Println("Could not create session by path:", sessionPath, "[", reason, "]")
		return nil
	} else {
		return &SessionTelegram{(*proto)}
	}

	return nil
}

// ConnectToServer sets connection to telegram service
func (account *SessionTelegram) ConnectToServer() error {
	if reason := account.Connect(); reason != nil {
		return errors.New(reason.Error())
	}
	return nil
}

//DisconnectFromServer drop connection from telegram service
func (account *SessionTelegram) DisconnectFromServer() error {
	if reason := account.Disconnect(); reason != nil {
		return errors.New(reason.Error())
	}
	return nil
}

//RegisterNewAccount create SessionTelegram in telegram
func (account *SessionTelegram) RegisterNewAccount(PhoneNumber string, smsCode string, Phone_code_hash string) error {
	// Let's attempt create SessionTelegram in telegram
	// Before try it, we need send sms code to PhoneNumber and check this code after
	// and if new SessionTelegram success created user make auto sign in (see: SentCode, GetSmsFromSimSms  methods)
	firstName, lastName := GenerateUsername()
	if tl, reason := account.InvokeSync(mtproto.TL_auth_signUp{Phone_number: PhoneNumber,
		Phone_code_hash: Phone_code_hash,
		Phone_code:      smsCode,
		First_name:      firstName,
		Last_name:       lastName,
	}); reason != nil {
		// invoke signup
		return errors.New(reason.Error())
	} else {
		fmt.Println("[+] Success registed account with name", tl)
	}

	return nil

}
