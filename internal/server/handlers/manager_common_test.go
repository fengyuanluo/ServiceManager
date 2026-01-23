package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fengyuanluo/go-backend-template/internal/service/manager"
	"github.com/gin-gonic/gin"
)

type responseEnvelope struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func decodeResponse(t *testing.T, rec *httptest.ResponseRecorder) responseEnvelope {
	t.Helper()
	var payload responseEnvelope
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	return payload
}

func TestListHandler(t *testing.T) {
	prev := managerExecSource
	managerExecSource = manager.ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		return []byte("ssh.service loaded active running OpenSSH"), nil
	}}
	defer func() { managerExecSource = prev }()

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/?user=root", nil)

	listHandler("systemd")(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status")
	}
	payload := decodeResponse(t, rec)
	if payload.Code != http.StatusOK {
		t.Fatalf("unexpected response code")
	}
	if !bytes.Contains(payload.Data, []byte("ssh")) {
		t.Fatalf("unexpected response data")
	}
}

func TestDetailHandler(t *testing.T) {
	out := "Id=ssh.service\nDescription=OpenSSH\nMainPID=123\nActiveState=active"
	prev := managerExecSource
	managerExecSource = manager.ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		return []byte(out), nil
	}}
	defer func() { managerExecSource = prev }()

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/?name=ssh", nil)

	detailHandler("systemd")(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status")
	}
	payload := decodeResponse(t, rec)
	if !bytes.Contains(payload.Data, []byte("ssh")) {
		t.Fatalf("unexpected response data")
	}
}

func TestLogsHandler(t *testing.T) {
	prev := managerExecSource
	managerExecSource = manager.ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		return []byte("info line\nerror line"), nil
	}}
	defer func() { managerExecSource = prev }()

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/?name=ssh&lines=10&keyword=error", nil)

	logsHandler("systemd")(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status")
	}
	payload := decodeResponse(t, rec)
	if !bytes.Contains(payload.Data, []byte("error")) {
		t.Fatalf("unexpected response data")
	}
}

func TestActionHandler(t *testing.T) {
	prev := managerExecSource
	managerExecSource = manager.ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		return []byte("ok"), nil
	}}
	defer func() { managerExecSource = prev }()

	body := bytes.NewBufferString(`{"name":"ssh","action":"restart"}`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", body)
	c.Request.Header.Set("Content-Type", "application/json")

	actionHandler("systemd")(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status")
	}
}

func TestBatchHandler(t *testing.T) {
	prev := managerExecSource
	managerExecSource = manager.ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		return []byte("ok"), nil
	}}
	defer func() { managerExecSource = prev }()

	body := bytes.NewBufferString(`{"items":[{"name":"ssh","action":"restart"}]}`)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", body)
	c.Request.Header.Set("Content-Type", "application/json")

	batchHandler("systemd")(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status")
	}
}

func TestInstallHandler(t *testing.T) {
	prevExec := managerExecSource
	prevEnv := managerEnvSource
	managerExecSource = manager.ExecSource{Run: func(name string, args ...string) ([]byte, error) {
		return []byte("ok"), nil
	}}
	managerEnvSource = manager.EnvSource{ReadFile: func(path string) ([]byte, error) {
		return []byte("ID=ubuntu\n"), nil
	}}
	defer func() {
		managerExecSource = prevExec
		managerEnvSource = prevEnv
	}()

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")

	installHandler("supervisor")(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status")
	}
}
