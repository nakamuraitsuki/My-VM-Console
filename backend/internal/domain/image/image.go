package image

import "github.com/google/uuid"

type ImageID string

type Image struct {
	id          ImageID
	alias       string // "ubuntu/24.04"
	fingerprint string // "sha256:..."
	serverURL   string // "https://images.linuxcontainers.org" (取得元URL)
	protocol    string // "simplestreams" or "incus"
	isPublic    bool
}

func NewImage(
	id ImageID,
	alias, fingerprint, serverURL, protocol string,
	isPublic bool,
) *Image {
	return &Image{
		id:          id,
		alias:       alias,
		fingerprint: fingerprint,
		serverURL:   serverURL,
		protocol:    protocol,
		isPublic:    isPublic,
	}
}

func NewID() ImageID {
	return ImageID("img-" + uuid.New().String())
}

// --- Getter ---
func (i *Image) ID() ImageID         { return i.id }
func (i *Image) Alias() string       { return i.alias }
func (i *Image) Fingerprint() string { return i.fingerprint }
func (i *Image) ServerURL() string   { return i.serverURL }
func (i *Image) Protocol() string    { return i.protocol }
func (i *Image) IsPublic() bool      { return i.isPublic }
