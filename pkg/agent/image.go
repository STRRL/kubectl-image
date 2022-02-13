package agent

import "time"

type ContainerImage struct {
	Repository string
	Tag        string
	Digest     string
	ImageID    string
	Created    time.Time
	Size       int64
}
