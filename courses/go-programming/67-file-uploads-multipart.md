# Chapter 67: File Uploads, Multipart Forms, and Static Files

File uploads are a common requirement: profile pictures, document attachments, CSV imports. Go's `net/http` handles multipart forms natively. This chapter covers single and multiple file uploads, size limits, type validation, streaming to object storage, and serving static files.

## Table of Contents

1. [Multipart Form Basics](#1-multipart-form-basics)
2. [File Upload Handler](#2-file-upload-handler)
3. [Validation and Security](#3-validation-and-security)
4. [Streaming to Storage](#4-streaming-to-storage)
5. [Static File Serving](#5-static-file-serving)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. Multipart Form Basics

A multipart form request has `Content-Type: multipart/form-data; boundary=...`. Each part is a separate field or file, separated by the boundary.

```
POST /upload HTTP/1.1
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW

------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="description"

My photo
------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="file"; filename="photo.jpg"
Content-Type: image/jpeg

<binary data>
------WebKitFormBoundary7MA4YWxkTrZu0gW--
```

```go
// Parse a multipart form
// maxMemory = how many bytes to keep in RAM (rest spills to disk)
if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
    http.Error(w, "form too large", http.StatusRequestEntityTooLarge)
    return
}

// Read a non-file field
description := r.FormValue("description")

// Read a single file
file, header, err := r.FormFile("file")
if err != nil {
    http.Error(w, "missing file field", http.StatusBadRequest)
    return
}
defer file.Close()

fmt.Println("filename:", header.Filename)
fmt.Println("size:", header.Size)
fmt.Println("type:", header.Header.Get("Content-Type"))
```

---

## 2. File Upload Handler

```go
package handler

import (
    "fmt"
    "io"
    "mime"
    "net/http"
    "path/filepath"
    "strings"
)

const (
    maxUploadSize = 5 << 20 // 5 MB
)

var allowedTypes = map[string]bool{
    "image/jpeg": true,
    "image/png":  true,
    "image/webp": true,
    "image/gif":  true,
}

type UploadResult struct {
    Filename    string `json:"filename"`
    Size        int64  `json:"size"`
    ContentType string `json:"contentType"`
    URL         string `json:"url"`
}

// HandleUpload handles a single file upload.
func HandleUpload(storage Storage) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Limit total request body size
        r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
        
        if err := r.ParseMultipartForm(maxUploadSize); err != nil {
            if strings.Contains(err.Error(), "http: request body too large") {
                writeError(w, http.StatusRequestEntityTooLarge,
                    fmt.Sprintf("file must be smaller than %d MB", maxUploadSize>>20))
                return
            }
            writeError(w, http.StatusBadRequest, "invalid form")
            return
        }

        file, header, err := r.FormFile("file")
        if err != nil {
            writeError(w, http.StatusBadRequest, "missing file field")
            return
        }
        defer file.Close()

        // Detect content type from file content (not client-supplied header)
        contentType, err := detectContentType(file)
        if err != nil {
            writeError(w, http.StatusBadRequest, "could not detect file type")
            return
        }

        if !allowedTypes[contentType] {
            writeError(w, http.StatusUnsupportedMediaType,
                fmt.Sprintf("unsupported file type %q, allowed: jpeg, png, webp, gif", contentType))
            return
        }

        // Build a safe filename
        ext := extensionForType(contentType)
        safeFilename := sanitizeFilename(header.Filename, ext)

        // Store the file
        url, err := storage.Store(r.Context(), safeFilename, contentType, file)
        if err != nil {
            writeError(w, http.StatusInternalServerError, "failed to store file")
            return
        }

        writeJSON(w, http.StatusCreated, UploadResult{
            Filename:    safeFilename,
            Size:        header.Size,
            ContentType: contentType,
            URL:         url,
        })
    }
}

// HandleMultiUpload handles multiple files in one request.
func HandleMultiUpload(storage Storage, maxFiles int) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        r.Body = http.MaxBytesReader(w, r.Body, int64(maxFiles)*maxUploadSize)
        
        if err := r.ParseMultipartForm(int64(maxFiles) * maxUploadSize); err != nil {
            writeError(w, http.StatusBadRequest, "invalid form")
            return
        }

        files := r.MultipartForm.File["files"]
        if len(files) == 0 {
            writeError(w, http.StatusBadRequest, "no files provided")
            return
        }
        if len(files) > maxFiles {
            writeError(w, http.StatusBadRequest, fmt.Sprintf("max %d files per request", maxFiles))
            return
        }

        results := make([]UploadResult, 0, len(files))
        var errs []string

        for _, fileHeader := range files {
            f, err := fileHeader.Open()
            if err != nil {
                errs = append(errs, fmt.Sprintf("%s: open error", fileHeader.Filename))
                continue
            }

            ct, err := detectContentType(f)
            if err != nil || !allowedTypes[ct] {
                f.Close()
                errs = append(errs, fmt.Sprintf("%s: unsupported type", fileHeader.Filename))
                continue
            }

            ext := extensionForType(ct)
            name := sanitizeFilename(fileHeader.Filename, ext)
            url, err := storage.Store(r.Context(), name, ct, f)
            f.Close()
            if err != nil {
                errs = append(errs, fmt.Sprintf("%s: storage error", fileHeader.Filename))
                continue
            }
            results = append(results, UploadResult{
                Filename: name, Size: fileHeader.Size,
                ContentType: ct, URL: url,
            })
        }

        code := http.StatusCreated
        if len(results) == 0 { code = http.StatusBadRequest }

        writeJSON(w, code, map[string]any{
            "uploaded": results,
            "errors":   errs,
        })
    }
}

// detectContentType reads the first 512 bytes to sniff the content type.
// It then seeks back to the start so the caller can read the full file.
func detectContentType(f interface {
    io.Reader
    io.Seeker
}) (string, error) {
    buf := make([]byte, 512)
    n, err := f.Read(buf)
    if err != nil && err != io.EOF { return "", err }
    if _, err := f.Seek(0, io.SeekStart); err != nil { return "", err }
    return http.DetectContentType(buf[:n]), nil
}

func extensionForType(ct string) string {
    exts, _ := mime.ExtensionsByType(ct)
    if len(exts) > 0 { return exts[0] }
    return ".bin"
}

func sanitizeFilename(original, ext string) string {
    // Strip directory traversal and non-safe characters
    name := filepath.Base(original)
    name = strings.Map(func(r rune) rune {
        if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' ||
            r == '-' || r == '_' || r == '.' {
            return r
        }
        return '_'
    }, name)
    // Force the correct extension based on detected type
    name = strings.TrimSuffix(name, filepath.Ext(name))
    return name + ext
}
```

---

## 3. Validation and Security

### Magic bytes check

The `Content-Type` header sent by the client can be spoofed. Always detect the type from file content:

```go
// For images, verify the magic bytes manually:
var magicBytes = map[string][]byte{
    "image/jpeg": {0xFF, 0xD8, 0xFF},
    "image/png":  {0x89, 0x50, 0x4E, 0x47},
    "image/gif":  {0x47, 0x49, 0x46},
    "image/webp": {0x52, 0x49, 0x46, 0x46}, // "RIFF"
}

func checkMagicBytes(data []byte, expected string) bool {
    magic, ok := magicBytes[expected]
    if !ok { return true } // unknown type, trust http.DetectContentType
    if len(data) < len(magic) { return false }
    for i, b := range magic {
        if data[i] != b { return false }
    }
    return true
}
```

### Size limit per file

```go
// Limit a single file read to maxSize bytes
func limitedRead(src io.Reader, maxSize int64) ([]byte, error) {
    limited := io.LimitReader(src, maxSize+1)
    data, err := io.ReadAll(limited)
    if err != nil { return nil, err }
    if int64(len(data)) > maxSize {
        return nil, fmt.Errorf("file exceeds maximum size of %d bytes", maxSize)
    }
    return data, nil
}
```

---

## 4. Streaming to Storage

Never buffer an entire upload in memory if you can help it. Stream directly to the destination.

```go
// Storage interface — implemented by local disk, S3, GCS, etc.
type Storage interface {
    Store(ctx context.Context, key, contentType string, body io.Reader) (url string, err error)
    Delete(ctx context.Context, key string) error
    URL(key string) string
}

// Local disk implementation
type DiskStorage struct {
    dir     string
    baseURL string
}

func NewDiskStorage(dir, baseURL string) *DiskStorage {
    return &DiskStorage{dir: dir, baseURL: baseURL}
}

func (s *DiskStorage) Store(ctx context.Context, key, contentType string, body io.Reader) (string, error) {
    path := filepath.Join(s.dir, key)
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
        return "", fmt.Errorf("mkdirall: %w", err)
    }
    
    f, err := os.CreateTemp(filepath.Dir(path), ".upload-*")
    if err != nil { return "", fmt.Errorf("create temp: %w", err) }
    
    _, err = io.Copy(f, body)
    f.Close()
    if err != nil {
        os.Remove(f.Name())
        return "", fmt.Errorf("write: %w", err)
    }
    
    if err := os.Rename(f.Name(), path); err != nil {
        os.Remove(f.Name())
        return "", fmt.Errorf("rename: %w", err)
    }
    
    return s.URL(key), nil
}

func (s *DiskStorage) Delete(ctx context.Context, key string) error {
    return os.Remove(filepath.Join(s.dir, key))
}

func (s *DiskStorage) URL(key string) string {
    return s.baseURL + "/" + key
}

// S3 streaming upload (sketch — uses multipart upload for large files)
type S3Storage struct {
    client *s3.Client
    bucket string
}

func (s *S3Storage) Store(ctx context.Context, key, contentType string, body io.Reader) (string, error) {
    _, err := s.client.PutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(key),
        Body:        body,           // streamed directly — no full buffer
        ContentType: aws.String(contentType),
    })
    if err != nil { return "", fmt.Errorf("s3 put: %w", err) }
    return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key), nil
}
```

---

## 5. Static File Serving

```go
// Serve static files from a directory
http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

// Serve embedded files (Go 1.16+)
//go:embed static
var staticFiles embed.FS

http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(staticFiles))))

// Single page app — serve index.html for all unmatched routes
type spaHandler struct {
    fs embed.FS
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Try to serve the file
    f, err := h.fs.Open(strings.TrimPrefix(r.URL.Path, "/"))
    if err == nil {
        f.Close()
        http.FileServer(http.FS(h.fs)).ServeHTTP(w, r)
        return
    }
    // Fallback to index.html for client-side routing
    http.ServeFileFS(w, r, h.fs, "index.html")
}

// With chi router:
r := chi.NewRouter()
r.Handle("/api/*", apiRouter)
r.Handle("/*", spaHandler{staticFiles})
```

### Cache headers for static assets

```go
// Serve static files with long-lived cache and content hash in URL
func hashFilename(path string) string {
    data, err := os.ReadFile(path)
    if err != nil { return path }
    h := fmt.Sprintf("%x", sha256.Sum256(data))[:8]
    ext := filepath.Ext(path)
    return strings.TrimSuffix(path, ext) + "." + h + ext
}

// Middleware to add Cache-Control headers
func cacheControl(maxAge int) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d, immutable", maxAge))
            next.ServeHTTP(w, r)
        })
    }
}

// Static files: cache for 1 year (365*24*60*60 = 31536000 seconds)
// API responses: no-cache
r.With(cacheControl(31536000)).Handle("/static/*", staticHandler)
```

---

## Summary

- Parse multipart forms with `r.ParseMultipartForm(maxBytes)` — memory-bounds the form
- **Always detect content type from file bytes** (`http.DetectContentType`), not from `Content-Type` header
- Use `http.MaxBytesReader` to reject oversized bodies before parsing
- Stream files to storage with `io.Copy` — never buffer the entire file in memory for large uploads
- Sanitize filenames: strip path separators, dangerous characters, and force the correct extension
- Serve embedded static files with `embed.FS` for single-binary deployments
- For SPAs: serve `index.html` as a fallback for unmatched routes

## Exercises

### Easy
1. Write a handler that accepts a CSV file upload, parses the CSV rows, and returns a JSON array of the rows.
2. Add a `downloads` endpoint that serves user-uploaded files from disk with the correct `Content-Disposition: attachment; filename=...` header so browsers download rather than display them.
3. Implement a file upload size meter: return progress updates as SSE events while storing the file. Use `io.TeeReader` to count bytes as they're written.

### Medium
4. Build an **image resizer**: on upload, generate three sizes (thumb 100×100, medium 400×400, full) using `golang.org/x/image` or a pure-Go imaging library. Store all three, return URLs for each.
5. Implement **resumable uploads**: break large files into chunks (each < 5MB). The client uploads each chunk with `Content-Range: bytes 0-4999999/total`. Server assembles chunks in order. Use a temp file per upload ID.
6. Add **virus scanning** integration: after receiving the file bytes, submit them to a local ClamAV socket or a cloud API. Reject the upload and delete the file if threats are detected.

### Hard
7. Implement **server-side multipart streaming**: instead of `ParseMultipartForm` (which buffers), use `r.MultipartReader()` to process parts one at a time as they arrive. Stream each file part directly to S3 using a multipart upload. This allows handling arbitrarily large files with bounded memory.
8. Build a **presigned URL system**: instead of proxying file uploads through your server, generate time-limited S3/GCS presigned PUT URLs. The client uploads directly to object storage. Your server only receives a notification callback when upload is complete and verifies the file was actually stored before updating the database.
