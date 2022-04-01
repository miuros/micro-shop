package router

import "micro-shop/internal/middleware"

var rules=[]middleware.Rule{
	{
		Key: "/api/user/",
		Limit: 500,
		Per: 50,
	},
	{
		Key: "/api/comment/",
		Limit: 500,
		Per: 50,
	},
	{
		Key: "/api/notice",
		Limit: 500,
		Per: 100,
	},
	{
		Key: "/api/product",
		Limit: 500,
		Per: 100,
	},
	{
		Key: "/api/order",
		Limit: 100,
		Per: 10,
	},
	{
		Key: "/api/shop",
		Limit: 500,
		Per: 50,
	},
	{
		Key: "/api/cart",
		Limit: 200,
		Per: 10,
	},
}

func NewLimiter()*middleware.RouterLimiter {
	rl:=middleware.NewRouterLimiter()
	rl.AddLimiter(rules...)
	return rl
}
