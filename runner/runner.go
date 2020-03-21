package runner

// Runner represents runner docker information
type Runner struct {
	ID      string `json:"id" bson:"id"`
	Port    string `json:"port" bson:"port"`
	RedisID string `json:"rid" bson:"redisid"`
}
