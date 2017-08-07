package simplestorage

import (
	"fmt"
	"io"

	model "FlexerAPI/Model"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

type SimpleStorage struct {
	Client     *storage.Client
	BucketName string
	Bucket     *storage.BucketHandle
	Config     model.Config
	Writer     io.Writer
	Ctx        context.Context
}

func (s *SimpleStorage) CreateFile(fileName string, file []byte) {
	fmt.Fprintf(s.Writer, "Creating file /%v/%v\n", s.BucketName, fileName)

	wc := s.Bucket.Object(fileName).NewWriter(s.Ctx)
	wc.ContentType = "text/plain"
	wc.Metadata = map[string]string{
		"x-goog-project-id": "1089819306042",
	}

	if _, err := wc.Write(file); err != nil {
		fmt.Println("createFile: unable to write data to bucket %q, file %q: %v", s.BucketName, fileName, err)
		return
	}
	if err := wc.Close(); err != nil {
		fmt.Println("createFile: unable to close bucket %q, file %q: %v", s.BucketName, fileName, err)
		return
	}
}
