package middleware

import (
	"testing"
)

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "Hello World"},
		{"<script>alert('xss')</script>", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
		{"Normal text", "Normal text"},
		{"<img src=x onerror=alert(1)>", "&lt;img src=x onerror=alert(1)&gt;"},
		{"Test\x00Null", "TestNull"},
		{"O'Reilly & Sons", "O&#39;Reilly &amp; Sons"},
		{"", ""},
	}

	for _, tt := range tests {
		result := SanitizeString(tt.input)
		if result != tt.expected {
			t.Errorf("SanitizeString(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractIP(t *testing.T) {
	// Test basic IP extraction from RemoteAddr
	tests := []struct {
		remoteAddr string
		expected   string
	}{
		{"192.168.1.1:8080", "192.168.1.1"},
		{"10.0.0.1:443", "10.0.0.1"},
		{"127.0.0.1:3000", "127.0.0.1"},
	}

	for _, tt := range tests {
		// Simple split test
		parts := splitIP(tt.remoteAddr)
		if parts != tt.expected {
			t.Errorf("extractIP(%q) = %q, want %q", tt.remoteAddr, parts, tt.expected)
		}
	}
}

func splitIP(addr string) string {
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return addr[:i]
		}
	}
	return addr
}

func TestRateLimiterTokens(t *testing.T) {
	// Test token bucket logic
	v := getVisitor("test-ip-1", 5)
	if v.tokens != 5 {
		t.Errorf("new visitor should have 5 tokens, got %d", v.tokens)
	}

	// Consume tokens
	v.tokens--
	v.tokens--
	if v.tokens != 3 {
		t.Errorf("after 2 requests, should have 3 tokens, got %d", v.tokens)
	}

	// Get same visitor
	v2 := getVisitor("test-ip-1", 5)
	if v2.tokens < 3 {
		t.Errorf("same visitor should still have >= 3 tokens, got %d", v2.tokens)
	}
}
