package utils

// import (
// 	// "bytes"
// 	"github.com/aliyun/aliyun-oss-go-sdk/oss"
//   // "fmt"
// )

// func OssInit() error{
//   endpoint := "https://oss-cn-guangzhou.aliyuncs.com"
// 	accessKeyID := "LTAI5tA7Rw75kZrJXb9hAzHw"
// 	accessKeySecret := "WWz8NPvgCN6zlp95waucHlFIlzJFyY"
// 	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
//   if err != nil {
//     return err
//   }

//   // 用defer关闭client
//   defer client.Close() 

//   // 获取bucket 
//   bucket, err := client.Bucket("tikfun")
//   if err != nil {
//     return err 
//   }

//   // 记录初始化成功日志
//   log.Printf("OSS client initialized, endpoint: %s, bucket: %s", endpoint, bucket)

//   return nil
// }