package service

/* Abstarct impelemtation of the user data and credential that the bosh/service need in order to authenticate the user */
type ServiceUser interface {
	/* Use to get the user key/password/token to connect to bosh server or/and to xmpp service provider */
	GetUserData() (string, interface{})

	/* The authentication request against the bosh server is made with base64 encode of all the credential details of the request */
	SetAuthRequestData(token string) (string, interface{})
}
