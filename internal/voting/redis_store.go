package voting

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)


type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb * redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func (s *RedisStore) AddProjects(p *Project) error {
    id := uuid.New().String()
    p.Id = id
    ctx := context.Background()
    key := "projects"  // Hash key
    
    val, err := json.Marshal(p)
    if err != nil {
        return err
    }
	if err := s.rdb.HSet(ctx, key, id, val).Err(); err != nil {
        return err
    }

// Initialize votes counter for this project
    votesKey := fmt.Sprintf("projects:%s:votes", id)
    return s.rdb.Set(ctx, votesKey, 0, 0).Err() // 0 means no expiration
 
}

func (s *RedisStore) ListProjects() ([]Project, error) {
    ctx := context.Background()
    key := "projects"
    
    // Get all fields and values
    val, err := s.rdb.HGetAll(ctx, key).Result()
    if err != nil {
        return nil, err
    }
    
    res := make([]Project, 0, len(val))
    for _, v := range val {
        var prj Project
        if err := json.Unmarshal([]byte(v), &prj); err != nil {
            // Better to return error instead of panic
            return nil, fmt.Errorf("failed to unmarshal project: %w", err)
        }
        res = append(res, prj)
    }
    return res, nil
}

func (s *RedisStore) CastVote(user_id string, project_id string) error {
    key := fmt.Sprintf("hasVoted:%s", user_id)
    ctx := context.Background()
    
    // Check if user already voted
    res := s.rdb.SetArgs(ctx, key, "voted", redis.SetArgs{Mode: "NX"})
    if res.Val() != "Ok" {
        return UserAlreadyVoted
    }
    
    // Increment vote counter for the project
    votesKey := fmt.Sprintf("projects:%s:votes", project_id)
    _, err := s.rdb.Incr(ctx, votesKey).Result()
    return err
}
