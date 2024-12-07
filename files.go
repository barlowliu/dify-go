package dify-go

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "mime/multipart"
    "net/http"
    "os"

    "github.com/hashicorp/go-retryablehttp"
)

// UploadFile uploads a file to the Dify API.
// Returns the uploaded file's information.
func (c *Client) UploadFile(ctx context.Context, filePath, user string) (*FileUploadResponse, error) {
    url := c.buildURL("/files/upload")

    // Open the file
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Create a buffer to write our multipart form
    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)

    // Add the file
    part, err := writer.CreateFormFile("file", filePath)
    if err != nil {
        return nil, err
    }
    _, err = io.Copy(part, file)
    if err != nil {
        return nil, err
    }

    // Add the user field
    if err := writer.WriteField("user", user); err != nil {
        return nil, err
    }

    // Close the writer to finalize the multipart form
    if err := writer.Close(); err != nil {
        return nil, err
    }

    // Create new request
    req, err := retryablehttp.NewRequest("POST", url, &requestBody)
    if err != nil {
        return nil, err
    }

    // Add headers
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
    req.Header.Set("Content-Type", writer.FormDataContentType())

    // Execute request
    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        var apiErr APIError
        if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
            return nil, fmt.Errorf("status code: %d", resp.StatusCode)
        }
        return nil, &apiErr
    }

    var uploadResp FileUploadResponse
    if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
        return nil, err
    }

    return &uploadResp, nil
}

