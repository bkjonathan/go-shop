#!/bin/bash
# Runs inside the LocalStack container at the "ready" stage.
# NOTE: the shebang above is required - LocalStack execs this file directly,
# and without it the run fails with "Exec format error" and no bucket is made.
set -euo pipefail

BUCKET="${AWS_S3_BUCKET:-ecommerce-uploads}"
REGION="${DEFAULT_REGION:-us-east-1}"

if awslocal s3api head-bucket --bucket "$BUCKET" 2>/dev/null; then
  echo "LocalStack S3 bucket '$BUCKET' already exists."
else
  awslocal s3 mb "s3://$BUCKET" --region "$REGION"
  echo "LocalStack S3 bucket '$BUCKET' created successfully."
fi

# Uploaded product images are served straight from the bucket URL, so allow
# anonymous reads (LocalStack only - never do this on a real AWS bucket).
awslocal s3api put-bucket-policy --bucket "$BUCKET" --policy "{
  \"Version\": \"2012-10-17\",
  \"Statement\": [{
    \"Sid\": \"PublicRead\",
    \"Effect\": \"Allow\",
    \"Principal\": \"*\",
    \"Action\": \"s3:GetObject\",
    \"Resource\": \"arn:aws:s3:::$BUCKET/*\"
  }]
}"
echo "LocalStack S3 bucket '$BUCKET' is ready."
