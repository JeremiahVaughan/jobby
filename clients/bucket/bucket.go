package bucket

import (
    "encoding/json"
    "crypto/x509"
    "crypto/ecdsa"
    "context"
    "bytes"
    "sync"
    "fmt"
    "os"
    "io"

    "github.com/JeremiahVaughan/jobby/config" 
    s3_config "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
    "github.com/go-acme/lego/v4/registration"
)

type Client struct {
    m sync.Mutex
    uploader *manager.Uploader
    downloader *manager.Downloader
    serviceName string
    configBucketName string
    bitBunkerBucketName string
    certName string
    certKeyName string
    legoRegistrationFileName string
    legoRegistrationKeyFileName string
    domainsFileName string
    useTestDir bool
    s3Client *s3.Client
}

func New(
    ctx context.Context,
    config config.Bucket,
    serviceName string,
) (*Client, error) {
    cfg, err := s3_config.LoadDefaultConfig(ctx)
    if err != nil {
        return nil, fmt.Errorf("error, when loading default config. Error: %v", err)
    }
    s3Client := s3.NewFromConfig(cfg)
    return &Client{
        uploader: manager.NewUploader(s3Client),
        downloader: manager.NewDownloader(s3Client),
        serviceName: serviceName,
        configBucketName: config.ConfigBunkerBucketName,
        bitBunkerBucketName: config.BitBunkerBucketName,
        legoRegistrationFileName: config.LegoRegistrationFileName,
        legoRegistrationKeyFileName: config.LegoRegistrationKeyFileName,
        certName: config.CertName,
        certKeyName: config.CertKeyName,
        domainsFileName: config.DomainsFileName,
        useTestDir: config.UseTestDir,
        s3Client: s3Client,
    }, nil
}
                                                                     

func (c *Client) UploadLegoRegistrationKeyToBucket(ctx context.Context, privateKey *ecdsa.PrivateKey) error {
	der, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
        return fmt.Errorf("error, when MarshalECPrivateKey() for GetMarshalledKey(). Error: %v", err)
	}
    err = c.UploadToConfigBunker(ctx, der, c.legoRegistrationKeyFileName)
    if err != nil {
        return fmt.Errorf("error, when UploadToConfigBunker() for UploadLegoRegistrationToBucket(). Error: %v", err)
    }
    return nil
}

func (c *Client) UploadLegoRegistrationToBucket(ctx context.Context, reg *registration.Resource) error {
    theBytes, err := json.Marshal(reg)
    if err != nil {
        return fmt.Errorf("error, when encoding registration for UploadLegoRegistrationToBucket(). Error: %v", err)
    }
    err = c.UploadToConfigBunker(ctx, theBytes, c.legoRegistrationFileName)
    if err != nil {
        return fmt.Errorf("error, when UploadToConfigBunker() for UploadLegoRegistrationToBucket(). Error: %v", err)
    }
    return nil
}

func (c *Client) DownloadDomainsFromBucket(ctx context.Context) ([]string, error) {
    theBytes, err := c.DownloadFromConfigBunker(ctx, c.domainsFileName)
    if err != nil {
        return nil, fmt.Errorf("error, when DownloadFromConfigBunker() for DownloadExistingDomainsFromBucket(). Error: %v", err)
    }
    var domains []string
    err = json.Unmarshal(theBytes, &domains)
    if err != nil {
        return nil, fmt.Errorf("error, when decoding for DownloadExistingDomainsFromBucket(). Error: %v", err)
    }
    return domains, nil 
}

func (c *Client) UploadDomainsToBucket(ctx context.Context, domains []string) error {
    theBytes, err := json.Marshal(domains)
    if err != nil {
        return fmt.Errorf("error, when decoding for UploadExistingDomainsToBucket(). Error: %v", err)
    }
    err = c.UploadToConfigBunker(ctx, theBytes, c.domainsFileName)
    if err != nil {
        return fmt.Errorf("error, when UploadToConfigBunker() for UploadExistingDomainsToBucket(). Error: %v", err)
    }
    return nil 
}


func (c *Client) DownloadLegoRegistrationKeyFromBucket(ctx context.Context) (*ecdsa.PrivateKey, error) {
    theBytes, err := c.DownloadFromConfigBunker(ctx, c.legoRegistrationKeyFileName)
    if err != nil {
        return nil, fmt.Errorf("error, when DownloadFromConfigBunker() for DownloadLegoRegistrationKeyFromBucket(). Error: %v", err)
    }
    key, err := x509.ParseECPrivateKey(theBytes)
    if err != nil {
        return nil, fmt.Errorf("error, when ParseECPrivateKey() for SetMarshalledKey(). Error: %v", err)
    }
    return key, nil 
}

func (c *Client) DownloadLegoRegistrationFromBucket(ctx context.Context) (*registration.Resource, error) {
    theBytes, err := c.DownloadFromConfigBunker(ctx, c.legoRegistrationFileName)
    if err != nil {
        return nil, fmt.Errorf("error, when DownloadFromConfigBunker() for DownloadLegoRegistrationFromBucket(). Error: %v", err)
    }
    var reg registration.Resource
    err = json.Unmarshal(theBytes, &reg)
    if err != nil {
        return nil, fmt.Errorf("error, when decoding registration for DownloadLegoRegistrationFromBucket(). Error: %v", err)
    }
    return &reg, nil
}

