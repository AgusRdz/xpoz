package updater

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// installRelease downloads, verifies, and installs the named release tag.
func installRelease(ctx context.Context, tag string) error {
	archive := platformArchive()

	checksums, err := downloadRelease(ctx, tag, "checksums.txt")
	if err != nil {
		return err
	}
	sigHex, err := downloadSig(ctx, tag)
	if err != nil {
		return err
	}
	if err := verifySignature(checksums, sigHex); err != nil {
		return err
	}
	expected, err := findChecksum(checksums, archive)
	if err != nil {
		return err
	}
	data, err := downloadRelease(ctx, tag, archive)
	if err != nil {
		return err
	}
	if err := verifyChecksum(data, expected); err != nil {
		return err
	}
	return extractAndReplace(data, archive)
}

func downloadSig(ctx context.Context, tag string) (string, error) {
	data, err := downloadRelease(ctx, tag, "checksums.txt.sig")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// platformArchive returns the archive filename for the current OS and arch.
func platformArchive() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}
	return fmt.Sprintf("xpoz-%s-%s%s", goos, goarch, ext)
}

// extractAndReplace extracts the binary from archive data and replaces the current binary.
func extractAndReplace(data []byte, archive string) error {
	binaryName := "xpoz"
	if runtime.GOOS == "windows" {
		binaryName = "xpoz.exe"
	}

	var tempPath string
	var err error
	if strings.HasSuffix(archive, ".zip") {
		tempPath, err = extractFromZip(data, binaryName)
	} else {
		tempPath, err = extractFromTarGz(data, binaryName)
	}
	if err != nil {
		return err
	}
	defer os.Remove(tempPath)

	current, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding current binary: %w", err)
	}
	return applyUpdate(tempPath, current)
}

func extractFromTarGz(data []byte, binaryName string) (string, error) {
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("opening gzip: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("reading tar: %w", err)
		}
		if filepath.Base(hdr.Name) == binaryName {
			return writeTempBinary(tr)
		}
	}
	return "", fmt.Errorf("binary %q not found in archive", binaryName)
}

func extractFromZip(data []byte, binaryName string) (string, error) {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("opening zip: %w", err)
	}
	for _, f := range r.File {
		if filepath.Base(f.Name) == binaryName {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("opening %s in zip: %w", f.Name, err)
			}
			defer rc.Close()
			return writeTempBinary(rc)
		}
	}
	return "", fmt.Errorf("binary %q not found in zip", binaryName)
}

func writeTempBinary(r io.Reader) (string, error) {
	tmp, err := os.CreateTemp("", "xpoz-update-*")
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}
	defer tmp.Close()
	if _, err := io.Copy(tmp, r); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("writing temp file: %w", err)
	}
	return tmp.Name(), nil
}

// applyUpdate replaces dst with src. On Unix this is an atomic rename.
// On Windows, src is staged as dst+".new" for ApplyPending on next startup.
func applyUpdate(src, dst string) error {
	if runtime.GOOS != "windows" {
		if err := os.Chmod(src, 0755); err != nil {
			return fmt.Errorf("chmod: %w", err)
		}
		return os.Rename(src, dst)
	}
	staged := dst + ".new"
	if err := os.Rename(src, staged); err != nil {
		return fmt.Errorf("staging update: %w", err)
	}
	fmt.Fprintf(os.Stderr, "  Update staged — restart xpoz to apply.\n")
	return nil
}
