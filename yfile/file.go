package yfile

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

// IsDir checks if the path is directory or not.
func IsDir(n string) bool {
	file, err := os.Stat(n)
	if err != nil {
		return false
	}
	return file.IsDir()
}

// IsExist checks if a file or directory exists.
func IsExist(n string) bool {
	_, err := os.Stat(n)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

// IsLink checks if a file is symbol link or not.
func IsLink(n string) bool {
	fi, err := os.Lstat(n)
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeSymlink != 0
}

// ModTime returns file modified time in seconds.
func ModTime(n string) (int64, error) {
	f, err := os.Stat(n)
	if err != nil {
		return 0, err
	}
	return f.ModTime().Unix(), nil
}

// FileSize returns file size in bytes.
func FileSize(n string) (int64, error) {
	f, err := os.Stat(n)
	if err != nil {
		return 0, err
	}
	return f.Size(), nil
}

// FileCopy copies a file.
func FileCopy(src, dst string) error {
	if !IsExist(src) {
		return errors.New("src does not exist")
	} else if IsExist(dst) {
		return errors.New("dst already exists")
	}
	srcFile, srcErr := os.Open(src)
	if srcErr != nil {
		return srcErr
	}
	defer srcFile.Close()
	dstFile, dstErr := os.Create(dst)
	if dstErr != nil {
		return dstErr
	}
	defer dstFile.Close()
	buf := make([]byte, 4<<10)
	var (
		cnt int
		err error
	)
	for {
		if cnt, err = srcFile.Read(buf); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		} else if _, err = dstFile.Write(buf[:cnt]); err != nil {
			return err
		}
	}
}

// DirCopy copies a directory.
func DirCopy(src, dst string) error {
	if srcStat, err := os.Stat(src); err != nil {
		return err
	} else if !srcStat.IsDir() {
		return errors.New("src is not a directory")
	} else if err = os.MkdirAll(dst, srcStat.Mode()); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf(`cannot read source directory: %w`, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if err = DirCopy(path.Join(src, entry.Name()), path.Join(dst, entry.Name())); err != nil {
				return err
			}
		} else {
			if err = FileCopy(path.Join(src, entry.Name()), path.Join(dst, entry.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}

// ClearFile clears a file.
func ClearFile(n string) error {
	return os.WriteFile(n, []byte{}, 0644)
}

// MIMEType returns the MIME type of file.
func MIMEType(v any) (string, error) {
	buf := make([]byte, 512)
	switch vv := v.(type) {
	case string:
		if f, e := os.Open(vv); e != nil {
			return "", e
		} else {
			defer f.Close()
			if _, e = f.Read(buf); e != nil {
				return "", e
			}
		}
	case []byte:
		_, _ = io.CopyN(bytes.NewBuffer(buf[:0]), bytes.NewReader(vv), 512)
	case io.Reader:
		if _, e := vv.Read(buf); e != nil {
			return "", e
		}
	default:
		return "", errors.New("invalid argument")
	}
	return http.DetectContentType(buf), nil
}
