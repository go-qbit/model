package model

import "context"

type modelCtxType int8

var derivableFieldsCtx modelCtxType = 0

func initDerivableFieldsCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, derivableFieldsCtx, make(map[string]interface{}))
}

func SetDerivableFieldsData(ctx context.Context, key string, data interface{}) {
	ctx.Value(derivableFieldsCtx).(map[string]interface{})[key] = data
}

func GetDerivableFieldsData(ctx context.Context, key string) interface{} {
	return ctx.Value(derivableFieldsCtx).(map[string]interface{})[key]
}
