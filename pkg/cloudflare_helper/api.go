package cloudflare_helper

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier-base/db/models"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"github.com/WQGroup/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"time"
)

// CloudFlareHelper 5000 人每天下载，每个人之多 66 次
type CloudFlareHelper struct {
	s3Client         *s3.Client
	preSignClient    *s3.PresignClient
	cloudFlareConfig settings.CloudFlareConfig
}

func NewCloudFlareHelper(cloudFlareConfig settings.CloudFlareConfig) *CloudFlareHelper {

	c := CloudFlareHelper{
		cloudFlareConfig: cloudFlareConfig,
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cloudFlareConfig.AccountId),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cloudFlareConfig.AccessKeyId,
				cloudFlareConfig.AccessKeySecret, "")),
	)
	if err != nil {
		logger.Panicln(err)
	}

	c.s3Client = s3.NewFromConfig(cfg)
	c.preSignClient = s3.NewPresignClient(c.s3Client)

	return &c
}

func (c CloudFlareHelper) UploadFile(houseKeepingConfig settings.HouseKeepingConfig, subtitleInfo *models.SubtitleInfo) error {

	body, err := subtitleInfo.GetSubtitleData(houseKeepingConfig.SubsSaveRootDirPath)
	if err != nil {
		return err
	}
	r2StoreKey := subtitleInfo.R2StoreKey()
	_, err = c.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(c.cloudFlareConfig.BucketName),
		Key:    aws.String(r2StoreKey),
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		return err
	}
	return nil
}

func (c CloudFlareHelper) GenerateDownloadUrl(subtitleInfo *models.SubtitleInfo) (string, error) {

	r2StoreKey := subtitleInfo.R2StoreKey()
	downloadTTL := c.cloudFlareConfig.DownloadFileTTL
	if downloadTTL <= 0 {
		downloadTTL = 60
	} else if downloadTTL >= 240 {
		downloadTTL = 240
	}

	preSignedHTTPRequest, err := c.preSignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(c.cloudFlareConfig.BucketName),
		Key:    aws.String(r2StoreKey),
	}, func(options *s3.PresignOptions) {
		options.Expires = time.Duration(downloadTTL) * time.Second
	})
	if err != nil {
		return "", err
	}
	//// 替换原有的域名
	//if settings.Get().CloudFlareConfig.DomainAccess == "" {
	//	return "", fmt.Errorf("CloudFlareConfig.DomainAccess is empty")
	//}
	//u, err := url.Parse(preSignedHTTPRequest.URL)
	//if err != nil {
	//	return "", err
	//}
	//// 这个域名的信息，需要在 R2 的具体一个桶里面的 Settings 去设置
	//u.Host = settings.Get().CloudFlareConfig.DomainAccess

	return preSignedHTTPRequest.URL, nil
}

func (c CloudFlareHelper) DeleteAllFile() error {

	times := 0
	for {
		// 删除一个桶里面所有的文件
		result, err := c.s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket: aws.String(c.cloudFlareConfig.BucketName),
		})
		if err != nil {
			return err
		}

		if len(result.Contents) > 0 {

			times++
			logger.Infoln(times, "Will Delete", len(result.Contents), "files")

			objectKeys := make([]string, 0, len(result.Contents))
			for _, obj := range result.Contents {
				objectKeys = append(objectKeys, *obj.Key)
			}

			var objectIds []types.ObjectIdentifier
			for _, key := range objectKeys {
				objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
			}
			logger.Infoln("Try Delete")
			_, err = c.s3Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
				Bucket: aws.String(c.cloudFlareConfig.BucketName),
				Delete: &types.Delete{Objects: objectIds},
			})
			if err != nil {
				return err
			}

			logger.Infoln("Delete OK")
		} else {
			break
		}
	}

	return nil
}
