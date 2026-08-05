package release

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"strings"
	"time"
)

const maxArchiveFileBytes = 128 << 20

type sourceEntry struct {
	header tar.Header
	data   []byte
}

func buildSourceArchive(ctx context.Context, config Config) ([]byte, error) {
	raw, err := gitBytes(ctx, config.Root, "archive", "--format=tar", config.Commit)
	if err != nil {
		return nil, fmt.Errorf("archive source commit: %w", err)
	}
	prefix := "starter-otel_" + strings.TrimPrefix(config.Version, "v") + "/"
	entries, err := readGitArchive(raw, prefix, config.Epoch)
	if err != nil {
		return nil, err
	}
	return writeSourceArchive(entries, config.Epoch)
}

func readGitArchive(data []byte, prefix string, epoch time.Time) ([]sourceEntry, error) {
	reader := tar.NewReader(bytes.NewReader(data))
	var result []sourceEntry
	for {
		header, err := reader.Next()
		if errors.Is(err, io.EOF) {
			if len(result) == 0 {
				return nil, errors.New("git source archive contains no files")
			}
			return result, nil
		}
		if err != nil {
			return nil, fmt.Errorf("read Git source archive: %w", err)
		}
		entry, include, entryErr := readSourceEntry(reader, header, prefix, epoch)
		if entryErr != nil {
			return nil, entryErr
		}
		if include {
			result = append(result, entry)
		}
	}
}

func readSourceEntry(
	reader io.Reader,
	header *tar.Header,
	prefix string,
	epoch time.Time,
) (sourceEntry, bool, error) {
	if header.Typeflag == tar.TypeDir || header.Typeflag == tar.TypeXGlobalHeader || header.Typeflag == tar.TypeXHeader {
		return sourceEntry{}, false, nil
	}
	if pathErr := validateArchivePath(header.Name); pathErr != nil {
		return sourceEntry{}, false, pathErr
	}
	if header.Size < 0 || header.Size > maxArchiveFileBytes {
		return sourceEntry{}, false, fmt.Errorf("source entry %q has invalid size %d", header.Name, header.Size)
	}
	entry := sourceEntry{header: normalizedHeader(*header, prefix, epoch)}
	switch header.Typeflag {
	case tar.TypeReg:
		data, readErr := io.ReadAll(io.LimitReader(reader, maxArchiveFileBytes+1))
		if readErr != nil {
			return sourceEntry{}, false, fmt.Errorf("read source entry %q: %w", header.Name, readErr)
		}
		if int64(len(data)) != header.Size {
			return sourceEntry{}, false, fmt.Errorf("read source entry %q: got %d bytes, require %d", header.Name, len(data), header.Size)
		}
		entry.data = data
	case tar.TypeSymlink:
		if linkErr := validateSymlink(header.Name, header.Linkname); linkErr != nil {
			return sourceEntry{}, false, linkErr
		}
	default:
		return sourceEntry{}, false, fmt.Errorf("source entry %q uses unsupported tar type %d", header.Name, header.Typeflag)
	}
	return entry, true, nil
}

func normalizedHeader(source tar.Header, prefix string, epoch time.Time) tar.Header {
	mode := int64(0o644)
	if source.Typeflag == tar.TypeSymlink {
		mode = 0o777
	} else if source.Mode&0o111 != 0 {
		mode = 0o755
	}
	return tar.Header{
		Name:       prefix + path.Clean(source.Name),
		Linkname:   source.Linkname,
		Size:       source.Size,
		Mode:       mode,
		ModTime:    epoch,
		AccessTime: epoch,
		ChangeTime: epoch,
		Typeflag:   source.Typeflag,
		Format:     tar.FormatPAX,
	}
}

func validateArchivePath(name string) error {
	clean := path.Clean(name)
	if name == "" || clean == "." || clean == ".." || strings.HasPrefix(clean, "../") ||
		path.IsAbs(clean) || strings.Contains(name, "\\") {
		return fmt.Errorf("source archive contains unsafe path %q", name)
	}
	return nil
}

func validateSymlink(name, target string) error {
	if target == "" || path.IsAbs(target) || strings.Contains(target, "\\") {
		return fmt.Errorf("source symlink %q has unsafe target %q", name, target)
	}
	resolved := path.Clean(path.Join(path.Dir(name), target))
	if resolved == ".." || strings.HasPrefix(resolved, "../") {
		return fmt.Errorf("source symlink %q escapes archive root", name)
	}
	return nil
}

func writeSourceArchive(entries []sourceEntry, epoch time.Time) ([]byte, error) {
	var output bytes.Buffer
	gzipWriter, err := gzip.NewWriterLevel(&output, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("construct source gzip writer: %w", err)
	}
	gzipWriter.ModTime = epoch
	gzipWriter.OS = 255
	tarWriter := tar.NewWriter(gzipWriter)
	for _, entry := range entries {
		if err := tarWriter.WriteHeader(&entry.header); err != nil {
			return nil, errors.Join(err, tarWriter.Close(), gzipWriter.Close())
		}
		if _, err := tarWriter.Write(entry.data); err != nil {
			return nil, errors.Join(err, tarWriter.Close(), gzipWriter.Close())
		}
	}
	if err := tarWriter.Close(); err != nil {
		return nil, errors.Join(err, gzipWriter.Close())
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, fmt.Errorf("close source gzip writer: %w", err)
	}
	return output.Bytes(), nil
}
