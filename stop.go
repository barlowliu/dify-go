package dify-go

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/go-retryablehttp"
)

// StopTask stops an ongoing stream task by its task_id.
func (c *Client) StopTask(ctx context.Context, taskID, user string) (*StopResponse, error) {
    endpoint := fmt.Sprintf("/chat-messages/%s/stop", taskID)
    url := c.buildURL(endpoint)

    // Prepare request body
    body := map[string]string{
        "user": user,
    }
    bodyBytes, err := json.Marshal(body)
    if err != nil {
        return nil, err
    }

    // Create new request
    req, err := retryablehttp.NewRequest("POST", url, bytes.NewReader(bodyBytes))
    if err != nil {
        return nil, err
    }

    // Add headers
    c.addHeaders(req)

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

    var stopResp StopResponse
    if err := json.NewDecoder(resp.Body).Decode(&stopResp); err != nil {
        return nil, err
    }

    return &stopResp, nil
}