func (c *Client) DownloadCertFromConfigBunker(ctx context.Context) ([]byte, error) {
    theBytes, err := c.DownloadFromConfigBunker(ctx, c.certName)
    if err != nil {
        return nil, fmt.Errorf("error, when DownloadFromConfigBunker() for DownloadCertFromConfigBunker(). Error: %v", err)
    }
    return theBytes, nil
}

func (c *Client) DownloadCertKeyFromConfigBunker(ctx context.Context) ([]byte, error) {
    theBytes, err := c.DownloadFromConfigBunker(ctx, c.certKeyName)
    if err != nil {
        return nil, fmt.Errorf("error, when DownloadFromConfigBunker() for DownloadCertKeyFromConfigBunker(). Error: %v", err)
    }
    return theBytes, nil
}

func (c *Client) UploadCertToConfigBunker(ctx context.Context, theBytes []byte) error {
    err := c.UploadToConfigBunker(ctx, theBytes, c.certName)
    if err != nil {
        return fmt.Errorf("error, when UploadToConfigBunker() for UploadCertToConfigBunker(). Error: %v", err)
    }
    return nil
}

func (c *Client) UploadCertKeyToConfigBunker(ctx context.Context, theBytes []byte) error {
    err := c.UploadToConfigBunker(ctx, theBytes, c.certKeyName)
    if err != nil {
        return fmt.Errorf("error, when UploadToConfigBunker() for UploadCertKeyToConfigBunker(). Error: %v", err)
    }
    return nil
}

func applyTestFolder(objectName string) string {
    return fmt.Sprintf("testing/%s", objectName) 
}

func (c *Client) UploadToConfigBunker(                                    
   ctx context.Context,                                              
   data []byte,                                                      
   fileName string,                                                
) error {                                                             
    c.m.Lock()                                                        
    defer c.m.Unlock()                                                
    dataReader := bytes.NewReader(data)                               
    objectName := fmt.Sprintf("%s/%s", c.serviceName, fileName)
    if c.useTestDir {
        objectName = applyTestFolder(objectName)
    }
    object := s3.PutObjectInput{                                      
        Bucket: aws.String(c.configBucketName),                            
        Key: aws.String(objectName),                                  
        Body: dataReader,                                             
        ContentLength: aws.Int64(int64(len(data))), 
    }                                                                 
    _, err := c.uploader.Upload(ctx, &object)                         
    if err != nil {                                                   
        return fmt.Errorf(                                            
            "error, when uploading data from memory. BucketName: %s. Object Name: %s. Error: %v",                                          
            c.configBucketName,                                            
            objectName,                                               
            err,                                                      
        )                                                             
    }                                                                 
    return nil                                                        
}                                                                     

func (c *Client) DownloadFromConfigBunker(
    ctx context.Context,
    objectName string,
) ([]byte, error) {                                                     
    if c.useTestDir {
        objectName = applyTestFolder(objectName)
    }
    configFile := fmt.Sprintf("%s/%s", c.serviceName, objectName)
    objectInput := &s3.GetObjectInput{                                
        Bucket: aws.String(c.configBucketName),                                   
        Key: aws.String(configFile),                                         
    }                                                                 
    result, err := c.s3Client.GetObject(ctx, objectInput)                    
    if err != nil {                                                   
        return nil, fmt.Errorf("error, failed to get object %s from bucket %s. Error: %v", configFile, c.configBucketName, err)                                            
    }                                                                 
    defer result.Body.Close()                                         
    body, err := io.ReadAll(result.Body)                          
    if err != nil {                                                   
        return nil, fmt.Errorf("error, failed to read object body. Error: %v", err) 
    }                                                                 
    return body, nil                                                  
}                                                                     

func (c *Client) UploadFromDisk(
    ctx context.Context,
    fileLocation string,
    objectName string,
) error {
    if c.useTestDir {
        objectName = applyTestFolder(objectName)
    }
    c.m.Lock()
    defer c.m.Unlock()
    file, err := os.Open(fileLocation)
    if err != nil {
        return fmt.Errorf("error, when opening file %s. Error: %v", fileLocation, err)
    }
    defer file.Close()
    object := s3.PutObjectInput{
        Bucket: aws.String(c.bitBunkerBucketName),
        Key: aws.String(objectName),
        Body: file,
    }
    _, err = c.uploader.Upload(ctx, &object) 
    if err != nil {
        return fmt.Errorf(
            "error, when uploading file. BucketName: %s. Object Name: %s. Error: %v",
            c.bitBunkerBucketName,
            fileLocation,
            err,
        )
    }
    return nil
}

func (c *Client) DownloadToDisk(                                            
   ctx context.Context,                                              
   objectName string,                                                
   outputLocation string,                                            
) error {                                                             
    if c.useTestDir {
        objectName = applyTestFolder(objectName)
    }
    c.m.Lock()                                                        
    defer c.m.Unlock()                                                
                                                                      
    outFile, err := os.Create(outputLocation)                         
    if err != nil {                                                   
        return fmt.Errorf("error, when creating output file %s. Error: %v", outputLocation, err)                                      
    }                                                                 
    defer outFile.Close()                                             
                                                                      
    object := &s3.GetObjectInput{                                     
        Bucket: aws.String(c.bitBunkerBucketName),                            
        Key: aws.String(objectName),                                  
    }                                                                 
                                                                      
    _, err = c.downloader.Download(ctx, outFile, object)           
    if err != nil {                                                   
        return fmt.Errorf("error, failed to download file. Error: %v", err)         
    }                                                                 
                                                                      
    return nil                                                        
}
