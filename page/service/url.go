package service

func GetLiveUrl(uri, room string) (string, error) {
	return uri+room, nil
}
