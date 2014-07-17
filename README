go2bosh
=========

> BOSH connection manager session client creator. 

The motivation for this project: http://metajack.im/2008/10/03/getting-attached-to-strophe/


How to use:

First you need to change the configuration files to feet your needs:

** This example configuration is to create xmpp session with facebook (facebook chat).

** You can see in the article above how to create the javascript code in order to manage the attach session. 

bosh.conf

    {
    	"BoshServer": "your_xmpp_bosh_server",
    	"BOSHMethod": "auth.xmpp_login",
    	"RequestTimeOut": 500000,
    	"RandomAuthKeyInterval": 10000000,
    	"ServicePort": 4040,
    	"ServiceDomain": "chat.facebook.com",
    	"Wait": "300",
    	"Hold": "1",
    	"XmlLang": "en",
    	"Content": "text/xml; charset=utf-8",
    	"Ver": "1.0",
    	"Mechanism": "X-FACEBOOK-PLATFORM",
    	"TlsXmlns": "urn:ietf:params:xml:ns:xmpp-tls",
    	"SaslXmlns": "urn:ietf:params:xml:ns:xmpp-sasl",
    	"BindXmlns": "urn:ietf:params:xml:ns:xmpp-bind",
    	"SessionXmlns": "urn:ietf:params:xml:ns:xmpp-session"
    }

Working example:
    
    package main

    import (
    	"crypto/md5"
    	"encoding/base64"
    	"encoding/hex"
    	"flag"
    	"fmt"
    	bosh "github.com/may215/go2bosh"
    	"regexp"
    	"strings"
    	"time"
    )
    
    var (
    	facebookSecret string
    	facebookAppId  string
    	boshMethod     string
    )
    
    type User struct{}
    
    /* Here you need to return the user access token from facebook */
    func (this *User) GetUserData() (string, interface{}) {
    	return "CAADXHJWnl5EBAGE89GGF9wzgKAKR2XXQoVrZAZCUcQLkm8EfmGjex3dZCgcrv7NA1kZA7rJrEviJFY7gtGIaASPvQVVg4gmuKJ6Teb7jDPF9VZBZAkpzXKRe6q", nil
    }
    
    /* Here we create the base64 authentication string */
    func (this *User) SetAuthRequestData(token string) (string, interface{}) {
    	t := time.Now()
    	tt := t.Format("20060102150405")
    	api_secret := facebookSecret
    	api_key := facebookAppId
    	call_id := string(tt)
    	method := boshMethod
    	nonce := encodeString("1234567890")
    	v := "1.0"
    	st := []string{"&api_key=", api_key, "&call_id=", call_id, "&method=", method, "&nonce=", nonce, "&access_token=", token, "&v=", v}
    	str := strings.Join(st, "")
    	reg, _ := regexp.Compile("&")
    	safe := reg.ReplaceAllString(str, "")
    	sig := encodeString(safe)
    	st = []string{str, "&sig=", sig + api_secret}
    	str = strings.Join(st, "")
    	str = base64.StdEncoding.EncodeToString([]byte(str))
    	return str, nil
    }
    
    func encodeString(data string) string {
    	hasher := md5.New()
    	hasher.Write([]byte(data))
    	return hex.EncodeToString(hasher.Sum(nil))
    }
    
    func main() {
    	flag.StringVar(&facebookSecret, "facebookSecret", "facebook app secret", "your facebook app secret")
    	flag.StringVar(&facebookAppId, "facebookAppId", "facebook app key", "your facebook app id")
    	flag.StringVar(&boshMethod, "boshMethod", "auth.xmpp_login", "endpoint xmpp auth method")
    
    	var user User
    	authResponse, err := bosh.Bosh_Connect(&user)
    	/* The output of the method is {Rid: rid, Sid: sid, Jid: jid} */
    	fmt.Println(authResponse)
    	fmt.Println(err)
    }

Version
----

1.0

License
----

MIT

Author
----

Meir Shamay @meir_shamay

**Free Software, Hell Yeah!**

[@meir_shamay]:https://www.twitter.com/meir_shamay