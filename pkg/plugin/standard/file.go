package standard

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"strings"

	blade "github.com/oligarch316/go-sickle-blade"
	"go.uber.org/zap"
)

const (
	fileDataHeadSize             = 512
	fileOpenFlag                 = os.O_WRONLY | os.O_CREATE | os.O_EXCL
	fileOpenMode     fs.FileMode = 0755
)

type fileBufferedData struct {
	head bytes.Buffer
	io.Reader
}

func (fbd *fileBufferedData) from(reader io.Reader) error {
	if _, err := io.CopyN(&fbd.head, reader, fileDataHeadSize); err != nil && err != io.EOF {
		return err
	}

	fbd.Reader = io.MultiReader(&fbd.head, reader)
	return nil
}

type (
	FileConsumableName      interface{ StdFileName() string }
	FileConsumableExtension interface{ StdFileExtension() string }
)

type FileConsumer struct{ logger blade.Logger }

func (FileConsumer) Namespace() string { return "file" }

func (FileConsumer) extensionFor(item blade.Item, dataHead []byte) (string, bool) {
	if consumable, ok := item.(FileConsumableExtension); ok {
		if ext := consumable.StdFileExtension(); ext != "" {
			if ext[0] != '.' {
				ext = "." + ext
			}

			return ext, true
		}
	}

	if dataHead != nil {
		contentType := http.DetectContentType(dataHead)
		extList, err := mime.ExtensionsByType(contentType)
		if err == nil && len(extList) != 0 {
			return extList[0], true
		}
	}

	return "", false
}

func (FileConsumer) basenameFor(item blade.Item) (string, error) {
	if consumable, ok := item.(FileConsumableName); ok {
		if name := consumable.StdFileName(); name != "" {
			return name, nil
		}
	}

	fingerprint, err := item.Fingerprint()
	if err != nil {
		return "", err
	}

	if len(fingerprint) > 10 {
		fingerprint = fingerprint[:10]
	}

	return base64.RawURLEncoding.EncodeToString(fingerprint), nil
}

func (fc FileConsumer) collectionNameFor(item blade.CollectionItem) (string, error) {
	basename, err := fc.basenameFor(item)
	if err != nil {
		return "", err
	}

	ext, ok := fc.extensionFor(item, nil)
	if !ok {
		ext = ".collection"
	}

	return basename + ext, nil
}

func (fc FileConsumer) mediaNameFor(item blade.MediaItem, dataHead []byte) (string, error) {
	basename, err := fc.basenameFor(item)
	if err != nil {
		return "", err
	}

	ext, ok := fc.extensionFor(item, dataHead)
	if !ok {
		ext = ".media"
	}

	return basename + ext, err
}

func (fc FileConsumer) writeFile(name string, data io.Reader) error {
	logger := fc.logger.With(zap.String("filename", name))

	file, openErr := os.OpenFile(name, fileOpenFlag, fileOpenMode)
	if openErr != nil {
		if openErr == fs.ErrExist {
			logger.Warn("file name already exists, skipping")
			return nil
		}
		return openErr
	}

	if _, copyErr := io.Copy(file, data); copyErr != nil {
		if closeErr := file.Close(); closeErr != nil {
			logger.Warn(
				"failed to close file after copy error, manual cleanup may be required",
				zap.NamedError("closeError", closeErr),
			)

			return copyErr
		}

		if removeErr := os.Remove(name); removeErr != nil {
			logger.Warn(
				"failed to remove file after copy error, manual cleanup may be required",
				zap.NamedError("removeError", removeErr),
			)
		}

		return copyErr
	}

	if closeErr := file.Close(); closeErr != nil {
		logger.Warn("failed to close file, manual cleanup may be required", zap.Error(closeErr))
		return closeErr
	}

	return nil
}

func (fc FileConsumer) ConsumeCollection(ctx context.Context, item blade.CollectionItem) error {
	fc.logger.Info("received collection item", zap.String("type", fmt.Sprintf("%T", item)))

	urls, err := item.ChildURLs(ctx)
	switch {
	case err != nil:
		return err
	case len(urls) == 0:
		fc.logger.Warn("collection contains no urls, skipping")
		return nil
	}

	filename, err := fc.collectionNameFor(item)
	if err != nil {
		return err
	}

	urlStrs := make([]string, len(urls))
	for i, url := range urls {
		urlStrs[i] = url.String()
	}

	data := strings.NewReader(strings.Join(urlStrs, "\n") + "\n")

	if err := fc.writeFile(filename, data); err != nil {
		return err
	}

	fc.logger.Info("wrote collection item to file", zap.String("filename", filename))
	return nil
}

func (fc FileConsumer) ConsumeMedia(ctx context.Context, item blade.MediaItem) error {
	fc.logger.Info("received media item", zap.String("type", fmt.Sprintf("%T", item)))

	itemReader, err := item.Data(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := itemReader.Close(); err != nil {
			fc.logger.Warn("failed to close media item data reader", zap.NamedError("closeError", err))
		}
	}()

	var data fileBufferedData
	if err := data.from(itemReader); err != nil {
		return err
	}

	filename, err := fc.mediaNameFor(item, data.head.Bytes())
	if err != nil {
		return err
	}

	if err := fc.writeFile(filename, data); err != nil {
		return err
	}

	fc.logger.Info("wrote media item to file", zap.String("filename", filename))
	return nil
}
