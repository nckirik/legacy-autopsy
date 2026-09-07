package canonical

import "testing"

func TestBasicMarkdownFingerprint(t *testing.T) {
	got, err := BasicMarkdownFingerprint("# Title\r\n\r\nText   \r\n")
	if err != nil {
		t.Fatal(err)
	}
	want := "sha256:594f76ae6a9c8fe694cef00288a9812fe09d58fecb75f3e619490684558c49e0"
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
