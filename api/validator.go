package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/YuanData/allegro-trade/util"
)

var validSymbol validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if symbol, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedSymbol(symbol)
	}
	return false
}
