package identicon

import (
	"image/png"
	"os"
	"testing"
)

func TestGenerateIdenticon(t *testing.T) {
	f, err := os.Create("identicon.png")
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	defer os.Remove(f.Name())
	f.Close()

	if err := GenerateIdenticon(12345, f.Name()); err != nil {
		t.Fatalf("GenerateIdenticon returned error: %v", err)
	}

	imgFile, err := os.Open(f.Name())
	if err != nil {
		t.Fatalf("failed to open generated file: %v", err)
	}
	defer imgFile.Close()

	img, err := png.Decode(imgFile)
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 420 || bounds.Dy() != 420 {
		t.Fatalf("unexpected image size: %dx%d", bounds.Dx(), bounds.Dy())
	}
}
