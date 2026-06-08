// Shared target + level definitions (matches Figma "Level UP" design)

export interface TargetDef {
  value: "boss" | "client" | "friend";
  label: string;
}

export interface LevelDef {
  value: number; // 1-5, sent to backend
  label: string; // playful Thai name
  feeling: string; // "ฟีลลิ่ง" description
  example: string; // example sentence
}

export const TARGETS: TargetDef[] = [
  { value: "boss", label: "หัวหน้า" },
  { value: "client", label: "ลูกค้า" },
  { value: "friend", label: "เพื่อน" },
];

export const LEVELS: LevelDef[] = [
  {
    value: 1,
    label: "สวมวิญญาณผู้ดี",
    feeling:
      "ประชดประชันด้วยความสุภาพขั้นสุด พูดจาขอรับ/คะขา แต่แววตาตัวละครคือนั่งกอดอกมองบน",
    example:
      "“กราบเรียนท่านผู้มีอุปการคุณ โปรดประทานข้อมูลอันเที่ยงตรงให้กระหม่อมด้วยเถิดขอรับ”",
  },
  {
    value: 2,
    label: "พูดดีด้วยละนะ",
    feeling:
      "เหมือนคนพยายามใจเย็นเหน็บแนมเบาๆ ฟีลแอดมินเพจที่ตอบลูกค้ามารอบที่ร้อยของวัน",
    example:
      "“แอดมินอธิบายไว้บรรทัดที่ 3 แล้วนะต๊ะ ลองเพ่งดูอีกที เผื่อเมื่อกี้ลืมลืมตา”",
  },
  {
    value: 3,
    label: "มนุษย์ปกติ",
    feeling:
      "ระดับมาตรฐาน ใช้ง่าย เป็นมิตร คุยง่ายเหมือนเพื่อนทั่วไป ไม่มีคำหยาบหรือคำสแลงเฉพาะกลุ่ม อ่านแล้วเข้าใจทันที ตัวละครลายเส้นดินสอจะยิ้มแย้มแบบปกติที่สุด",
    example:
      "“เรียบร้อยครับ! รบกวนช่วยตรวจสอบความถูกต้องของข้อมูลอีกครั้งนะ”",
  },
  {
    value: 4,
    label: "นึกว่าสนิท",
    feeling:
      "ข้ามขั้นความเกรงใจไปไกล คุยเหมือนสนิทกันมาตั้งแต่ชาติที่แล้ว ชวนทะเลาะและจิกกัดตลอดเวลา",
    example:
      "“เห้ยแก๊! กดเข้ามาทำไมตั้งหลายรอบ เช็กข้อมูลด่วนเลย พิมพ์ผิดมาจะขำให้”",
  },
  {
    value: 5,
    label: "ตัวมัมซัมซุง",
    feeling:
      "โวยวาย พิมพ์ผิดๆ ถูกๆ เน้นอารมณ์ร่วมแบบบ้าคลั่ง ลายเส้นดินสอจะขยุกขยิกสั่นๆ เหมือนคนสติหลุด",
    example:
      "“กรี๊ดดดดดดดดด! เข้าระบบดั้ยล้าววว! เช็กข้อมูลด่วนน ถ้ากดผิดชั้นจะทุบหลังแก!!”",
  },
];

export const targetLabel = (v: string) =>
  TARGETS.find((t) => t.value === v)?.label ?? v;

export const levelLabel = (v: number) =>
  LEVELS.find((l) => l.value === v)?.label ?? "มนุษย์ปกติ";
