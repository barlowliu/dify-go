package dify-go

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "io"

    "github.com/hashicorp/go-retryablehttp"
)

// SendCompletionMessage sends a text completion message to the Dify API.
// It supports both blocking and streaming response modes.
// For streaming, it returns a channel of ChunkChatCompletionResponse.
func (c *Client) SendCompletionMessage(ctx context.Context, reqBody CompletionMessageRequest) (*CompletionResponse, <-chan ChunkChatCompletionResponse, error) {
    url := c.buildURL("/completion-messages")
    var respBody CompletionResponse
    var streamChan <-chan ChunkChatCompletionResponse

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
        streamChan = make(chan ChunkChatCompletionResponse)
        go func() {
            defer close(streamChan)
            resp, err := c.HTTPClient.Do(req)
            if err != nil {
                streamChan <- ChunkChatCompletionResponse{
                    Event:   "error",
                    Message: err.Error(),
                }
                return
            }
            defer resp.Body.Close()

            if resp.StatusCode != 200 {
                var apiErr APIError
                if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
                    streamChan <- ChunkChatCompletionResponse{
                        Event:   "error",
                        Message: fmt.Sprintf("status code: %d", resp.StatusCode),
                    }
                    return
                }
                streamChan <- ChunkChatCompletionResponse{
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
                        streamChan <- ChunkChatCompletionResponse{
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
                var chunk ChunkChatCompletionResponse
                if err := json.Unmarshal([]byte(jsonStr), &chunk); err != nil {
                    streamChan <- ChunkChatCompletionResponse{
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

