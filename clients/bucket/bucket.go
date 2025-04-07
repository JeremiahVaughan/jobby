package bucket

import (
    "context"
    "sync"
    "fmt"
    "os"

    "github.com/JeremiahVaughan/jobby/config" 
    s3_config "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

type Client struct {
    m sync.Mutex
    uploader *manager.Uploader
    config config.Bucket
}

func New(ctx context.Context, config config.Bucket) (*Client, error) {
    cfg, err := s3_config.LoadDefaultConfig(ctx)
    if err != nil {
        return nil, fmt.Errorf("error, when loading default config. Error: %v", err)
    }
    c := s3.NewFromConfig(cfg)
    uploader := manager.NewUploader(c)
    return &Client{
        uploader: uploader,
        config: config,
    }, nil
}

func (c *Client) BackupFile(
    ctx context.Context,
    fileName string,
    objectName string,
) error {
    c.m.Lock()
    defer c.m.Unlock()
    file, err := os.Open(fileName)
    if err != nil {
        return fmt.Errorf("error, when opening file %s. Error: %v", fileName, err)
    }
    defer file.Close()
    object := s3.PutObjectInput{
        Bucket: aws.String(c.config.Name),
        Key: aws.String(objectName),
        Body: file,
    }
    _, err = c.uploader.Upload(ctx, &object) 
    if err != nil {
        return fmt.Errorf(
            "error, when uploading file. BucketName: %s. Object Name: %s. Error: %v",
            c.config.Name,
            fileName,
            err,
        )
    }
    return nil
}
