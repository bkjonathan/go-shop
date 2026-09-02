package providers

import (
	"context"
	"fmt"
	"mime"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appconfig "github.com/bkjonathan/go-shop/internal/config"
)

type S3Provider struct {
	client   *s3.Client
	bucket   string
	endpoint string
	region   string
}

func NewS3Provider(cfg *appconfig.Config) *S3Provider {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AWS.AccessKeyId, cfg.AWS.SecretAccessKey, "")),
	)
	if err != nil {
		panic("failed to load AWS config: " + err.Error())
	}

	endpoint := strings.TrimSuffix(cfg.AWS.S3Endpoint, "/")

	// Configure for localstack if endpoint is provided. LocalStack has no
	// per-bucket DNS, so addressing has to be path style.
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = true
		}
	})

	return &S3Provider{
		client:   client,
		bucket:   cfg.AWS.S3Bucket,
		endpoint: endpoint,
		region:   cfg.AWS.Region,
	}
}

func (p *S3Provider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	key := p.objectKey(path)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("open upload %q: %w", file.Filename, err)
	}
	defer src.Close()

	_, err = p.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(p.bucket),
		Key:           aws.String(key),
		Body:          src,
		ContentLength: aws.Int64(file.Size),
		ContentType:   aws.String(contentType(file, key)),
	})
	if err != nil {
		return "", fmt.Errorf("put s3://%s/%s: %w", p.bucket, key, err)
	}

	// Callers store this as the image URL, so hand back something fetchable
	// rather than the bare object key.
	return p.publicURL(key), nil
}

func (p *S3Provider) DeleteFile(path string) error {
	key := p.objectKey(path)

	_, err := p.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("delete s3://%s/%s: %w", p.bucket, key, err)
	}
	return nil
}

// contentType prefers the extension whenever the client sent nothing useful, so
// images are not stored as application/octet-stream and downloaded by browsers
// instead of rendered.
func contentType(file *multipart.FileHeader, key string) string {
	ct := file.Header.Get("Content-Type")
	if ct == "" || ct == "application/octet-stream" {
		if byExt := mime.TypeByExtension(filepath.Ext(key)); byExt != "" {
			return byExt
		}
	}
	if ct == "" {
		return "application/octet-stream"
	}
	return ct
}

// publicURL builds the address an uploaded object is served from: path style
// against a custom endpoint (LocalStack, MinIO), virtual host style on real AWS.
func (p *S3Provider) publicURL(key string) string {
	if p.endpoint != "" {
		return fmt.Sprintf("%s/%s/%s", p.endpoint, p.bucket, key)
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", p.bucket, p.region, key)
}

// objectKey accepts either a plain key or a URL previously returned by
// publicURL, so deleting works on whatever was stored on the record.
func (p *S3Provider) objectKey(path string) string {
	if u, err := url.Parse(path); err == nil && u.Scheme != "" && u.Host != "" {
		path = u.Path
	}
	path = strings.TrimPrefix(path, "/")
	return strings.TrimPrefix(path, p.bucket+"/")
}
