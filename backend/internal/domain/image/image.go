package image

type ImageID string

type Image struct {
	id          ImageID
	alias       string // 例: "ubuntu/24.04"
	fingerprint string
	isPublic    bool
}
