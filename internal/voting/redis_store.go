package voting

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisStore handles Redis operations for the voting system
type RedisStore struct {
	rdb *redis.Client
}

// NewRedisStore creates a new RedisStore instance
func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

// AddProjects:: Implemented by Suresh
func (s *RedisStore) AddProjects(p *Project) error {
	ctx := context.Background()
	prJson, err := json.Marshal(p)
	if err != nil {
		return err
	}
	projectKey := "project"
	rerr := s.rdb.HSet(ctx, projectKey, p.Id, prJson).Err()
	if rerr != nil {
		return rerr
	}
	//setup voting field for the project entry
	voteKey := "votes"
	voteErr := s.rdb.HSet(ctx, voteKey, p.Id, 0).Err()
	if voteErr != nil {
		return voteErr
	}
	return nil
	}

// ListProjects retrieves all projects from Redis
func (s *RedisStore) ListProjects() ([]Project, error) {
	ctx := context.Background()
	key := "project"

	// Check if projects key exists
	exists, err := s.rdb.Exists(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to check projects existence: %w", err)
	}
	if exists == 0 {
		return []Project{}, nil // Return empty slice if no projects
	}

	// Get all fields and values
	val, err := s.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch projects: %w", err)
	}

	projects := make([]Project, 0, len(val))
	for id, v := range val {
		var prj Project
		if err := json.Unmarshal([]byte(v), &prj); err != nil {
			return nil, fmt.Errorf("failed to unmarshal project (ID: %s): %w", id, err)
		}
		projects = append(projects, prj)
	}
	return projects, nil
}



// CastVote allows a user to vote for a project
func (s *RedisStore) CastVote(userID string, projectID string) error {
	ctx := context.Background()
	// Check for user
	userKey := "users"
	exists, existsErr := s.rdb.HExists(ctx, userKey, userID). Result()
	if existsErr != nil {
		return existsErr
	}
	if !exists {
		return UserNotFound
	}

	//check hasVoted + set vote
		votedKey := fmt.Sprintf("voted:%s", userID)
	_, err:= s.rdb.SetArgs(ctx, votedKey, true, redis.SetArgs{
		Mode: "NX",
	}).Result()
	if err != nil {
		return UserAlreadyVoted
	}

	//increase vote by 1
	voteKey := fmt.Sprintf("votes:%s", projectID)
	_, incrErr :=s.rdb.Incr(ctx, voteKey).Result()
	if incrErr != nil {
		deleteKey := fmt.Sprintf("hasVoted:%s", userID)
		s.rdb.Del(ctx, deleteKey)
		return errors.New("failed to cast vote")
	} 

	return nil
}



// AddUser creates a new user with email uniqueness check
func (s *RedisStore) AddUser(user *User) (*User, error) {
	ctx := context.Background()
	usersKey := "users"
	emailKey := "users:email" // Separate hash for email lookups

	// Check if email already exists
	exists, err := s.rdb.HExists(ctx, emailKey, user.Email).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email %s already exists", user.Email)
	}

	// Marshal user to JSON
	val, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}
	// Add User
	hsetErr := s.rdb.HSet(ctx, usersKey, user.Id, val).Err()
	if hsetErr != nil {
		return nil , hsetErr
	}
	//Set email key
	emailSetErr := s.rdb.HSet(ctx, emailKey, user.Email, user.Id).Err()
	if emailSetErr != nil {
		return nil, emailSetErr
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *RedisStore) GetUser(userID string) (*User, error) {	

	ctx := context.Background()
	usersKey := "users"

	val, err := s.rdb.HGet(ctx, usersKey, userID).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("user %s not found", userID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var user User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *RedisStore) GetUserByEmail(email string) (*User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	ctx := context.Background()
	emailKey := "users:email"

	// Get user ID from email mapping
	userID, err := s.rdb.HGet(ctx, emailKey, email).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("user with email %s not found", email)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID by email: %w", err)
	}

	// Get user by ID
	return s.GetUser(userID)
}

func (s *RedisStore) CreateAdmin(email string, password_hash string) error {
	ctx := context.Background()
	adminKey := fmt.Sprintf("admin:%s", email)
	_, err:= s.rdb.SetArgs(ctx, adminKey, password_hash, redis.SetArgs{
		Mode: "NX",
	}).Result()
	if err != nil {
		return AdminAlreadyExists
	}

	return nil
}

func (s *RedisStore) GetAdminPassword(email string) (string, error) {
    ctx := context.Background()
    adminKey := fmt.Sprintf("admin:%s", email)
    
    // Get the password hash (stored as string)
    passwordHash, err := s.rdb.Get(ctx, adminKey).Result()
    if err == redis.Nil {
        return "", AdminNotFound // User doesn't exist
    } else if err != nil {
        return "", err // Redis error
    }
    
    return passwordHash, nil
}






