package backend

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/anxiousmodernman/grpcweb-boilerplate/proto/server"
	"github.com/asdine/storm"
)

// Backend should be used to implement the server interface
// exposed by the generated server proto.
type Backend struct {
}

// Ensure struct implements interface
var _ server.BackendServer = (*Backend)(nil)

// We implement this
//type ProxyServer interface {
//	State(context.Context, *StateRequest) (*ProxyState, error)
//	Put(context.Context, *BackendT) (*OpResult, error)
//	Remove(context.Context, *BackendT) (*OpResult, error)
//}

type Proxy struct {
	DB *storm.DB
}

var _ server.ProxyServer = (*Proxy)(nil)

func (p *Proxy) State(context.Context, *server.StateRequest) (*server.ProxyState, error) {
	return nil, nil
}

func (p *Proxy) Put(ctx context.Context, b *server.BackendT) (*server.OpResult, error) {
	var bd BackendData
	err := p.DB.One("Domain", b.Domain, &bd)

	if err != nil {
		if err == storm.ErrNotFound {
			// do nothing, always overwrite
		} else {
			return &server.OpResult{}, errors.New("")
		}
	}
	bd.Domain = b.Domain
	bd.IPs = combine(bd.IPs, b.Ips)

	err = p.DB.Save(&bd)
	if err != nil {
		return nil, fmt.Errorf("save: %v", err)
	}

	resp := &server.OpResult{200, "Ok"}

	return resp, nil
}

func combine(a, b []string) []string {
	// let's pre-allocate enough space
	both := make([]string, 0, len(a)+len(b))
	both = append(both, a...)
	both = append(both, b...)
	sort.Strings(both)
	var val string
	var res []string
	for _, x := range res {
		if strings.TrimSpace(x) == strings.TrimSpace(val) {
			continue
		}
		val = x
		res = append(res, strings.TrimSpace(x))
	}
	return res
}

func (p *Proxy) Remove(context.Context, *server.BackendT) (*server.OpResult, error) { return nil, nil }

// BackendData is our type for the storm ORM. We can define field-level
// constraints and indexes on struct tags.
type BackendData struct {
	ID     int
	Domain string `storm:"unique"`
	IPs    []string
}

func (bd BackendData) PutIP(ip string) error {

	return nil
}
