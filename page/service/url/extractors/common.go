package extractors

type PlatformEnum int32

const (
	_PlatformNone PlatformEnum = 0
	PlatformDouyu PlatformEnum = 1
)

type LiveInfo struct {
	Platform  PlatformEnum
	RoomName  string
	OwnerName string
	RealUrl   string
}

type LiveInterface interface {
	Do(string) (*LiveInfo, error)
}
