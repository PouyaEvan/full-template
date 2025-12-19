resource "aws_s3_bucket" "app_storage" {
  bucket = "my-app-storage-bucket-unique-id"

  tags = {
    Name        = "App Storage"
    Environment = "Production"
  }
}

resource "aws_s3_bucket_ownership_controls" "app_storage" {
  bucket = aws_s3_bucket.app_storage.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "app_storage" {
  depends_on = [aws_s3_bucket_ownership_controls.app_storage]

  bucket = aws_s3_bucket.app_storage.id
  acl    = "private"
}
