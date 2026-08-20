// Copyright 2026 Q-Tech Team
// Licensed under the GNU AGPLv3 License.
// See LICENSE file in the project root for full license information.

package services

import (
	"context"
	"mime/multipart"
	"testing"
	"bytes"
	"os"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

type mockRekognitionClient struct {
	output *rekognition.DetectModerationLabelsOutput
	err    error
}

func (m *mockRekognitionClient) DetectModerationLabels(ctx context.Context, params *rekognition.DetectModerationLabelsInput, optFns ...func(*rekognition.Options)) (*rekognition.DetectModerationLabelsOutput, error) {
	return m.output, m.err
}

func createDummyFile(size int64) *multipart.FileHeader {
	return &multipart.FileHeader{
		Filename: "test.jpg",
		Size:     size,
	}
}

func createTestFile(t *testing.T, content []byte) *multipart.FileHeader {
	// Create a temporary file to mock multipart.FileHeader's Open()
	tmpfile, err := os.CreateTemp("", "test*.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// We can't easily mock multipart.FileHeader.Open() in standard library because it's a struct with private fields.
	// But in fiber tests or unit tests for it, we usually mock the data differently. 
	// For now, testing the mock client directly is hard without a valid multipart file.
	// Let's just create a valid multipart.FileHeader using a multipart.Writer
	
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	part.Write(content)
	writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(1024)
	if err != nil {
		t.Fatal(err)
	}
	return form.File["file"][0]
}

func TestNSFWService_CheckNSFW_Fallback(t *testing.T) {
	service := NewNSFWService(nil)

	// Test case: Fallback without client
	isNSFW, ratio, err := service.CheckNSFW(context.Background(), createDummyFile(1024))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if isNSFW {
		t.Errorf("Expected isNSFW=false for fallback, got true")
	}
	if ratio != 0.10 {
		t.Errorf("Expected ratio=0.10, got %v", ratio)
	}
}

func TestNSFWService_CheckNSFW_FileSizeLimit(t *testing.T) {
	service := &nsfwService{client: &mockRekognitionClient{}}

	// Test case: File too large
	_, _, err := service.CheckNSFW(context.Background(), createDummyFile(6*1024*1024))
	if err == nil || err.Error() != "file too large, max size is 5MB" {
		t.Errorf("Expected file too large error, got %v", err)
	}
}

func TestNSFWService_CheckNSFW_WithMock(t *testing.T) {
	mockOutput := &rekognition.DetectModerationLabelsOutput{
		ModerationLabels: []types.ModerationLabel{
			{
				Name:       aws.String("Explicit Nudity"),
				Confidence: aws.Float32(95.5),
			},
		},
	}
	service := &nsfwService{client: &mockRekognitionClient{output: mockOutput}}

	file := createTestFile(t, []byte("fake image data"))

	isNSFW, ratio, err := service.CheckNSFW(context.Background(), file)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !isNSFW {
		t.Errorf("Expected isNSFW=true, got false")
	}
	if ratio != 0.955 {
		t.Errorf("Expected ratio=0.955, got %v", ratio)
	}
}

func TestNSFWService_CheckNSFW_Safe(t *testing.T) {
	mockOutput := &rekognition.DetectModerationLabelsOutput{
		ModerationLabels: []types.ModerationLabel{
			{
				Name:       aws.String("Weapons"),
				Confidence: aws.Float32(80.0),
			},
		},
	}
	service := &nsfwService{client: &mockRekognitionClient{output: mockOutput}}

	file := createTestFile(t, []byte("fake safe image data"))

	isNSFW, ratio, err := service.CheckNSFW(context.Background(), file)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if isNSFW {
		t.Errorf("Expected isNSFW=false, got true")
	}
	if ratio != 0 {
		t.Errorf("Expected ratio=0, got %v", ratio)
	}
}

func TestNSFWService_CheckNSFW_OpenError(t *testing.T) {
	service := &nsfwService{client: &mockRekognitionClient{}}
	// createDummyFile has no content, so Open() will fail
	_, _, err := service.CheckNSFW(context.Background(), createDummyFile(1024))
	if err == nil {
		t.Errorf("Expected open error, got nil")
	}
}

func TestNSFWService_CheckNSFW_RekognitionError(t *testing.T) {
	importErr := &nsfwService{client: &mockRekognitionClient{err: errors.New("rekognition err")}}
	
	file := createTestFile(t, []byte("fake safe image data"))

	_, _, err := importErr.CheckNSFW(context.Background(), file)
	if err == nil {
		t.Errorf("Expected rekognition error, got nil")
	}
}
