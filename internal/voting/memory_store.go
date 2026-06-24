package voting



type MemoryStore struct {
	voters []User
	projects []Project

}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		voters: []User{},
		projects: []Project{},
	}
}

func (s *MemoryStore) ListProjects() []Project {
	//return the list of projects
	return s.projects
}
func (s *MemoryStore) GetUser(id string) (User, error) {
	for _, val := range s.voters{
		if (val.Id == id){
			return val, nil
		}
	}
	return User{}, UserNotFound
}
