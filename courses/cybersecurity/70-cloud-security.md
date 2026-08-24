# Chapter 70: Cloud Security — Securing AWS, GCP, and Azure

*Cloud misconfigurations are the #1 cause of data breaches today. An S3 bucket open to the internet, an overly-permissive IAM role, or exposed metadata service can expose millions of records.*

---

## Cloud Security Model: Shared Responsibility

```
Traditional:       Cloud (IaaS):        Cloud (SaaS):
You own:           You own:             You own:
- Hardware         - OS & runtime       - Data
- Network          - Middleware         - Access control
- OS               - Application
- Application      - Data
- Data             
                   Provider owns:       Provider owns:
                   - Hardware           - Everything else
                   - Network
                   - Hypervisor
```

---

## AWS Security Fundamentals

### IAM (Identity and Access Management)

```json
// IAM Policy — least privilege example
// BAD: Admin access
{
    "Effect": "Allow",
    "Action": "*",
    "Resource": "*"
}

// GOOD: Specific S3 read-only
{
    "Version": "2012-10-17",
    "Statement": [{
        "Effect": "Allow",
        "Action": [
            "s3:GetObject",
            "s3:ListBucket"
        ],
        "Resource": [
            "arn:aws:s3:::my-specific-bucket",
            "arn:aws:s3:::my-specific-bucket/*"
        ]
    }]
}
```

```bash
# AWS CLI security commands
# Find overly-permissive IAM policies
aws iam list-policies --scope Local
aws iam get-policy-version --policy-arn ARN --version-id v1

# Find public S3 buckets
aws s3api list-buckets
aws s3api get-bucket-acl --bucket BUCKET_NAME
aws s3api get-bucket-policy --bucket BUCKET_NAME

# Find unused access keys (rotate or delete)
aws iam list-access-keys --user-name username
aws iam get-access-key-last-used --access-key-id KEY_ID

# AWS Security Hub — aggregated findings
aws securityhub get-findings
```

---

## Common AWS Misconfigurations

### SSRF → Metadata Service → Credential Theft

```bash
# If a web app has SSRF and runs on EC2:
curl http://169.254.169.254/latest/meta-data/iam/security-credentials/
# Returns: EC2 role name

curl http://169.254.169.254/latest/meta-data/iam/security-credentials/my-role
# Returns:
{
    "AccessKeyId": "ASIA...",
    "SecretAccessKey": "wJalrXUtn...",
    "Token": "AQoXnyc4PI...",
    "Expiration": "2024-01-01T00:00:00Z"
}
# These temporary credentials have the EC2 role's permissions!

# Defense: use IMDSv2 (requires PUT to get token first, SSRF can't do this)
aws ec2 modify-instance-metadata-options --instance-id i-xxx --http-tokens required
```

### Public S3 Bucket

```bash
# Check if bucket is public
aws s3api get-bucket-acl --bucket company-backups 2>/dev/null
# If "EVERYONE" has READ access → data breach

# Bug bounty finds: often contain:
# - Database backups (.sql, .dump)
# - Source code
# - Customer data CSVs
# - Environment files (.env)

# Automated scanning
gitleaks detect --source s3://company-bucket
```

### Wildcard IAM Policies

```bash
# Find roles that can assume any role
aws iam list-roles | jq '.Roles[] | select(.AssumeRolePolicyDocument.Statement[].Principal.AWS == "*")'

# Find roles with s3:* or iam:*
aws iam simulate-principal-policy \
    --policy-source-arn arn:aws:iam::123456789:role/my-role \
    --action-names "s3:*" "iam:CreateUser" \
    --resource-arns "*"
```

---

## CloudTrail — AWS Audit Logging

```bash
# Enable CloudTrail (if not enabled — big red flag)
aws cloudtrail create-trail --name my-trail \
    --s3-bucket-name my-audit-bucket \
    --is-multi-region-trail

# Investigate suspicious API calls
aws cloudtrail lookup-events \
    --lookup-attributes AttributeKey=EventName,AttributeValue=ConsoleLogin \
    --start-time 2025-01-01T00:00:00Z

# Find who created/deleted IAM users
aws cloudtrail lookup-events \
    --lookup-attributes AttributeKey=EventName,AttributeValue=CreateUser

# Find someone calling ListBuckets (enumeration)
aws cloudtrail lookup-events \
    --lookup-attributes AttributeKey=EventName,AttributeValue=ListBuckets
```

---

