package tg

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

// APIEndpoint is Telegram's current Bot API base url endpoint
const APIEndpoint = "https://api.telegram.org/"

// WebhookHandler is a function that handles updates
type WebhookHandler func(APIUpdate)

// Telegram is the API client for the Telegram Bot API
type Telegram struct {
	Token string
}

// MakeAPIClient creates a Telegram instance from a Bot API token
func MakeAPIClient(token string) *Telegram {
	tg := new(Telegram)
	tg.Token = token
	return tg
}

// SetWebhook sets the webhook address so that Telegram knows where to send updates
func (t Telegram) SetWebhook(webhook string) {
	resp, err := http.PostForm(t.apiURL("setWebhook"), url.Values{"url": {webhook}})
	if !checkerr("SetWebhook/http.PostForm", err) {
		defer resp.Body.Close()
		var result APIResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			log.Println("[SetWebhook] Could not read reply: " + err.Error())
			return
		}
		if result.Ok {
			log.Println("Webhook successfully set!")
		} else {
			log.Printf("[SetWebhook] Error setting webhook (errcode %d): %s\n", *(result.ErrCode), *(result.Description))
			panic(errors.New("Cannot set webhook"))
		}
	}
}

// HandleWebhook is a webhook HTTP handler for standalone bots
func (t Telegram) HandleWebhook(bind string, webhook string, handler WebhookHandler) error {
	whmux := http.NewServeMux()
	whmux.HandleFunc(webhook, func(rw http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		// Re-encode request to ensure conformity
		var update APIUpdate
		err := json.NewDecoder(req.Body).Decode(&update)
		if err != nil {
			log.Println("[webhook] Received incorrect request: " + err.Error())
			return
		}

		handler(update)
	})
	return http.ListenAndServe(bind, whmux)
}

// SendTextMessage sends an HTML-styled text message to a specified chat
func (t Telegram) SendTextMessage(data ClientTextMessageData) {
	postdata := url.Values{
		"chat_id":    {strconv.FormatInt(data.ChatID, 10)},
		"text":       {data.Text},
		"parse_mode": {"HTML"},
	}
	if data.ReplyID != nil {
		postdata["reply_to_message_id"] = []string{strconv.FormatInt(*(data.ReplyID), 10)}
	}

	_, err := http.PostForm(t.apiURL("sendMessage"), postdata)
	checkerr("SendTextMessage/http.PostForm", err)
}

// SendPhoto sends a picture to a chat as a photo
func (t Telegram) SendPhoto(data ClientPhotoData) {
	// Decode photo from b64
	photolen := base64.StdEncoding.DecodedLen(len(data.Bytes))
	photobytes := make([]byte, photolen)
	decoded, err := base64.StdEncoding.Decode(photobytes, []byte(data.Bytes))
	if checkerr("SendPhoto/base64.Decode", err) {
		return
	}

	// Write file into multipart buffer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("photo", data.Filename)
	if checkerr("SendPhoto/multipart.CreateFormFile", err) {
		return
	}
	part.Write(photobytes[0:decoded])

	// Write other fields
	writer.WriteField("chat_id", strconv.FormatInt(data.ChatID, 10))

	if data.ReplyID != nil {
		writer.WriteField("reply_to_message_id", strconv.FormatInt(*data.ReplyID, 10))
	}

	if data.Caption != "" {
		writer.WriteField("caption", data.Caption)
	}

	err = writer.Close()
	if checkerr("SendPhoto/writer.Close", err) {
		return
	}

	// Create HTTP client and execute request
	client := &http.Client{}
	req, err := http.NewRequest("POST", t.apiURL("sendPhoto"), body)
	if checkerr("SendPhoto/http.NewRequest", err) {
		return
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())

	_, err = client.Do(req)
	checkerr("SendPhoto/http.Do", err)
}

// ForwardMessage forwards an existing message to a chat
func (t Telegram) ForwardMessage(data ClientForwardMessageData) {
	postdata := url.Values{
		"chat_id":      {strconv.FormatInt(data.ChatID, 10)},
		"from_chat_id": {strconv.FormatInt(data.FromChatID, 10)},
		"message_id":   {strconv.FormatInt(data.MessageID, 10)},
	}

	_, err := http.PostForm(t.apiURL("forwardMessage"), postdata)
	checkerr("ForwardMessage/http.PostForm", err)
}

