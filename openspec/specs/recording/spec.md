# Recording & Upload

## Purpose

WebRTC stream recording with automatic upload to S3/MinIO object storage for later playback.

## Requirements

### Requirement: Video Recording

The system SHALL record video tracks to IVF format when recording is enabled.

#### Scenario: VP8 recording
- **WHEN** publisher sends VP8 video track and RECORD_ENABLED=1
- **THEN** system records to IVF file with naming {room}_{trackID}_{timestamp}.ivf

#### Scenario: VP9 recording
- **WHEN** publisher sends VP9 video track and RECORD_ENABLED=1
- **THEN** system records to IVF file with naming {room}_{trackID}_{timestamp}.ivf

#### Scenario: Recording disabled
- **WHEN** RECORD_ENABLED=0 (default)
- **THEN** system does not record any tracks

### Requirement: Audio Recording

The system SHALL record audio tracks to OGG format when recording is enabled.

#### Scenario: Opus recording
- **WHEN** publisher sends Opus audio track and RECORD_ENABLED=1
- **THEN** system records to OGG file at 48kHz stereo

### Requirement: S3 Upload

The system SHALL upload recordings to S3/MinIO when configured.

#### Scenario: Upload enabled
- **WHEN** UPLOAD_RECORDINGS=1 and S3 credentials configured
- **THEN** system uploads recording files to configured S3 bucket

#### Scenario: Upload with prefix
- **WHEN** S3_PREFIX is configured
- **THEN** uploaded objects use the prefix in their key

#### Scenario: Upload failure
- **WHEN** S3 upload fails
- **THEN** system logs error and preserves local file

### Requirement: Local File Management

The system SHALL manage recording files in configurable local directory.

#### Scenario: Custom record directory
- **WHEN** RECORD_DIR is configured
- **THEN** system stores recordings in specified directory

#### Scenario: Delete after upload
- **WHEN** DELETE_RECORDING_AFTER_UPLOAD=1 and upload succeeds
- **THEN** system deletes local recording file

### Requirement: Recording Query

The system SHALL provide API to list recording files.

#### Scenario: List recordings
- **WHEN** client calls GET /api/records
- **THEN** system returns list of recordings with metadata (room, trackID, filename, size, createdAt)
