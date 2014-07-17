package go2bosh

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/may215/go2bosh/service"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var (
	jabber_rid int
	jabber_sid string
	config     = new(BoshConfiguration)
	user_id    string
)

/* Configuration struct to contain all the bosh configuration service and the xml data */
type BoshConfiguration struct {
	BoshServer            string
	BOSHMethod            string
	RequestTimeOut        int
	RandomAuthKeyInterval int
	ServicePort           int
	ServiceDomain         string
	Wait                  string
	Hold                  string
	XmlLang               string
	Content               string
	Ver                   string
	Mechanism             string
	TlsXmlns              string
	SaslXmlns             string
	BindXmlns             string
	SessionXmlns          string
}

/* Holds the auth response data */
type AuthResponse struct {
	Jid string /* the current service authenticated user id, e.g. (-12345678@chat.facebook.com) */
	Sid string /* bosh session id */
	Rid string /* auto generated request id */
}

/* Return the configuration data from the bosh config file */
func getConfigBoshData() *handlerError {
	file, err := os.Open("bosh.conf")
	if err != nil {
		return &handlerError{err, "Unable to open file for read", 100001}
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var conf = strings.Join(lines, " ")
	b := []byte(conf)
	m_err := json.Unmarshal(b, &config)
	if m_err != nil {
		return &handlerError{m_err, "Unable to marshal configuration file data", 100002}
	}

	return &handlerError{err, "success reading config file", 0}
}

/* Increments the rid by one to send the next payload of xml. */
func jabber_get_next_rid() int {
	jabber_rid = rand.Intn(config.RandomAuthKeyInterval)
	return jabber_rid
}

/* Connect to punjab server flow handler  */
func Bosh_Connect(su ServiceUser) (*AuthResponse, *handlerError) {
	/* Get the service configuration */
	conf_err := getConfigBoshData()
	if conf_err.Error != nil {
		return nil, &handlerError{conf_err.Error, conf_err.Message, 100003}
	}

	/* Set the start session rid */
	jabber_get_next_rid()

	/* Increment the rid */
	jabber_rid += 1

	/* Stage 1 */
	resp_one, err_one := getSidRequest()
	/* Validate stage response and data */
	if err_one != nil {
		return nil, &handlerError{err_one.Error, err_one.Message, err_one.Code}
	}
	jabber_sid = getXmlData(resp_one, "sid")
	if jabber_sid == "" {
		return nil, &handlerError{errors.New("Unable to get sid value from the response"), "", 100004}
	}
	/* Increment the rid */
	jabber_rid += 1

	/* Stage 2 */
	resp_two, err_two := getStartChallengeRequest()
	/* Validate stage response and data */
	if err_two != nil {
		return nil, &handlerError{err_two.Error, err_two.Message, err_two.Code}
	}
	challenge, _ := getElementText(resp_two, "challenge")
	if challenge == "" {
		return nil, &handlerError{errors.New("Unable to make challenge request"), "", 100005}
	}
	/* Increment the rid */
	jabber_rid += 1
	/* Stage 3 */
	resp_three, err_three := getAuthenticateRequest(su)
	/* Validate stage response and data */
	if err_three != nil {
		return nil, &handlerError{err_three.Error, err_three.Message, err_three.Code}
	}
	fmt.Println(string(resp_three))
	_, a_exist := getElementText(resp_three, "success")

	if !a_exist {
		return nil, &handlerError{errors.New("Unable to start auth request"), "", 100006}
	}
	/* Increment the rid */
	jabber_rid += 1

	/* Stage 4 */
	resp_four, err_four := getRestartRequest()

	/* Validate stage response and data */
	if err_four != nil {
		return nil, &handlerError{err_four.Error, err_four.Message, err_four.Code}
	}

	_, r_exists := getElementText(resp_four, "bind")
	if !r_exists {
		return nil, &handlerError{errors.New("Unable to set restart request"), "", 100007}
	}
	/* Increment the rid */
	jabber_rid += 1

	/* Stage 5 */
	resp_five, err_five := getBindRequest()

	/* Validate stage response and data */
	if err_five != nil {
		return nil, &handlerError{err_five.Error, err_five.Message, err_five.Code}
	}
	bind := getXmlData(resp_five, "from")
	if bind == "" {
		return nil, &handlerError{errors.New("Unable to make bind request"), "", 100008}
	}
	/* Increment the rid */
	jabber_rid += 1

	/* Stage 6 */
	resp_six, err_six := setSessionRequest()
	/* Validate stage response and data */
	if err_six != nil {
		return nil, &handlerError{err_six.Error, err_six.Message, err_six.Code}
	}

	// Get the jid from the final response
	jid := getXmlData(resp_six, "from")
	if jid == "" {
		return nil, &handlerError{errors.New("Unable to get jid"), "", 100009}
	}
	auth_resp := AuthResponse{Rid: strconv.Itoa(jabber_rid), Sid: jabber_sid, Jid: jid}

	return &auth_resp, nil
}

/* Step 1: Create a request to get a SID and establish communications. Send your requests using HttpWebRequest to your jabber server url. Once you should get a response that contains a SID and the AuthId. Save the SID. Check XEP 206 for more details on responses.*/
func getSidRequest() ([]byte, *handlerError) {
	data := "<body rid='" + strconv.Itoa(jabber_rid) + "' xmlns='http://jabber.org/protocol/httpbind' to='" + config.ServiceDomain + "' xml:lang='" + config.XmlLang + "' wait='" + config.Wait + "' hold='" + config.Hold + "' content='" + config.Content + "' ver='1.6' xmpp:version='" + config.Ver + "' xmlns:xmpp='urn:xmpp:xbosh'/>"
	res_data, err := makeRequest(data, "")
	if err != nil {
		return nil, &handlerError{err.Error, "Unable to create request to get sid", 100010}
	}
	return res_data, nil
}

/* Step 2: Start challenge request for the specific service provider we want to connect to, e.g. X-FACEBOOK-PLATFORM */
func getStartChallengeRequest() ([]byte, *handlerError) {
	data := "<body content='" + config.Content + "' xml:lang='" + config.XmlLang + "' rid='" + strconv.Itoa(jabber_rid) + "' xmlns='http://jabber.org/protocol/httpbind' sid='" + jabber_sid + "'><auth xmlns='" + config.SaslXmlns + "' mechanism='" + config.Mechanism + "'/></body>"
	res_data, err := makeRequest(data, "")
	if err != nil {
		return nil, &handlerError{err.Error, "Unable to create challenge request to start challenge", 100011}
	}
	return res_data, nil
}

/* Step 3: Authenticate. Essentially login your user in this step. Here we need to decide what your mechanism for encrypting is going to be, but for simplicity, we'll use PLAIN. We need to convert authentication to a base64. Here we start passing the SID back to the server. We MUST increment the RID each time. On this call, we'll be looking for a "success" response. */
func getAuthenticateRequest(su ServiceUser) ([]byte, *handlerError) {
	auth_data, err_auth := setAuthRequestData(su)
	if err_auth.Error != nil {
		return nil, &handlerError{err_auth.Error, err_auth.Message, err_auth.Code}
	}
	auth_data = "<response xmlns='urn:ietf:params:xml:ns:xmpp-sasl'>" + auth_data + "</response>"
	data := "<body content='" + config.Content + "' xml:lang='" + config.XmlLang + "' rid='" + strconv.Itoa(jabber_rid) + "' xmlns='http://jabber.org/protocol/httpbind' sid='" + jabber_sid + "'>" + auth_data + "</body>"
	res_data, err := makeRequest(data, "")
	if err != nil {
		return nil, &handlerError{err.Error, "Unable to create authentication request", 100012}
	}
	return res_data, nil
}

/* Step 4: Restart. Send a restart command back to the jabber server. This should return a stream features element if you are successful. */
func getRestartRequest() ([]byte, *handlerError) {
	data := "<body rid='" + strconv.Itoa(jabber_rid) + "' xmlns='http://jabber.org/protocol/httpbind' sid='" + jabber_sid + "' to='chat.facebook.com' xml:lang='en' xmpp:restart='true' xmlns:xmpp='urn:xmpp:xbosh'/>"
	res_data, err := makeRequest(data, "")
	if err != nil {
		return nil, &handlerError{err.Error, "Unable to create restart request", 100013}
	}
	return res_data, nil
}

/* Step 5: Bind. This will bind our resource to the jabber id. The resource is what we use to represent where our user. This is where we'll get the JID returned. */
func getBindRequest() ([]byte, *handlerError) {
	data := "<body rid='" + strconv.Itoa(jabber_rid) + "' xmlns='http://jabber.org/protocol/httpbind' sid='" + jabber_sid + "'><iq type='set' id='_bind_auth_2' xmlns='jabber:client'><bind xmlns='" + config.BindXmlns + "'/></iq></body>"
	res_data, err := makeRequest(data, "")
	if err != nil {
		return nil, &handlerError{err.Error, "Unable to create bind request", 100014}
	}
	return res_data, nil
}

/* Step 6: Set the Session. This will establish the session on the jabber server with our shiny new JID. This returns a session id. */
func setSessionRequest() ([]byte, *handlerError) {
	data := "<body rid='" + strconv.Itoa(jabber_rid) + "' xmlns='http://jabber.org/protocol/httpbind' sid='" + jabber_sid + "'><iq type='set' id='_session_auth_2' xmlns='jabber:client'><session xmlns='" + config.SessionXmlns + "'/></iq></body>"
	res_data, err := makeRequest(data, "")
	if err != nil {
		return nil, &handlerError{err.Error, "Unable to create authentication session request to bosh", 100015}
	}
	return res_data, nil
}

/* Step two: create the authenticated session against the endpoint service  */
func setAuthRequestData(su ServiceUser) (string, *handlerError) {
	//Get the user access token
	token, err := su.GetUserData()
	orig_err, _ := err.(handlerError)

	if orig_err.Error != nil {
		return "", &handlerError{orig_err.Error, orig_err.Message, orig_err.Code}
	}
	auth_data, auth_err := su.SetAuthRequestData(token)
	orig_err, _ = auth_err.(handlerError)

	return auth_data, &orig_err
}
