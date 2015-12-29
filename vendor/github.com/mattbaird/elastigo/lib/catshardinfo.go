package elastigo

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type CatShards []CatShardInfo

// Stringify the shards
func (s *CatShards) String() string {
	var buffer bytes.Buffer

	if s != nil {
		for _, cs := range *s {
			buffer.WriteString(fmt.Sprintf("%v\n", cs))
		}
	}
	return buffer.String()
}

var ErrInvalidShardLine = errors.New("Cannot parse shardline")

// Create a CatShard from a line of the raw output of a _cat/shards
func NewCatShardInfo(rawCat string) (catshard *CatShardInfo, err error) {

	split := strings.Fields(rawCat)
	if len(split) < 4 {
		return nil, ErrInvalidShardLine
	}
	catshard = &CatShardInfo{}
	catshard.IndexName = split[0]
	catshard.Shard, err = strconv.Atoi(split[1])
	if err != nil {
		catshard.Shard = -1
	}
	catshard.Primary = split[2]
	catshard.State = split[3]
	if len(split) == 4 {
		return catshard, nil
	}

	catshard.Docs, err = strconv.ParseInt(split[4], 10, 64)
	if err != nil {
		catshard.Docs = 0
	}
	if len(split) == 5 {
		return catshard, nil
	}
	catshard.Store, err = strconv.ParseInt(split[5], 10, 64)
	if err != nil {
		catshard.Store = 0
	}
	if len(split) == 6 {
		return catshard, nil
	}
	catshard.NodeIP = split[6]
	if len(split) == 7 {
		return catshard, nil
	}
	catshard.NodeName = split[7]
	if len(split) > 8 {
	loop:
		for i, moreName := range split {
			if i > 7 {
				if moreName == "->" {
					break loop
				}
				catshard.NodeName += " "
				catshard.NodeName += moreName
			}
		}
	}

	return catshard, nil
}

// Print shard info
func (s *CatShardInfo) String() string {
	if s == nil {
		return ":::::::"
	}
	return fmt.Sprintf("%v:%v:%v:%v:%v:%v:%v:%v", s.IndexName, s.Shard, s.Primary,
		s.State, s.Docs, s.Store, s.NodeIP, s.NodeName)
}

// Get all the shards, even the bad ones
func (c *Conn) GetCatShards() (shards CatShards) {
	shards = make(CatShards, 0)
	args := map[string]interface{}{"bytes": "b"}
	s, err := c.DoCommand("GET", "/_cat/shards", args, nil)
	if err == nil {
		catShardLines := strings.Split(string(s[:]), "\n")
		for _, shardLine := range catShardLines {
			shard, _ := NewCatShardInfo(shardLine)
			if nil != shard {
				shards = append(shards, *shard)
			}
		}
	}
	return shards
}
