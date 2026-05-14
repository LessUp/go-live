# Recording

Built-in recording with optional S3/MinIO upload.

## Recording Architecture

```mermaid
flowchart TB
    subgraph Source
        P[Publisher<br/>RTP Packets]
    end

    subgraph SFU
        TF[TrackFanout] -->|Video| VREC[Video Recorder<br/>IVF Writer]
        TF -->|Audio| AREC[Audio Recorder<br/>OGG Writer]
    end

    subgraph Local
        VREC --> FS[(Filesystem<br/>RECORD_DIR)]
        AREC --> FS
    end

    subgraph Upload
        FS -->|UPLOAD_RECORDINGS=1| S3C[S3 Client]
        S3C --> BUCKET[(S3/MinIO)]
    end

    P -->|RTP| TF
```

## Configuration

```bash
# Enable recording
RECORD_ENABLED=1

# Recording directory
RECORD_DIR=records

# Upload to S3/MinIO
UPLOAD_RECORDINGS=1
DELETE_RECORDING_AFTER_UPLOAD=1

# S3 Configuration
S3_ENDPOINT=s3.amazonaws.com
S3_REGION=us-east-1
S3_BUCKET=my-recordings
S3_ACCESS_KEY=$AWS_ACCESS_KEY_ID
S3_SECRET_KEY=$AWS_SECRET_ACCESS_KEY
S3_USE_SSL=1
S3_PREFIX=recordings/
```

## Supported Codecs

| Codec | Container | Extension |
|-------|-----------|-----------|
| VP8 | IVF | `.ivf` |
| VP9 | IVF | `.ivf` |
| Opus | OGG | `.ogg` |

## File Naming

```
{room}_{trackID}_{timestamp}.{ext}
```

Example:
```
demo_video0_1710123456.ivf
demo_audio1_1710123456.ogg
```

## Recording Flow

```mermaid
sequenceDiagram
    participant TF as TrackFanout
    participant R as Recorder
    participant FS as Filesystem
    participant S3 as S3/MinIO

    Note over TF: Publisher connects<br/>Recording starts

    loop Every RTP Packet
        TF->>R: WriteRTP(packet)
        R->>FS: Write to file
    end

    Note over TF: Publisher disconnects

    TF->>R: Close()
    R->>FS: Close file
    
    alt UPLOAD_RECORDINGS=1
        TF->>S3: PutObject(bucket, key, file)
        
        alt DELETE_RECORDING_AFTER_UPLOAD=1
            TF->>FS: Remove local file
        end
    end
```

## S3/MinIO Configuration

### AWS S3

```bash
S3_ENDPOINT=s3.amazonaws.com
S3_REGION=us-east-1
S3_BUCKET=my-bucket
S3_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE
S3_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
S3_USE_SSL=1
```

### MinIO

```bash
S3_ENDPOINT=minio.example.com:9000
S3_BUCKET=recordings
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_USE_SSL=0
S3_PATH_STYLE=1  # Required for MinIO
```

### Object Key Format

```
{S3_PREFIX}{room}_{trackID}_{timestamp}.{ext}
```

Example with `S3_PREFIX=recordings/`:
```
recordings/demo_video0_1710123456.ivf
```

## Recording Lifecycle

```mermaid
stateDiagram-v2
    [*] --> Idle: RECORD_ENABLED=0
    [*] --> Ready: RECORD_ENABLED=1
    
    Ready --> Recording: Publisher Connects
    Recording --> Finalizing: Publisher Disconnects
    Finalizing --> Uploading: UPLOAD_RECORDINGS=1
    Finalizing --> Ready: UPLOAD_RECORDINGS=0
    
    Uploading --> Cleaning: Upload Success
    Uploading --> Ready: Upload Failed
    
    Cleaning --> Ready: DELETE_RECORDING_AFTER_UPLOAD=1
    Cleaning --> Ready: DELETE_RECORDING_AFTER_UPLOAD=0
```

## API Endpoints

### List Recordings

```http
GET /api/records
```

Response:
```json
[
  {
    "name": "demo_video0_1710123456.ivf",
    "size": 1048576,
    "modTime": "2024-03-10T12:34:56Z",
    "url": "/records/demo_video0_1710123456.ivf"
  }
]
```

### Download Recording

```http
GET /records/{filename}
```

## Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| No recordings | `RECORD_ENABLED=0` | Set to `1` |
| Empty files | Codec not supported | Use VP8/VP9/Opus |
| Permission denied | Cannot write to `RECORD_DIR` | Check directory permissions |
| Upload fails | Invalid S3 credentials | Verify `S3_*` configuration |

### Debug Commands

```bash
# Check recording directory
ls -la records/

# Verify S3 configuration
echo $S3_ENDPOINT $S3_BUCKET

# Test S3 connection (AWS CLI)
aws s3 ls --endpoint-url http://minio:9000
```

## Storage Considerations

### File Size Estimation

| Resolution | Bitrate | 1 Hour Size |
|------------|---------|-------------|
| 720p | 2 Mbps | ~900 MB |
| 1080p | 4 Mbps | ~1.8 GB |
| 4K | 15 Mbps | ~6.75 GB |

### Disk Space Management

- Monitor `RECORD_DIR` disk usage
- Implement retention policy (e.g., delete after 7 days)
- Use `DELETE_RECORDING_AFTER_UPLOAD=1` for S3 storage
