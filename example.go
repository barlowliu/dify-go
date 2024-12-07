package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/barlowliu/dify-go.go"
)

func main() {
    // 初始化客户端
    client := dify-go.NewClient("https://dify-go.ai/v1", "your_api_key_here")

    // 设置超时时间
    client.SetTimeout(30 * time.Second)

    // 发送对话消息（阻塞模式）
    chatReq := dify-go.ChatMessageRequest{
        Query:         "What are the specs of the iPhone 13 Pro Max?",
        ResponseMode:  "blocking",
        User:          "abc-123",
        ConversationID: "",
        Files: []dify-go.FileUploadInfo{
            {
                Type:           "image",
                TransferMethod: "remote_url",
                URL:            "https://cloud.dify-go.ai/logo/logo-site.png",
            },
        },
    }

    chatResp, _, err := client.SendChatMessage(context.Background(), chatReq)
    if err != nil {
        log.Fatalf("Error sending chat message: %v", err)
    }

    fmt.Printf("Chat Response: %+v\n", chatResp)

    // 上传文件
    uploadResp, err := client.UploadFile(context.Background(), "path/to/your/image.png", "abc-123")
    if err != nil {
        log.Fatalf("Error uploading file: %v", err)
    }

    fmt.Printf("Uploaded File: %+v\n", uploadResp)

    // 执行工作流（流式模式）
    workflowReq := dify-go.WorkflowRunRequest{
        Inputs:       map[string]interface{}{},
        ResponseMode: "streaming",
        User:         "abc-123",
    }

    _, streamChan, err := client.RunWorkflow(context.Background(), workflowReq)
    if err != nil {
        log.Fatalf("Error running workflow: %v", err)
    }

    // 处理流式响应
    for chunk := range streamChan {
        fmt.Printf("Workflow Chunk: %+v\n", chunk)
    }

    // 停止任务
    stopResp, err := client.StopTask(context.Background(), "task_id_here", "abc-123")
    if err != nil {
        log.Fatalf("Error stopping task: %v", err)
    }

    fmt.Printf("Stop Task Response: %+v\n", stopResp)
}

