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
    
    // Store with project ID as the field name
    return s.rdb.HSet(ctx, key, id, val).Err()
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
	key := fmt.Sprintf("hasVote:%s", user_id)
	ctx := context.Background()
	projectsKey := fmt.Sprintf("projects:%s", project_id)
	res := s.rdb.SetArgs(ctx, key, "voted", redis.SetArgs{
		Mode: "NX",
	})
	ok := res.Val() == "Ok"
	if !ok {
		return UserAlreadyVoted
	}
	_, err := s.rdb.HIncrBy(ctx, projectsKey, "votes", 1).Result()
 return err
}