## Go: AWS Security Scanner

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3SecurityFinding struct {
    Bucket   string
    Issue    string
    Severity string
}

func scanS3Buckets(ctx context.Context) []S3SecurityFinding {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    client := s3.NewFromConfig(cfg)
    
    // List all buckets
    buckets, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
    if err != nil {
        log.Fatal(err)
    }
    
    var findings []S3SecurityFinding
    
    for _, bucket := range buckets.Buckets {
        name := *bucket.Name
        
        // Check public access block settings
        pab, err := client.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{
            Bucket: &name,
        })
        
        if err != nil {
            // If GetPublicAccessBlock fails with NoSuchPublicAccessBlockConfiguration,
            // the bucket has no public access block = potentially public
            findings = append(findings, S3SecurityFinding{
                Bucket:   name,
                Issue:    "No public access block configured",
                Severity: "HIGH",
            })
            continue
        }
        
        config := pab.PublicAccessBlockConfiguration
        if config != nil {
            if !*config.BlockPublicAcls {
                findings = append(findings, S3SecurityFinding{
                    Bucket:   name,
                    Issue:    "BlockPublicAcls is disabled",
                    Severity: "HIGH",
                })
            }
            if !*config.BlockPublicPolicy {
                findings = append(findings, S3SecurityFinding{
                    Bucket:   name,
                    Issue:    "BlockPublicPolicy is disabled",
                    Severity: "HIGH",
                })
            }
        }
        
        // Check if bucket has versioning enabled (important for integrity)
        versioning, _ := client.GetBucketVersioning(ctx, &s3.GetBucketVersioningInput{
            Bucket: &name,
        })
        if versioning != nil && versioning.Status != types.BucketVersioningStatusEnabled {
            findings = append(findings, S3SecurityFinding{
                Bucket:   name,
                Issue:    "Versioning not enabled (ransomware risk)",
                Severity: "MEDIUM",
            })
        }
    }
    
    return findings
}

func main() {
    ctx := context.Background()
    findings := scanS3Buckets(ctx)
    
    if len(findings) == 0 {
        fmt.Println("No S3 security issues found")
        return
    }
    
    for _, f := range findings {
        fmt.Printf("[%s] %s: %s\n", f.Severity, f.Bucket, f.Issue)
    }
}
```

---

## Cloud Security Tools

```bash
# ScoutSuite — multi-cloud security auditing
scout aws --profile my-profile

# Prowler — AWS security tool
./prowler -M csv -g all

# CloudSploit / AquaSecurity
cloudsploit scan

# Pacu — AWS exploitation framework (pentesting)
pacu
# Pacu > run iam__enum_users_roles_policies_groups
# Pacu > run s3__bucket_finder
```

---

## Security Best Practices

```
IAM:
✓ Enable MFA for all users, especially root
✓ Delete root access keys (never use root for CLI)
✓ Use IAM roles for EC2/Lambda (no hardcoded keys)
✓ Rotate access keys regularly
✓ Use permission boundaries

Networking:
✓ Principle of least privilege on security groups (0.0.0.0/0 only for HTTP/HTTPS)
✓ VPC with private subnets for databases
✓ VPC Flow Logs enabled
✓ WAF in front of web applications

Data:
✓ Encrypt S3 at rest (SSE-S3 or SSE-KMS)
✓ Block all public S3 access
✓ S3 Object Lock for compliance
✓ Enable versioning

Monitoring:
✓ CloudTrail in all regions
✓ AWS Config for configuration compliance
✓ GuardDuty for threat detection
✓ Security Hub for centralized findings
```

---

## Summary

| Risk Area | Common Issue | Detection | Fix |
|-----------|-------------|-----------|-----|
| S3 | Public bucket | S3 Access Analyzer | Block public access |
| IAM | Wildcard policies | IAM Access Advisor | Least privilege |
| EC2 | SSRF → metadata | GuardDuty | IMDSv2 |
| API keys | Long-lived keys | CloudTrail | Rotate, use roles |
| Logs | CloudTrail disabled | AWS Config rule | Enable multi-region trail |

---

## Exercises

1. Set up a free AWS account and run Prowler — fix the top 5 findings
2. Configure a vulnerable S3 bucket in a test environment and verify scanner detects it
3. Enable GuardDuty in a test AWS account and simulate a finding (e.g., `bitcoin mining` activity)
4. Build the Go S3 scanner and extend it to check for bucket logging enabled
