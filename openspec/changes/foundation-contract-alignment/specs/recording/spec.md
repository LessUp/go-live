# Delta Spec: Recording

## MODIFIED Requirements

### Requirement: Recording Query — Response Model

The system SHALL return recording metadata using the following canonical JSON field names
when responding to `GET /api/records`:

| Field     | Type    | Description |
|-----------|---------|-------------|
| `name`    | string  | Filename of the recording (e.g. `room_trackID_1700000000.ivf`) |
| `size`    | integer | File size in bytes at query time |
| `modTime` | string  | File last-modified timestamp in RFC 3339 format, UTC timezone |
| `url`     | string  | Relative URL path to download the file (e.g. `/records/filename`) |

> **Replaces:** the base spec's scenario "returns list of recordings with metadata
> (room, trackID, filename, size, createdAt)" which described a schema that was never
> implemented. The actual implementation emits `name`, `size`, `modTime`, `url`.

#### Scenario: Recording list returns canonical fields
- **WHEN** client calls `GET /api/records` and one or more recording files exist in `RECORD_DIR`
- **THEN** system returns a `200 OK` JSON array where each entry has the fields `name`, `size`, `modTime`, and `url`

#### Scenario: Recording list is empty when directory missing
- **WHEN** client calls `GET /api/records` and `RECORD_DIR` does not exist
- **THEN** system returns `200 OK` with an empty JSON array `[]`

#### Scenario: Recording list excludes non-media files
- **WHEN** `RECORD_DIR` contains files with extensions other than `.ivf` and `.ogg`
- **THEN** those files are excluded from the response array

#### Scenario: Recording list sorted newest-first
- **WHEN** multiple recordings exist
- **THEN** entries are sorted by `modTime` descending; ties are broken by `name` ascending
