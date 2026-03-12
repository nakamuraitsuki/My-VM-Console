package gateway

import "example.com/m/internal/domain/compute"

type IngressID string

type IngressRoute struct {
	id         IngressID
	subdomain  string // 例: "app.example.com"の "app" 部分
	portName   string // 例: "http", "https"
	targetIP   string
	targetPort int
	instanceID compute.InstanceID // どのインスタンスにルーティングするか
	ownerID    string             // user.UserIDなど
}

// --- Constructor ---
func NewIngressRoute(id IngressID, subdomain, portName, targetIP string, targetPort int, ownerID string, instanceID compute.InstanceID) *IngressRoute {
	return &IngressRoute{
		id:         id,
		subdomain:  subdomain,
		portName:   portName,
		targetIP:   targetIP,
		targetPort: targetPort,
		instanceID: instanceID,
		ownerID:    ownerID,
	}
}

// --- Getter ---
func (r *IngressRoute) ID() IngressID     { return r.id }
func (r *IngressRoute) Subdomain() string { return r.subdomain }
func (r *IngressRoute) TargetIP() string  { return r.targetIP }
func (r *IngressRoute) PortName() string  { return r.portName }
func (r *IngressRoute) TargetPort() int   { return r.targetPort }
func (r *IngressRoute) OwnerID() string   { return r.ownerID }
func (r *IngressRoute) InstanceID() compute.InstanceID { return r.instanceID }
