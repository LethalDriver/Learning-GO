type RoomRepository interface {
	FindAll() ([]*Room, error)
	FindById(id int) (*Room, error)
	Create(room *Room) error
}

type RedisRoomRepository struct {
	redisClient *redis.Client
}

func NewRedisRoomRepository(redisClient *redis.Client) *RedisRoomRepository {
	return &RedisRoomRepository{redisClient: redisClient}
}

func (r *RedisRoomRepository) FindAll() ([]*Room, error) {
	var rooms []*Room
	keys, err := r.redisClient.Keys("room:*").Result()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		id, _ := strconv.Atoi(strings.TrimPrefix(key, "room:"))
		room, err := r.FindById(id)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (r *RedisRoomRepository) FindById(id int) (*Room, error) {
	room := &Room{}
	err := r.redisClient.Get(fmt.Sprintf("room:%d", id)).Scan(room)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *RedisRoomRepository) Create(room *Room) error {
	return r.redisClient.Set(fmt.Sprintf("room:%d", room.Id), room, 0).Err()
}
