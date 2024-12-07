package dify-go

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "io"

    "github.com/hashicorp/go-retryablehttp"
)

// RunWorkflow executes a workflow.
// It supports both blocking and streaming response modes.
// For streaming, it returns a channel of ChunkCompletionResponse.
func (c *Client) RunWorkflow(ctx context.Context, reqBody WorkflowRunRequest) (*WorkflowCompletionResponse, <-chan ChunkCompletionResponse, error) {
    url := c.buildURL("/workflows/run")
    var respBody WorkflowCompletionResponse
    var streamChan <-chan ChunkCompletionResponse

    // Marshal request body
    bodyBytes, err := json.Marshal(reqBody)
    if err != nil {
        return nil, nil, err
    }

    // Create new request
    req, err := retryablehttp.NewRequest("POST", url, bytes.NewReader(bodyBytes))
    if err != nil {
        return nil, nil, err
    }

    // Add headers
    c.addHeaders(req)

    // Determine if streaming
    if reqBody.ResponseMode == "streaming" {
        // Handle streaming
        streamChan = make(chan ChunkCompletionResponse)
        go func() {
            defer close(streamChan)
            resp, err := c.HTTPClient.Do(req)
            if err != nil {
                streamChan <- ChunkCompletionResponse{
                    Event:   "error",
                    Message: err.Error(),
                }
                return
            }
            defer resp.Body.Close()

            if resp.StatusCode != 200 {
                var apiErr APIError
                if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
                    streamChan <- ChunkCompletionResponse{
                        Event:   "error",
                        Message: fmt.Sprintf("status code: %d", resp.StatusCode),
                    }
                    return
                }
                streamChan <- ChunkCompletionResponse{
                    Event:   "error",
                    Code:    apiErr.Code,
                    Message: apiErr.Message,
                }
                return
            }

            reader := bufio.NewReader(resp.Body)
            for {
                line, err := reader.ReadString('\n')
                if err != nil {
                    if err != io.EOF {
                        streamChan <- ChunkCompletionResponse{
                            Event:   "error",
                            Message: err.Error(),
                        }
                    }
                    break
                }

                if len(line) < 6 || line[:6] != "data: " {
                    continue
                }

                jsonStr := line[6:]
                var chunk ChunkCompletionResponse
                if err := json.Unmarshal([]byte(jsonStr), &chunk); err != nil {
                    streamChan <- ChunkCompletionResponse{
                        Event:   "error",
                        Message: err.Error(),
                    }
                    continue
                }

                select {
                case streamChan <- chunk:
                case <-ctx.Done():
                    return
                }
            }
        }()
        return nil, streamChan, nil
    } else if reqBody.ResponseMode == "blocking" {
        // Handle blocking
        resp, err := c.HTTPClient.Do(req)
        if err != nil {
            return nil, nil, err
        }
        defer resp.Body.Close()

        if resp.StatusCode != 200 {
            var apiErr APIError
            if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
                return nil, nil, fmt.Errorf("status code: %d", resp.StatusCode)
            }
            return nil, nil, &apiErr
        }

        if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
            return nil, nil, err
        }

        return &respBody, nil, nil
    } else {
        return nil, nil, fmt.Errorf("invalid response_mode: %s", reqBody.ResponseMode)
    }
}

// GetWorkflowStatus retrieves the status of a workflow execution by its ID.
func (c *Client) GetWorkflowStatus(workflowRunID string) (*WorkflowStatusResponse, error) {
    endpoint := fmt.Sprintf("/workflows/run/%s", workflowRunID)
    url := c.buildURL(endpoint)

    // Create new request
    req, err := retryablehttp.NewRequest("GET", url, nil)
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

    var statusResp WorkflowStatusResponse
    if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
        return nil, err
    }

    return &statusResp, nil
}

