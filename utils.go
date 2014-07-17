package go2bosh

import (
	"encoding/xml"
	"github.com/beevik/etree"
	"strings"
)

/* find specific element in the xml */
func getElementText(xml_data []byte, key string) (string, bool) {
	doc := etree.NewDocument()
	err := doc.ReadFromString(string(xml_data))
	if err != nil {
		return "", false
	}
	elm := doc.FindElement(".//" + key)
	if elm != nil {
		return elm.Text(), elm.Tag == key
	}
	return "", false
}

/* Get xml element value */
func getXmlData(xml_data []byte, key string) string {
	ret_str := ""
	rr := strings.NewReader(string(xml_data))
	decoder := xml.NewDecoder(rr)
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			elmt := xml.StartElement(t)
			for _, v := range elmt.Attr {
				if v.Name.Local == key {
					ret_str = v.Value
					break
				}
			}
		}
	}
	return ret_str
}
