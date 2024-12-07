package dify-go

import (
    "time"
)

// ChatMessageRequest represents the request body for sending chat messages.
type ChatMessageRequest struct {
    Query           string                 `json:"query"`
    Inputs          map[string]interface{} `json:"inputs,omitempty"`
    ResponseMode    string                 `json:"response_mode"`
    User            string                 `json:"user"`
    ConversationID  string                 `json:"conversation_id,omitempty"`
    Files           []FileUploadInfo       `json:"files,omitempty"`
    AutoGenerateName bool                  `json:"auto_generate_name,omitempty"`
}

// ChatCompletionResponse represents the response for blocking chat messages.
type ChatCompletionResponse struct {
    MessageID      string       `json:"message_id"`
    ConversationID string       `json:"conversation_id"`
    Mode           string       `json:"mode"`
    Answer         string       `json:"answer"`
    Metadata       Metadata     `json:"metadata"`
    CreatedAt      int64        `json:"created_at"`
}

// ChunkChatCompletionResponse represents each chunk in streaming chat messages.
type ChunkChatCompletionResponse struct {
    Event          string                 `json:"event"`
    TaskID         string                 `json:"task_id"`
    MessageID      string                 `json:"message_id"`
    ConversationID string                 `json:"conversation_id"`
    Answer         string                 `json:"answer,omitempty"`
    CreatedAt      int64                  `json:"created_at"`
    // Additional fields for different event types
    Metadata       *Metadata              `json:"metadata,omitempty"`
    Status         int                    `json:"status,omitempty"`
    Code           string                 `json:"code,omitempty"`
    Message        string                 `json:"message,omitempty"`
    Audio          string                 `json:"audio,omitempty"`
    // ... other fields as per API documentation
}

// Metadata contains usage and retriever resources information.
type Metadata struct {
    Usage              Usage               `json:"usage"`
    RetrieverResources []RetrieverResource `json:"retriever_resources"`
}

// Usage represents the usage details in metadata.
type Usage struct {
    PromptTokens       int     `json:"prompt_tokens"`
    PromptUnitPrice    string  `json:"prompt_unit_price"`
    PromptPriceUnit    string  `json:"prompt_price_unit"`
    PromptPrice        string  `json:"prompt_price"`
    CompletionTokens   int     `json:"completion_tokens"`
    CompletionUnitPrice string `json:"completion_unit_price"`
    CompletionPriceUnit string `json:"completion_price_unit"`
    CompletionPrice    string  `json:"completion_price"`
    TotalTokens        int     `json:"total_tokens"`
    TotalPrice         string  `json:"total_price"`
    Currency           string  `json:"currency"`
    Latency            float64 `json:"latency"`
}

// RetrieverResource represents individual retriever resources.
type RetrieverResource struct {
    Position        int     `json:"position"`
    DatasetID       string  `json:"dataset_id"`
    DatasetName     string  `json:"dataset_name"`
    DocumentID      string  `json:"document_id"`
    DocumentName    string  `json:"document_name"`
    SegmentID       string  `json:"segment_id"`
    Score           float64 `json:"score"`
    Content         string  `json:"content"`
}

// FileUploadInfo represents information about uploaded files.
type FileUploadInfo struct {
    Type            string `json:"type"`
    TransferMethod  string `json:"transfer_method"`
    URL             string `json:"url,omitempty"`
    UploadFileID    string `json:"upload_file_id,omitempty"`
}

// WorkflowRunRequest represents the request body for running workflows.
type WorkflowRunRequest struct {
    Inputs       map[string]interface{} `json:"inputs,omitempty"`
    ResponseMode string                 `json:"response_mode"`
    User         string                 `json:"user"`
    Files        []FileUploadInfo       `json:"files,omitempty"`
}

// CompletionMessageRequest represents the request body for text completion.
type CompletionMessageRequest struct {
    Inputs        map[string]interface{} `json:"inputs"`
    ResponseMode  string                 `json:"response_mode"`
    User          string                 `json:"user"`
    Files         []FileUploadInfo       `json:"files,omitempty"`
}

// CompletionResponse represents the response for blocking completion messages.
type CompletionResponse struct {
    ID        string `json:"id"`
    Answer    string `json:"answer"`
    CreatedAt int64  `json:"created_at"`
}

// FileUploadResponse represents the response after uploading a file.
type FileUploadResponse struct {
    ID         string `json:"id"`
    Name       string `json:"name"`
    Size       int    `json:"size"`
    Extension  string `json:"extension"`
    MimeType   string `json:"mime_type"`
    CreatedBy  string `json:"created_by"`
    CreatedAt  int64  `json:"created_at"`
}

// StopResponse represents the response after stopping a task.
type StopResponse struct {
    Result string `json:"result"`
}

// WorkflowCompletionResponse represents the response for blocking workflow runs.
type WorkflowCompletionResponse struct {
    WorkflowRunID string          `json:"workflow_run_id"`
    TaskID        string          `json:"task_id"`
    Data          WorkflowRunData `json:"data"`
}

// WorkflowRunData contains detailed information about workflow execution.
type WorkflowRunData struct {
    ID           string  `json:"id"`
    WorkflowID   string  `json:"workflow_id"`
    Status       string  `json:"status"`
    Outputs      *string `json:"outputs,omitempty"`
    Error        *string `json:"error,omitempty"`
    ElapsedTime  float64 `json:"elapsed_time,omitempty"`
    TotalTokens  int     `json:"total_tokens,omitempty"`
    TotalSteps   int     `json:"total_steps"`
    CreatedAt    int64   `json:"created_at"`
    FinishedAt    int64   `json:"finished_at,omitempty"`
}

// WorkflowStatusResponse represents the response for getting workflow status.
type WorkflowStatusResponse struct {
    ID          string  `json:"id"`
    WorkflowID  string  `json:"workflow_id"`
    Status      string  `json:"status"`
    Inputs      string  `json:"inputs"`
    Outputs     *string `json:"outputs,omitempty"`
    Error       *string `json:"error,omitempty"`
    TotalSteps  int     `json:"total_steps"`
    TotalTokens int     `json:"total_tokens"`
    CreatedAt   string  `json:"created_at"`
    FinishedAt  string  `json:"finished_at"`
    ElapsedTime float64 `json:"elapsed_time"`
}

