package gomini

import (
	"context"
	"encoding/json"
	"errors"
	"google.golang.org/grpc"
	"github.com/tidwall/sjson"
	"log"
	"net/http"
	"strings"
)

type handler struct {
	server  interface{}
	methods []grpc.MethodDesc
	cache   map[string]grpc.MethodDesc
}

type Router struct {
	models map[string]*handler
}

func (r *Router) Handle(req *http.Request) ([]byte, error) {
	header, js, tokenErr := Req2json(req)
	if tokenErr != nil {
		return nil, tokenErr
	}
	url := req.URL.Path
	log.Println(req.URL.Path, js)
	dec := func(in interface{}) error {
		err := json.Unmarshal([]byte(js), in)
		if err != nil {
			log.Println(err)
		}
		return nil
	}
	urs := strings.Split(strings.ToLower(url), "/")
	method := urs[len(urs)-1]
	service := strings.Join(urs[1:len(urs)-1], ".")
	if h, ok := r.models[service]; ok {
		if fn := h.cache[method].Handler; fn != nil {
			ctx := context.WithValue(context.Background(), "header", header)
			rs, err := fn(h.server, ctx, dec, nil)
			js, _ := json.Marshal(rs)
			js,_ = sjson.DeleteBytes(js,"req")
			return js, err
		}
	}
	return []byte(""), errors.New("服务不存在:" + service)
}

func (r *Router) Register(srv interface{}, service grpc.ServiceDesc) {
	if r.models == nil {
		r.models = make(map[string]*handler)
	}
	h := &handler{
		server:  srv,
		methods: service.Methods,
	}
	r.models[strings.ToLower(service.ServiceName)] = h
	h.cache = make(map[string]grpc.MethodDesc)
	for _, m := range h.methods {
		h.cache[strings.ToLower(m.MethodName)] = m
	}
	log.Println("register", service.ServiceName, "ok")
}
