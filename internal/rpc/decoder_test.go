package rpc

import (
	"testing"
)

func TestDecodeBatchResponseNullResult(t *testing.T) {
	body := ")]}'\n\n105\n" +
		`[["wrb.fr","cFji9",null,null,null,[3],"generic"],["di",136],["af.httprm",136,"7219101052150421406",47]]` +
		"\n25\n" +
		`[["e",4,null,null,141]]`

	raw, err := DecodeBatchResponse(body, "cFji9")
	if err != nil {
		t.Fatalf("expected no error for null result, got: %v", err)
	}
	if raw != nil {
		t.Fatalf("expected nil raw for null result, got: %s", string(raw))
	}
	t.Log("PASS: null result decoded as success")
}

func TestDecodeBatchResponseStringResult(t *testing.T) {
	body := ")]}'\n\n100\n" +
		`[["wrb.fr","wXbhsf","[[\"test\",null,\"id-123\"]]",null,null,null,"generic"]]` +
		"\n25\n" +
		`[["di",72]]`

	raw, err := DecodeBatchResponse(body, "wXbhsf")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if raw == nil {
		t.Fatal("expected non-nil result")
	}
	t.Logf("PASS: result=%s", string(raw))
}
