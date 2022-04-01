package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewBnUseCase, NewPdUseCase, NewCartUseCase, NewShopUseCase, NewCgUseCase)