// SendChatAction sends a 5 second long action (X is writing, sending a photo ecc.)
func (t Telegram) SendChatAction(data ClientChatActionData) {
	postdata := url.Values{
		"chat_id": {strconv.FormatInt(data.ChatID, 10)},
		"action":  {string(data.Action)},
	}

	_, err := http.PostForm(t.apiURL("sendChatAction"), postdata)
	checkerr("SendChatAction/http.PostForm", err)
}

// AnswerInlineQuery replies to an inline query
func (t Telegram) AnswerInlineQuery(data InlineQueryResponse) {
	jsonresults, err := json.Marshal(data.Results)
	if checkerr("AnswerInlineQuery/json.Marshal", err) {
		return
	}
	postdata := url.Values{
		"inline_query_id": {data.QueryID},
		"results":         {string(jsonresults)},
	}
	if data.CacheTime != nil {
		postdata["cache_time"] = []string{strconv.Itoa(*data.CacheTime)}
	}
	if data.IsPersonal {
		postdata["is_personal"] = []string{"true"}
	}
	if data.NextOffset != "" {
		postdata["next_offset"] = []string{data.NextOffset}
	}
	if data.PMText != "" {
		postdata["switch_pm_text"] = []string{data.PMText}
	}
	if data.PMParam != "" {
		postdata["switch_pm_parameter"] = []string{data.PMParam}
	}

	_, err = http.PostForm(t.apiURL("answerInlineQuery"), postdata)
	checkerr("AnswerInlineQuery/http.PostForm", err)
}

// GetFile sends a "getFile" API call to Telegram's servers and fetches the file
// specified afterward. The file will be then send back to the client that requested it
// with the specified callback id.
func (t Telegram) GetFile(data FileRequestData, client net.Conn, callback int) {
	fail := func(msg string) {
		errmsg, _ := json.Marshal(BrokerUpdate{
			Type:     BError,
			Error:    &msg,
			Callback: &callback,
		})
		fmt.Fprintln(client, string(errmsg))
	}

	postdata := url.Values{
		"file_id": {data.FileID},
	}
	resp, err := http.PostForm(t.apiURL("getFile"), postdata)
	if checkerr("GetFile/post", err) {
		fail("Server didn't like my request")
		return
	}
	defer resp.Body.Close()

	var filespecs = struct {
		Ok     bool     `json:"ok"`
		Result *APIFile `json:"result,omitempty"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&filespecs)
	if checkerr("GetFile/json.Decode", err) {
		fail("Server sent garbage (or error)")
		return
	}
	if filespecs.Result == nil {
		fail("Server didn't send a file info, does the file exist?")
		return
	}
	result := *filespecs.Result

	path := APIEndpoint + "file/bot" + t.Token + "/" + *result.Path
	fileresp, err := http.Get(path)
	if checkerr("GetFile/get", err) {
		fail("Could not retrieve file from Telegram's servers")
		return
	}
	defer fileresp.Body.Close()

	rawdata, err := ioutil.ReadAll(fileresp.Body)
	if checkerr("GetFile/ioutil.ReadAll", err) {
		fail("Could not read file data")
		return
	}

	rawlen := len(rawdata)
	if rawlen != *result.Size {
		// ???
		log.Printf("[GetFile] WARN ?? Downloaded file does not match provided filesize: %d != %d\n", rawlen, *result.Size)
	}
	b64data := base64.StdEncoding.EncodeToString(rawdata)

	clientmsg, err := json.Marshal(BrokerUpdate{
		Type:     BFile,
		Bytes:    &b64data,
		Callback: &callback,
	})
	if checkerr("GetFile/json.Marshal", err) {
		fail("Could not serialize reply JSON")
		return
	}

	fmt.Fprintln(client, string(clientmsg))
}

func (t Telegram) apiURL(method string) string {
	return APIEndpoint + "bot" + t.Token + "/" + method
}

func checkerr(method string, err error) bool {
	if err != nil {
		log.Printf("[%s] Error: %s\n", method, err.Error())
		return true
	}
	return false
}
