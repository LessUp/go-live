# 录制功能

内置录制，支持可选 S3/MinIO 上传。

## 录制架构

```mermaid
flowchart TB
    subgraph Source
        P[发布者<br/>RTP 包]
    end

    subgraph SFU
        TF[TrackFanout] -->|视频| VREC[视频录制器<br/>IVF 写入器]
        TF -->|音频| AREC[音频录制器<br/>OGG 写入器]
    end

    subgraph Local
        VREC --> FS[(文件系统<br/>RECORD_DIR)]
        AREC --> FS
    end

    subgraph Upload
        FS -->|UPLOAD_RECORDINGS=1| S3C[S3 客户端]
        S3C --> BUCKET[(S3/MinIO)]
    end

    P -->|RTP| TF
```

## 配置

```bash
# 启用录制
RECORD_ENABLED=1

# 录制目录
RECORD_DIR=records

# 上传到 S3/MinIO
UPLOAD_RECORDINGS=1
DELETE_RECORDING_AFTER_UPLOAD=1

# S3 配置
S3_ENDPOINT=s3.amazonaws.com
S3_REGION=us-east-1
S3_BUCKET=my-recordings
S3_ACCESS_KEY=$AWS_ACCESS_KEY_ID
S3_SECRET_KEY=$AWS_SECRET_ACCESS_KEY
S3_USE_SSL=1
S3_PREFIX=recordings/
```

## 支持的编解码器

| 编解码器 | 容器 | 扩展名 |
|----------|------|--------|
| VP8 | IVF | `.ivf` |
| VP9 | IVF | `.ivf` |
| Opus | OGG | `.ogg` |

## 文件命名

```
{room}_{trackID}_{timestamp}.{ext}
```

示例：
```
demo_video0_1710123456.ivf
demo_audio1_1710123456.ogg
```

## S3/MinIO 配置

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
S3_PATH_STYLE=1  # MinIO 必需
```

## 故障排除

| 问题 | 原因 | 解决方案 |
|------|------|----------|
| 无录制文件 | `RECORD_ENABLED=0` | 设置为 `1` |
| 空文件 | 编解码器不支持 | 使用 VP8/VP9/Opus |
| 权限拒绝 | 无法写入 `RECORD_DIR` | 检查目录权限 |
| 上传失败 | S3 凭证无效 | 验证 `S3_*` 配置 |

### 调试命令

```bash
# 检查录制目录
ls -la records/

# 验证 S3 配置
echo $S3_ENDPOINT $S3_BUCKET

# 测试 S3 连接 (AWS CLI)
aws s3 ls --endpoint-url http://minio:9000
```

## 存储考虑

### 文件大小估算

| 分辨率 | 比特率 | 1 小时大小 |
|--------|--------|------------|
| 720p | 2 Mbps | ~900 MB |
| 1080p | 4 Mbps | ~1.8 GB |
| 4K | 15 Mbps | ~6.75 GB |

### 磁盘空间管理

- 监控 `RECORD_DIR` 磁盘使用
- 实现保留策略（如 7 天后删除）
- 使用 `DELETE_RECORDING_AFTER_UPLOAD=1` 用于 S3 存储
