package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// chatCompleter is the slice of the OpenAI/DeepSeek client we depend on.
// The real *openai.Client satisfies it; tests inject a fake.
type chatCompleter interface {
	CreateChatCompletion(
		context.Context,
		openai.ChatCompletionRequest,
	) (openai.ChatCompletionResponse, error)
}

type TranslationService struct {
	aiClient chatCompleter
}

func NewTranslationService(apiKey string) (*TranslationService, error) {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.deepseek.com/v1"

	client := openai.NewClientWithConfig(config)
	return &TranslationService{aiClient: client}, nil
}

// newWithClient builds a service around an injected client (used in tests).
func newWithClient(c chatCompleter) *TranslationService {
	return &TranslationService{aiClient: c}
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
// Personality per level, based on the "5 ระดับภาษา" design spec.
var levelTones = map[int]string{
	1: `บุคลิกระดับ 1 — "สวมวิญญาณผู้ดี": ` +
		`ประชดประชันด้วยความสุภาพขั้นสุด ใช้คำราชาศัพท์หรือภาษาสูงเต็มขั้น พูดจาขอรับ/คะขา ` +
		`ฟังดูนอบน้อมสมบูรณ์แบบ แต่แฝงความ condescending เบาๆ — เหมือนผู้ดีที่นั่งกอดอกมองลงมาพร้อมรอยยิ้ม ` +
		`สายตาบอกว่า "โอ้โห... ต้องมาอธิบายถึงเรื่องนี้ด้วยเหรอ" ` +
		`ห้ามประชดตรงๆ — ทุกอย่างต้องฟังดูสุภาพอย่างสมบูรณ์แบบ`,

	2: `บุคลิกระดับ 2 — "พูดดีด้วยละนะ": ` +
		`เขียนเหมือนแอดมินที่ตอบคำถามเดิมมาร้อยครั้งในวันเดียว — ใจเย็น สุภาพ เป็นมืออาชีพ ` +
		`แต่ทุกประโยคแฝงความเหนื่อยล้าและการเหน็บแนมเบาๆ ` +
		`ให้คนอ่านรู้สึกว่า "เขาบอกว่าฉันโง่... แต่อย่างสุภาพมาก" ` +
		`ห้ามประชดตรงๆ แต่ต้องให้สัมผัส "อีกแล้วเหรอ" อยู่ในทุกประโยค`,

	3: `บุคลิกระดับ 3 — "มนุษย์ปกติ (Casual & Clean)": ` +
		`ใช้ภาษาสุภาพทั่วไป เป็นมิตร ตรงไปตรงมา ไม่มีการเสียดสีหรือประชดใดๆ ` +
		`เหมือนคุยกับเพื่อนทั่วไปที่มีมารยาทดี อ่านแล้วเข้าใจทันที ไม่ต้องตีความ`,

	4: `บุคลิกระดับ 4 — "นึกว่าสนิท": ` +
		`ข้ามขั้นความเกรงใจไปโดยสิ้นเชิง คุยเหมือนสนิทกันมาตั้งแต่ชาติที่แล้ว ` +
		`ชวนทะเลาะและจิกกัดตลอด ใช้ภาษาห้าว มีความตลกในแบบเพื่อนสนิทที่แกล้งกัน ` +
		`เจ็บแต่รู้ว่าไม่ได้จริงจัง สาระสำคัญยังครบ`,

	5: `บุคลิกระดับ 5 — "ตัวมัมซัมซุง": ` +
		`โวยวายเต็มที่ ใช้อักษรซ้ำเกินจำเป็น (เช่น "กรี๊ดดดดดด" "ล้าวววว" "ดั้ยยยย") ` +
		`พิมพ์ผิดๆ ถูกๆ ตัวพิมพ์ใหญ่เล็กสลับกันได้ เน้นอารมณ์ร่วมแบบบ้าคลั่ง ` +
		`เหมือนคนสติหลุดกำลังพิมพ์อยู่บน keyboard ที่สั่น ` +
		`ห้ามสุภาพ ห้ามเรียบร้อย แต่สาระสำคัญยังต้องครบ`,
}

// PurifyText translates angry text using target context + level tone.
// target: "boss" | "client" | "teacher" | "friend"
// level:  1-5
// lang:   "th" (Thai output) | "en" (English output)
func (s *TranslationService) PurifyText(
	ctx context.Context,
	inputRawText string,
	target string,
	level int,
	lang string,
) (string, error) {
	// Defaults
	if level < 1 || level > 5 {
		level = 3
	}
	if target == "" {
		target = "boss"
	}
	if lang == "" {
		lang = "th"
	}

	ctxDesc, ok := targetContext[target]
	if !ok {
		ctxDesc = targetContext["boss"]
	}
	toneDesc, ok := levelTones[level]
	if !ok {
		return "", fmt.Errorf("unsupported level: %d", level)
	}

	langInstr := "ตอบเป็นภาษาไทยเท่านั้น"
	if lang == "en" {
		langInstr = "Respond in English only. Apply the same personality/tone described above, adapted naturally to English."
	}

	systemPrompt := fmt.Sprintf(
		`คุณคือผู้เชี่ยวชาญด้านการสื่อสารที่เชี่ยวชาญบุคลิกภาษาหลากหลายรูปแบบ

%s
%s

หน้าที่ของคุณ: รับข้อความที่มีอารมณ์ คำหยาบ หรือความโกรธ แล้วแปลงเป็นข้อความที่เหมาะสมกับผู้รับและบุคลิกที่กำหนด
โดยยังคงสาระสำคัญ (Core Message) ไว้ครบถ้วน

%s

กฎเหล็ก: ห้ามใส่คำอธิบาย คำนำ หรือบทสรุปใดๆ — ตอบเฉพาะข้อความที่แปลเสร็จแล้วเท่านั้น`,
		ctxDesc, toneDesc, langInstr,
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
