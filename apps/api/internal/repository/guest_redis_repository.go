package repository

import (
	"context"
	"corporate-translator-api/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"github.com/redis/go-redis/v9"
)

const (
	guestKeyPrefix     = "guest:"
	guestInitialCredit = 3.00
	guestTTL           = 24 * time.Hour
)

var (
	 ErrInsufficientCredit = errors.New("insufficient credit")
    ErrGuestNotFound      = errors.New("guest not found")
)

type guestRedisRepository struct {
	client *redis.Client
}

func NewGuestRepository(client *redis.Client) GuestRepository {
	return &guestRedisRepository{client: client}
}

// ---- Helper ---------
func (r *guestRedisRepository) key(guestID string) string {
	return guestKeyPrefix + guestID
}

func (r *guestRedisRepository) save(ctx context.Context, key string, g *model.Guest) error {
	data, err := json.Marshal(g)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, guestTTL).Err()
}

func (r *guestRedisRepository) GetOrCreate(ctx context.Context, guestID string) (*model.Guest, error ) {
	key := r.key(guestID)

	val, err := r.client.Get(ctx, key).Result()

	if err == redis.Nil {
		guest := &model.Guest{
			Credit: 3.00,
			CreatedAt: time.Now().UTC(),
		}
		if err := r.save(ctx,key, guest); err != nil {
			return nil, fmt.Errorf("GetorCreate: %w", err)
		}
		return guest, nil
	}

	if err != nil {
        return nil, fmt.Errorf("GetOrCreate: %w", err)
	}

	var guest model.Guest

	if err := json.Unmarshal([]byte(val), &guest); err != nil {
	return nil, fmt.Errorf("GetOrCreate unmarshal: %w", err)
    }
    return &guest, nil

}

// DeductCredit — ตัด credit แบบ Atomic ด้วย Lua Script
// Atomic = ทำทีเดียวพร้อมกัน ไม่มีใครแทรกกลางได้
// ถ้าไม่ใช้ Lua: GET → เช็ค → SET (3 ขั้น มีช่องว่าง)
// → request 2 ตัวมาพร้อมกัน GET ได้ credit=1 ทั้งคู่ → ติดลบได้!
func (r *guestRedisRepository) DeductCredit(ctx context.Context, guestID string, amount float64) error {
    key := r.key(guestID)

    luaScript := redis.NewScript(`
        local val = redis.call("GET", KEYS[1])
        if not val then
            return redis.error_reply("GUEST_NOT_FOUND")
        end
        local data = cjson.decode(val)
        if data.credit < tonumber(ARGV[1]) then
            return redis.error_reply("INSUFFICIENT_CREDIT")
        end
        data.credit = data.credit - tonumber(ARGV[1])
        redis.call("SET", KEYS[1], cjson.encode(data), "EX", ARGV[2])
        return redis.status_reply("OK")
    `)
    // KEYS[1] = key ของ guest เช่น "guest:abc-123"
    // ARGV[1] = จำนวน credit ที่ตัด เช่น 1.0
    // ARGV[2] = TTL วินาที — reset นาฬิกาทุกครั้งที่ใช้งาน

    err := luaScript.Run(ctx, r.client,
        []string{key},
        amount,
        int(guestTTL.Seconds()),
    ).Err()

    if err != nil {
        if err.Error() == "INSUFFICIENT_CREDIT" {
            return ErrInsufficientCredit
        }
        if err.Error() == "GUEST_NOT_FOUND" {
            return ErrGuestNotFound
        }
        return fmt.Errorf("DeductCredit: %w", err)
    }
    return nil
}

// RefundCredit — คืน credit กลับ (ใช้ตอน AI fail)
func (r *guestRedisRepository) RefundCredit(ctx context.Context, guestID string, amount float64) error {
    key := r.key(guestID)

    luaScript := redis.NewScript(`
        local val = redis.call("GET", KEYS[1])
        if not val then
            return redis.status_reply("OK")
        end
        local data = cjson.decode(val)
        data.credit = data.credit + tonumber(ARGV[1])
        redis.call("SET", KEYS[1], cjson.encode(data), "EX", ARGV[2])
        return redis.status_reply("OK")
    `)

    return luaScript.Run(ctx, r.client,
        []string{key},
        amount,
        int(guestTTL.Seconds()),
    ).Err()
}

// Delete — ลบ guest session (ใช้ตอน login)
func (r *guestRedisRepository) Delete(ctx context.Context, guestID string) error {
    return r.client.Del(ctx, r.key(guestID)).Err()
}