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

// ── Context descriptions ──────────────────────────────────────────────────────
// Who the message is going TO — shapes relationship dynamic in the prompt.
var targetContext = map[string]string{
	"boss": `ผู้รับข้อความคือ "เจ้านาย" ซึ่งมีอำนาจเหนือกว่าและสามารถกระทบต่ออนาคตการทำงานของเราได้ ` +
		`ต้องรักษาความสัมพันธ์ไว้ ฟังดูเคารพ แต่ยังสื่อสาระสำคัญได้ครบ`,

	"client": `ผู้รับข้อความคือ "ลูกค้า" ซึ่งเป็นแหล่งรายได้ของบริษัท ` +
		`ต้องรักษาภาพลักษณ์มืออาชีพ ใช้น้ำเสียง Customer Service `,

	"teacher": `ผู้รับข้อความคือ "อาจารย์" ซึ่งเป็นผู้ให้ความรู้และมีอำนาจในการให้เกรด ` +
		`ต้องแสดงความนับถือในฐานะนักเรียน ใช้ภาษาสุภาพ เป็นทางการ `,

	"friend": `ผู้รับข้อความคือ "เพื่อน" ซึ่งมีความสัมพันธ์เท่าเทียมกัน ` +
		`สามารถใช้ภาษาที่ไม่เป็นทางการได้ แต่ยังต้องสื่อสาระสำคัญให้ชัดเจน `,
}

// ── Level tone modifiers ──────────────────────────────────────────────────────
// How polite/aggressive to be — layered on top of the target context.
var levelTones = map[int]string{
	1: `ระดับความสุภาพ: สูงสุด — ใช้คำราชาศัพท์หรือภาษาทางการสูงสุดเท่าที่เหมาะสม ` +
		`อ่อนน้อมถ่อมตนสุดๆ ไม่มีการเสียดสีหรือประชดใดๆ ทั้งสิ้น`,

	2: `ระดับความสุภาพ: สุภาพ — ใช้ภาษาธุรกิจระดับมืออาชีพ เป็นทางการ ไม่มีการเสียดสี`,

	3: `ระดับความสุภาพ: ปกติ — ใช้ภาษาสุภาพทั่วไปที่ใช้ในการทำงาน ตรงไปตรงมา ไม่มีการเสียดสี`,

	4: `ระดับความสุภาพ: แรงนิดๆ — ใช้ภาษาสุภาพแต่แอบประชดอยู่นิดหน่อย ` +
		`คนอ่านรู้สึกว่า "นี่มันฟังดูดี... แต่ทำไมรู้สึกโดนด่าน้อยๆ" ` +
		`ห้ามด่าตรงๆ แต่ต้องให้คนอ่านรู้สึกอึดอัดเล็กน้อย`,

	5: `ระดับความสุภาพ: แรงสุด — Passive-Aggressive ขั้นเทพ ` +
		`ทุกประโยคต้องแทงใจ แต่ปฏิเสธว่าเสียดสีไม่ได้ ` +
		`เทคนิคที่ต้องใช้: "ขอบคุณที่...(ตามด้วยสิ่งน่าหงุดหงิด)", ` +
		`"ด้วยความเข้าใจอย่างยิ่ง(ทั้งที่ไม่เข้าใจเลย)", ` +
		`"เชื่อว่าทุกคนมีเหตุผลของตัวเอง" ` +
		`คนอ่านต้องรู้สึกเจ็บแต่จับผิดไม่ได้`,
}

// PurifyText translates angry text using target context + level tone.
// target: "boss" | "client" | "teacher" | "friend"
// level:  1-5
func (s *TranslationService) PurifyText(
	ctx context.Context,
	inputRawText string,
	target string,
	level int,
) (string, error) {
	// Defaults
	if level < 1 || level > 5 {
		level = 3
	}
	if target == "" {
		target = "boss"
	}

	ctxDesc, ok := targetContext[target]
	if !ok {
		ctxDesc = targetContext["boss"]
	}
	toneDesc, ok := levelTones[level]
	if !ok {
		return "", fmt.Errorf("unsupported level: %d", level)
	}

	systemPrompt := fmt.Sprintf(
		`คุณคือผู้เชี่ยวชาญด้านการสื่อสารและนักเขียน Passive-Aggressive ระดับมืออาชีพ

%s
%s

หน้าที่ของคุณ: รับข้อความที่มีอารมณ์ คำหยาบ หรือความโกรธ แล้วแปลงเป็นข้อความที่เหมาะสมกับผู้รับและระดับที่กำหนด
โดยยังคงสาระสำคัญ (Core Message) ไว้ครบถ้วน

กฎเหล็ก: ห้ามใส่คำอธิบาย คำนำ หรือบทสรุปใดๆ — ตอบเฉพาะข้อความที่แปลเสร็จแล้วเท่านั้น`,
		ctxDesc, toneDesc,
	)

	resp, err := s.aiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "deepseek-chat",
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
				{Role: openai.ChatMessageRoleUser, Content: inputRawText},
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
