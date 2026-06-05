package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type TranslationService struct {
	aiClient *openai.Client
}

func NewTranslationService(apiKey string) (*TranslationService, error) {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.deepseek.com/v1"

	client := openai.NewClientWithConfig(config)
	return &TranslationService{aiClient: client}, nil
}

// Level-based system prompts — from ultra-polite to maximum passive-aggression.
var levelPrompts = map[int]string{
	1: `คุณคือผู้เชี่ยวชาญภาษาระดับราชการ หน้าที่ของคุณคือรับข้อความที่มีอารมณ์ ` +
		`แล้วเปลี่ยนให้เป็นภาษาระดับราชาศัพท์ สุภาพที่สุดเท่าที่จะเป็นไปได้ ` +
		`เหมือนเขียนจดหมายถึงท่านรัฐมนตรี อ่อนน้อมถ่อมตน ใช้คำลงท้ายอย่าง "ครับ/ค่ะ" ` +
		`หรือ "ด้วยความเคารพอย่างสูง" ` +
		`ห้ามใส่คำอธิบายใดๆ ตอบเฉพาะข้อความที่แปลเสร็จแล้วเท่านั้น`,

	2: `คุณคือผู้เชี่ยวชาญด้านการสื่อสารในองค์กร หน้าที่ของคุณคือรับข้อความที่มีอารมณ์ ` +
		`แล้วเรียบเรียงใหม่ให้เป็นภาษาธุรกิจระดับมืออาชีพ สุภาพ เป็นทางการ ` +
		`เหมาะสำหรับส่ง Email หาหัวหน้าหรือลูกค้า ` +
		`ยังคงสาระสำคัญไว้ครบถ้วน ` +
		`ห้ามใส่คำอธิบายใดๆ ตอบเฉพาะข้อความที่แปลเสร็จแล้วเท่านั้น`,

	3: `คุณคือผู้เชี่ยวชาญด้านการสื่อสารในองค์กร หน้าที่ของคุณคือรับข้อความที่มีคำด่า คำหยาบคาย หรืออารมณ์ฉุนเฉียว ` +
		`แล้วนำมาเรียบเรียงใหม่ให้กลายเป็นภาษาเขียนระดับมืออาชีพ สุภาพ อ่อนน้อม แต่ยังคงรักษาเนื้อหาหรือสาระสำคัญ (Core Message) ที่ผู้พูดต้องการสื่อไว้ครบถ้วน ` +
		`เหมาะสำหรับนำไปส่งต่อให้หัวหน้าหรือลูกค้าอ่านได้จริงในแชททำงาน ` +
		`ห้ามใส่คำอธิบายใดๆ ตอบเฉพาะข้อความที่แปลเสร็จแล้วเท่านั้น`,

	4: `คุณคือนักเขียนที่เก่งเรื่อง Passive-Aggressive หน้าที่ของคุณคือรับข้อความที่มีอารมณ์ ` +
		`แล้วเรียบเรียงให้ฟังดูสุภาพ... แต่จริงๆ แล้วมีมุกแอบแดก มีประชดนิดๆ ` +
		`คนอ่านต้องรู้สึกว่า "เฮ้ย นี่มันดูดีนะ... แต่ทำไมรู้สึกโดนด่า" ` +
		`เน้นเสียดสีอย่างมีสไตล์ ปฏิเสธไม่ได้ว่าไม่สุภาพ แต่คนเขียนจะรู้สึกสะใจ ` +
		`ห้ามใส่คำอธิบายใดๆ ตอบเฉพาะข้อความที่แปลเสร็จแล้วเท่านั้น`,

	5: `คุณคือราชาแห่งการแดก ระดับ Passive-Aggressive ขั้นเทพ ` +
		`หน้าที่ของคุณคือรับข้อความที่มีอารมณ์โกรธ แล้วแปลงมันเป็นข้อความที่ "ดูสุภาพ" ` +
		`แต่ทุกประโยคเต็มไปด้วยการประชดประชัน เหน็บแนม แรงสุดขีด ` +
		`คนอ่านต้องรู้สึกเจ็บแต่จับผิดไม่ได้ เหมือนโดนตบด้วยถุงมือกำมะหยี่ ` +
		`ใช้เทคนิค: ขอบคุณที่... (แล้วตามด้วยสิ่งที่น่าหงุดหงิด), "ด้วยความเข้าใจ" (ทั้งที่ไม่เข้าใจเลย) ` +
		`ห้ามใส่คำอธิบายใดๆ ตอบเฉพาะข้อความที่แปลเสร็จแล้วเท่านั้น`,
}

// PurifyText translates angry text according to the selected level (1-5).
func (s *TranslationService) PurifyText(ctx context.Context, inputRawText string, level int) (string, error) {
	if level < 1 || level > 5 {
		level = 3 // default to casual professional
	}

	systemPrompt, ok := levelPrompts[level]
	if !ok {
		return "", fmt.Errorf("unsupported level: %d", level)
	}

	resp, err := s.aiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "deepseek-chat",
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

	if len(resp.Choices) == 0 {
		return "", errors.New("no response from translation model")
	}

	return resp.Choices[0].Message.Content, nil
}
