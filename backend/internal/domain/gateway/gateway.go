package gateway

type IngressID string

type IngressRoute struct {
	id         IngressID
	fqdn       string // 例: "app.example.com"
	targetIP   string
	targetPort int
	ownerID    string // user.UserIDなど
}
