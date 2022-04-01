// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v2.1.4

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

type NtSrvHTTPServer interface {
	CreateNt(context.Context, *CreateNtReq) (*CreateNtReply, error)
	DeleteNt(context.Context, *DeleteNtReq) (*DeleteNtReply, error)
	GetNt(context.Context, *GetNtReq) (*GetNtReply, error)
	ListNt(context.Context, *ListNtReq) (*ListNtReply, error)
	UpdateStatus(context.Context, *UpdateStatusReq) (*UpdateStatusReply, error)
}

func RegisterNtSrvHTTPServer(s *http.Server, srv NtSrvHTTPServer) {
	r := s.Route("/")
	r.POST("/api/notice/create", _NtSrv_CreateNt0_HTTP_Handler(srv))
	r.PUT("/api/notice/update", _NtSrv_UpdateStatus0_HTTP_Handler(srv))
	r.GET("/api/notice/list/{userUuid}/{status}/{limit}/{page}/{type}", _NtSrv_ListNt0_HTTP_Handler(srv))
	r.DELETE("/api/notice/delete/{userUuid}/{id}", _NtSrv_DeleteNt0_HTTP_Handler(srv))
	r.GET("/api/notice/get/{id}/{userUuid}", _NtSrv_GetNt0_HTTP_Handler(srv))
}

func _NtSrv_CreateNt0_HTTP_Handler(srv NtSrvHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CreateNtReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.user.ntSrv/CreateNt")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CreateNt(ctx, req.(*CreateNtReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CreateNtReply)
		return ctx.Result(200, reply)
	}
}

func _NtSrv_UpdateStatus0_HTTP_Handler(srv NtSrvHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateStatusReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.user.ntSrv/UpdateStatus")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateStatus(ctx, req.(*UpdateStatusReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateStatusReply)
		return ctx.Result(200, reply)
	}
}

func _NtSrv_ListNt0_HTTP_Handler(srv NtSrvHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListNtReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.user.ntSrv/ListNt")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListNt(ctx, req.(*ListNtReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListNtReply)
		return ctx.Result(200, reply)
	}
}

func _NtSrv_DeleteNt0_HTTP_Handler(srv NtSrvHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteNtReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.user.ntSrv/DeleteNt")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteNt(ctx, req.(*DeleteNtReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DeleteNtReply)
		return ctx.Result(200, reply)
	}
}

func _NtSrv_GetNt0_HTTP_Handler(srv NtSrvHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetNtReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.user.ntSrv/GetNt")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetNt(ctx, req.(*GetNtReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetNtReply)
		return ctx.Result(200, reply)
	}
}

type NtSrvHTTPClient interface {
	CreateNt(ctx context.Context, req *CreateNtReq, opts ...http.CallOption) (rsp *CreateNtReply, err error)
	DeleteNt(ctx context.Context, req *DeleteNtReq, opts ...http.CallOption) (rsp *DeleteNtReply, err error)
	GetNt(ctx context.Context, req *GetNtReq, opts ...http.CallOption) (rsp *GetNtReply, err error)
	ListNt(ctx context.Context, req *ListNtReq, opts ...http.CallOption) (rsp *ListNtReply, err error)
	UpdateStatus(ctx context.Context, req *UpdateStatusReq, opts ...http.CallOption) (rsp *UpdateStatusReply, err error)
}

type NtSrvHTTPClientImpl struct {
	cc *http.Client
}

func NewNtSrvHTTPClient(client *http.Client) NtSrvHTTPClient {
	return &NtSrvHTTPClientImpl{client}
}

func (c *NtSrvHTTPClientImpl) CreateNt(ctx context.Context, in *CreateNtReq, opts ...http.CallOption) (*CreateNtReply, error) {
	var out CreateNtReply
	pattern := "/api/notice/create"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.user.ntSrv/CreateNt"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *NtSrvHTTPClientImpl) DeleteNt(ctx context.Context, in *DeleteNtReq, opts ...http.CallOption) (*DeleteNtReply, error) {
	var out DeleteNtReply
	pattern := "/api/notice/delete/{userUuid}/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.user.ntSrv/DeleteNt"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *NtSrvHTTPClientImpl) GetNt(ctx context.Context, in *GetNtReq, opts ...http.CallOption) (*GetNtReply, error) {
	var out GetNtReply
	pattern := "/api/notice/get/{id}/{userUuid}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.user.ntSrv/GetNt"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *NtSrvHTTPClientImpl) ListNt(ctx context.Context, in *ListNtReq, opts ...http.CallOption) (*ListNtReply, error) {
	var out ListNtReply
	pattern := "/api/notice/list/{userUuid}/{status}/{limit}/{page}/{type}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.user.ntSrv/ListNt"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *NtSrvHTTPClientImpl) UpdateStatus(ctx context.Context, in *UpdateStatusReq, opts ...http.CallOption) (*UpdateStatusReply, error) {
	var out UpdateStatusReply
	pattern := "/api/notice/update"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.user.ntSrv/UpdateStatus"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
