package url

import (
	"live-url/page/service/url/extractors"
	"net/url"
	"strings"
	"errors"
)

func GetLiveUrl(uri string) (*extractors.LiveInfo, error) {
	var liveDo extractors.LiveInterface

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(u.Hostname(), ".")
	domain := parts[len(parts)-2]
	switch strings.ToLower(domain) {
	case "douyu":
		liveDo = new(extractors.LiveDouyu)
	case "huya":
		liveDo = new(extractors.LiveHuya)
	default:
		return nil, errors.New("")
	}

	return liveDo.Do(uri)
}
