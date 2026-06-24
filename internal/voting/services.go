package voting

type Service struct{
	store *RedisStore
}


func NewService(rdb *RedisStore) *Service {
	return &Service{
		store: rdb,
	}
}


func (s *Service) CastVote(id string) {

}


