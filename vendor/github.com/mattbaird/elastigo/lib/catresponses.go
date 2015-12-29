package elastigo

type CatIndexInfo struct {
	Health   string
	Name     string
	Shards   int
	Replicas int
	Docs     CatIndexDocs
	Store    CatIndexStore
}

type CatIndexDocs struct {
	Count   int64
	Deleted int64
}

type CatIndexStore struct {
	Size    int64
	PriSize int64
}

type CatShardInfo struct {
	IndexName string
	Shard     int
	Primary   string
	State     string
	Docs      int64
	Store     int64
	NodeIP    string
	NodeName  string
}
