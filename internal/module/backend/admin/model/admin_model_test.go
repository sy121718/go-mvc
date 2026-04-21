package adminmodel

import (
	"context"
	"testing"
)

func TestAllowModifyIsAdminDoesNotPanicWhenFlagMissing(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("allowModifyIsAdmin 不应 panic: %v", r)
		}
	}()

	if allowModifyIsAdmin(context.Background()) {
		t.Fatalf("缺少内部标记时不应允许修改")
	}
}

func TestAllowModifyIsAdminReturnsTrueWhenFlagEnabled(t *testing.T) {
	ctx := context.WithValue(context.Background(), allowModifyIsAdminKey, true)
	if !allowModifyIsAdmin(ctx) {
		t.Fatalf("存在内部标记时应允许修改")
	}
}
