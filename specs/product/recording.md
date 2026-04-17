# Recording & Upload

## Overview

The project implements WebRTC stream recording with automatic upload to S3/MinIO object storage.

---

## User Stories

### As a Stream Administrator
- I want to record publisher streams for later playback
- I want recordings to be automatically uploaded to S3/MinIO
- I want to configure local file cleanup after upload
- I want to list available recordings via API

### As a Developer
- I want recordings in standard formats (IVF for video, OGG for audio)
- I want recording to not impact streaming performance
- I want graceful handling of recording failures

---

## Requirements

### Functional Requirements

1. **Video Recording**
   - VP8/VP9 → IVF format
   - Per-track recording
   - File naming: `{room}_{trackID}_{timestamp}.{ext}`

2. **Audio Recording**
   - Opus → OGG format
   - 48kHz, stereo

3. **S3/MinIO Upload**
   - Configurable via `UPLOAD_RECORDINGS`
   - S3 endpoint, region, bucket, credentials
   - SSL support (`S3_USE_SSL`)
   - Path-style addressing (`S3_PATH_STYLE`)
   - Object key prefix (`S3_PREFIX`)

4. **Local File Management**
   - Configurable output directory (`RECORD_DIR`)
   - Optional deletion after upload (`DELETE_RECORDING_AFTER_UPLOAD`)

---

## Acceptance Criteria

1. ✅ Video tracks are recorded to IVF files
2. ✅ Audio tracks are recorded to OGG files
3. ✅ Recordings are uploaded to S3 when enabled
4. ✅ Local files are deleted after upload when configured
5. ✅ Recordings can be listed via API
6. ✅ Recording can be enabled/disabled via `RECORD_ENABLED`

---

## Edge Cases

1. **Server Crash**: Recording files should be recoverable
2. **Upload Failure**: Retry logic or error logging
3. **Disk Space**: Monitor and handle disk exhaustion
4. **Concurrent Writes**: Multiple tracks recording simultaneously
5. **Large Files**: Handle long-running recordings

---

## Configuration Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `RECORD_ENABLED` | `0` | Enable recording (`1` to enable) |
| `RECORD_DIR` | `records` | Recording output directory |
| `UPLOAD_RECORDINGS` | `0` | Enable S3 upload (`1` to enable) |
| `DELETE_RECORDING_AFTER_UPLOAD` | `0` | Delete local file after upload |
| `S3_ENDPOINT` | - | S3/MinIO endpoint |
| `S3_REGION` | - | S3 region |
| `S3_BUCKET` | - | Target bucket name |
| `S3_ACCESS_KEY` | - | Access Key ID |
| `S3_SECRET_KEY` | - | Secret Access Key |
| `S3_USE_SSL` | `1` | Use HTTPS for S3 connection |
| `S3_PATH_STYLE` | `0` | Use path-style addressing |
| `S3_PREFIX` | - | Object key prefix |

---

## Out of Scope

- HLS/DASH streaming from recordings
- Video transcoding
- Recording playback via web interface
