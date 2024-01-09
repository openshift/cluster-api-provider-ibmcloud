package iamgo

import (
	"net/url"

	"github.com/liamg/jfather"
)

func Parse(policy []byte) (*Document, error) {
	if len(policy) > 0 && policy[0] == '%' {
		decoded, err := url.QueryUnescape(string(policy))
		if err != nil {
			return nil, err
		}
		policy = []byte(decoded)
	}
	var doc Document
	if err := jfather.Unmarshal(policy, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func ParseString(policy string) (*Document, error) {
	return Parse([]byte(policy))
}
