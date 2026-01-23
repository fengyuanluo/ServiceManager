package assets

import "testing"

func TestEmbeddedIndex(t *testing.T) {
	_, err := StaticFS.Open("index.html")
	if err != nil {
		t.Fatalf("expected embedded index")
	}
}
