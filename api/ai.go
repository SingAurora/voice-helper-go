package api

import (
	"context"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/utils"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"io"
)

/**
 * Authentication
 * 1.If you authorize your endpoint using an API key, you can set your api key to environment variable "ARK_API_KEY"
 * client := arkruntime.NewClientWithApiKey(os.Getenv("ARK_API_KEY"))
 * Note: If you use an API key, this API key will not be refreshed.
 * To prevent the API from expiring and failing after some time, choose an API key with no expiration date.
 *
 * 2.If you authorize your endpoint with Volcengine Identity and Access Management（IAM), set your api key to environment variable "VOLC_ACCESSKEY", "VOLC_SECRETKEY"
 * client := arkruntime.NewClientWithAkSk(os.Getenv("VOLC_ACCESSKEY"), os.Getenv("VOLC_SECRETKEY"))
 * To get your ak&sk, please refer to this document(https://www.volcengine.com/docs/6291/65568)
 * For more information，please check this document（https://www.volcengine.com/docs/82379/1263279）
 */

func Ai(content string) (*utils.ChatCompletionStreamReader, error) {
	client := arkruntime.NewClientWithApiKey("a90969af-9170-4e7a-913a-631374b256a3")
	ctx := context.Background()

	req := model.ChatCompletionRequest{
		Model: "ep-20240809104616-xz2tz",
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("你是豆包，是由字节跳动开发的 AI 人工智能助手"),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(content),
				},
			},
		},
	}
	return client.CreateChatCompletionStream(ctx, req)
	//if err != nil {
	//	fmt.Printf("stream chat error: %v\n", err)
	//	return
	//}
	//defer stream.Close()
	//
	//for {
	//	recv, err := stream.Recv()
	//	if err == io.EOF {
	//		return
	//	}
	//	if err != nil {
	//		fmt.Printf("Stream chat error: %v\n", err)
	//		return
	//	}
	//
	//	if len(recv.Choices) > 0 {
	//		fmt.Print(recv.Choices[0].Delta.Content)
	//	}
	//}
}

func AiTest() {
	client := arkruntime.NewClientWithApiKey("a90969af-9170-4e7a-913a-631374b256a3")
	ctx := context.Background()

	req := model.ChatCompletionRequest{
		Model: "ep-20240809104616-xz2tz",
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("你是豆包，是由字节跳动开发的 AI 人工智能助手"),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("常见的十字花科植物有哪些？"),
				},
			},
		},
	}
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("stream chat error: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Printf("Stream chat error: %v\n", err)
			return
		}

		if len(recv.Choices) > 0 {
			fmt.Print(recv.Choices[0].Delta.Content)
		}
	}
}
