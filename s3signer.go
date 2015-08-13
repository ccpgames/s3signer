package main

import (
    "io"
    "os"
    "log"
    "time"
    "errors"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/mitchellh/goamz/aws"
    "github.com/mitchellh/goamz/s3"
)

type S3Client struct {
    client  *s3.S3
    buckets []s3.Bucket
}

func S3init(region aws.Region) (S3Client) {
    auth, err := aws.EnvAuth()
    if err != nil {
      log.Fatal(err)
    }
    client := s3.New(auth, region)
    buckets, bucketerr := client.ListBuckets()
    if bucketerr != nil {
        log.Fatal(bucketerr)
    }

    return S3Client{
        client: client,
        buckets: buckets.Buckets,
    }
}

func (s3 *S3Client) BucketFileHandler(w http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    bucket := vars["bucket"]
    filename := vars["filename"]

    for i := 0; i < len(s3.buckets); i++ {
        if bucket == s3.buckets[i].Name {

            keys, err := s3.buckets[i].GetBucketContents()
            if err != nil {
                log.Fatal("Error listing bucket: ", err)
            }

            for _, key := range *keys {
                if filename == key.Key {
                    expires := time.Now()
                    expires.Add(time.Duration(5))
                    url := s3.buckets[i].SignedURL(filename, expires)
                    io.WriteString(w, url)
                    return
                }
            }
            return
        }
    }
}

func GetRegion() (region aws.Region, err error) {
    req_region := os.Getenv("AWS_REGION")
    if req_region == "" {
        err = errors.New("missing AWS_REGION environment variable")
        return
    }

    for name, aws_region := range aws.Regions {
        if req_region == name {
            region = aws_region
            return
        }
    }

    err = errors.New("invalid region: " + req_region)
    return
}

func main() {
    region, err := GetRegion()
    if err != nil {
        log.Fatal(err)
    }

    s3 := S3init(region)

    r := mux.NewRouter()
    r.HandleFunc("/{bucket}/{filename}/", s3.BucketFileHandler)
    http.Handle("/", r)

    httperr := http.ListenAndServe(":8080", nil)
    if httperr != nil {
        log.Fatal("http error: ", httperr)
    }
}
