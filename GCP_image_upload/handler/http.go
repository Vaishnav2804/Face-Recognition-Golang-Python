package Handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

type handlehttp struct {
}

func New() *handlehttp {
	return &handlehttp{}
}

const (
	projectID  = "face-recognition-311111" // Project ID
	bucketName = "fr_vaishnav"             // Bucket_Name
)

var uploader *ClientUploader

type responseErr struct {
	StatusCode int    `json:"code"`
	Err        string `json:"message"`
}

type responseSuccess struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

func init() {
	_ = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/Users/vaishnav/Downloads/face-recognition-311111-277732e769d4.json") // JSON PATH TO VALIDATE
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	uploader = &ClientUploader{
		cl:         client,
		bucketName: bucketName,
		projectID:  projectID,
	}
}

func (c ClientUploader) UploadFile(file multipart.File, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func (h handlehttp) Upload(w http.ResponseWriter, r *http.Request) {
	_, header1, err := r.FormFile("pic1") // for pic 1
	if err != nil {
		resp := responseErr{StatusCode: 500, Err: "wrong header/wrong file format"}
		finalResp, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(finalResp)
		return
	}

	_, header2, err := r.FormFile("pic2") // for pic 2
	if err != nil {
		resp := responseErr{StatusCode: 500, Err: "wrong header/wrong file format"}
		finalResp, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(finalResp)
		return
	}

	fmt.Println(header1.Filename)
	fmt.Println(header2.Filename)

	blobFile, err := header1.Open()
	if err != nil {
		resp := responseErr{StatusCode: 500, Err: "Cannot input file"}
		finalResp, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(finalResp)
		return
	}

	err = uploader.UploadFile(blobFile, header1.Filename)
	if err != nil {
		resp := responseErr{StatusCode: 500, Err: "Error in uploading file to GCP"}
		finalResp, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(finalResp)
		return
	}

	blobFile2, err := header2.Open()
	if err != nil {
		resp := responseErr{StatusCode: 500, Err: "Cannot input file"}
		finalResp, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(finalResp)
		return
	}

	err = uploader.UploadFile(blobFile2, header2.Filename)
	if err != nil {
		resp := responseErr{StatusCode: 500, Err: "Error in uploading file to GCP"}
		finalResp, _ := json.Marshal(resp)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(finalResp)
		return
	}

	resp := responseSuccess{StatusCode: 200, Message: "Success"}
	finalResp, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(finalResp)
	return
}
