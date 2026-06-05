package service

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type TranslationService struct {
	aiClient *openai.Client
}

func NewTranslationService(apiKey string) (*TranslationService, error) {
	// ⭐️ ทริคเด็ด: ใช้โครงสร้าง OpenAI แต่ออกเดินทางไปหา DeepSeek เซิร์ฟเวอร์
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.deepseek.com/v1" 

	client := openai.NewClientWithConfig(config)
	return &TranslationService{aiClient: client}, nil
}

func (s *TranslationService) PurifyText(ctx context.Context, inputRawText string) (string, error) {
	systemPrompt := `คุณคือผู้เชี่ยวชาญด้านการสื่อสารในองค์กร หน้าที่ของคุณคือการรับข้อความที่มีคำด่า คำหยาบคาย หรืออารมณ์ฉุนเฉียว 
แล้วนำมาเรียบเรียงใหม่ให้กลายเป็นภาษาเขียนระดับมืออาชีพ สุภาพ อ่อนน้อม แต่ยังคงรักษาเนื้อหาหรือสาระสำคัญ (Core Message) ที่ผู้พูดต้องการสื่อไว้ครบถ้วน 
เหมาะสำหรับนำไปส่งต่อให้หัวหน้าหรือลูกค้าอ่านได้จริงในแชททำงาน โดยห้ามใส่คำอธิบายเพิ่มเติมใดๆ ให้ตอบเฉพาะข้อความที่ปรับปรุงเสเคราะห์เสร็จแล้วเท่านั้น`

	// เรียกใช้งานโมเดลราคาถูกและเก่งของ DeepSeek
	resp, err := s.aiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "deepseek-chat", // ใช้โมเดลแชทมาตรฐานของเขาได้เลย
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: inputRawText,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}