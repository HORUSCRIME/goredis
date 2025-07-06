package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/HORUSCRIME/goredis/database"
	"github.com/HORUSCRIME/goredis/resp"
)

func PingCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.NewSimpleString("PONG")
	}
	if len(args) == 1 {
		return resp.NewBulkString(args[0].Bulk)
	}
	return resp.NewError("ERR wrong number of arguments for 'ping' command")
}

func EchoCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'echo' command")
	}
	return resp.NewBulkString(args[0].Bulk)
}

func SetCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.NewError("ERR wrong number of arguments for 'set' command")
	}
	key := string(args[0].Bulk)
	value := string(args[1].Bulk)
	ttl := time.Duration(0)

	if len(args) > 2 {
		for i := 2; i < len(args); i++ {
			option := strings.ToUpper(string(args[i].Bulk))
			switch option {
			case "EX":
				if i+1 >= len(args) {
					return resp.NewError("ERR syntax error")
				}
				seconds, err := strconv.ParseInt(string(args[i+1].Bulk), 10, 64)
				if err != nil || seconds <= 0 {
					return resp.NewError("ERR invalid expire time in 'SET' command")
				}
				ttl = time.Duration(seconds) * time.Second
				i++
			default:
				return resp.NewError(fmt.Sprintf("ERR unknown option '%s' for 'SET' command", option))
			}
		}
	}

	db.Set(key, database.NewString(value), ttl)
	return resp.NewSimpleString("OK")
}

func GetCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'get' command")
	}
	key := string(args[0].Bulk)

	val, ok := db.Get(key)
	if !ok {
		return resp.NewNullBulkString()
	}

	strVal, isString := val.(*database.String)
	if !isString {
		return resp.NewError(fmt.Sprintf("WRONGTYPE Operation against a key holding the wrong kind of value"))
	}
	return resp.NewBulkString([]byte(strVal.Val))
}

func DelCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.NewError("ERR wrong number of arguments for 'del' command")
	}
	deletedCount := 0
	for _, arg := range args {
		key := string(arg.Bulk)
		if db.Delete(key) {
			deletedCount++
		}
	}
	return resp.NewInteger(int64(deletedCount))
}

func ExistsCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.NewError("ERR wrong number of arguments for 'exists' command")
	}
	foundCount := 0
	for _, arg := range args {
		key := string(arg.Bulk)
		if db.Exists(key) {
			foundCount++
		}
	}
	return resp.NewInteger(int64(foundCount))
}

func TypeCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'type' command")
	}
	key := string(args[0].Bulk)
	return resp.NewSimpleString(db.Type(key))
}

func LPushCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.NewError("ERR wrong number of arguments for 'lpush' command")
	}
	key := string(args[0].Bulk)
	elements := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		elements[i] = string(arg.Bulk)
	}

	val, ok := db.Get(key)
	var list *database.List
	if !ok {
		list = database.NewList()
		db.Set(key, list, 0)
	} else {
		existingList, isList := val.(*database.List)
		if !isList {
			return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		list = existingList
	}
	newLen := list.LPush(elements...)
	return resp.NewInteger(int64(newLen))
}

func RPushCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.NewError("ERR wrong number of arguments for 'rpush' command")
	}
	key := string(args[0].Bulk)
	elements := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		elements[i] = string(arg.Bulk)
	}

	val, ok := db.Get(key)
	var list *database.List
	if !ok {
		list = database.NewList()
		db.Set(key, list, 0)
	} else {
		existingList, isList := val.(*database.List)
		if !isList {
			return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		list = existingList
	}
	newLen := list.RPush(elements...)
	return resp.NewInteger(int64(newLen))
}

func LPopCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'lpop' command")
	}
	key := string(args[0].Bulk)

	val, ok := db.Get(key)
	if !ok {
		return resp.NewNullBulkString()
	}
	list, isList := val.(*database.List)
	if !isList {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	element, success := list.LPop()
	if !success {
		return resp.NewNullBulkString()
	}
	return resp.NewBulkString([]byte(element))
}

func RPopCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'rpop' command")
	}
	key := string(args[0].Bulk)

	val, ok := db.Get(key)
	if !ok {
		return resp.NewNullBulkString()
	}
	list, isList := val.(*database.List)
	if !isList {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	element, success := list.RPop()
	if !success {
		return resp.NewNullBulkString()
	}
	return resp.NewBulkString([]byte(element))
}

func LLenCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'llen' command")
	}
	key := string(args[0].Bulk)

	val, ok := db.Get(key)
	if !ok {
		return resp.NewInteger(0)
	}
	list, isList := val.(*database.List)
	if !isList {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return resp.NewInteger(int64(list.LLen()))
}

func HSetCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) < 3 || len(args)%2 != 1 {
		return resp.NewError("ERR wrong number of arguments for 'hset' command")
	}
	key := string(args[0].Bulk)

	val, ok := db.Get(key)
	var hash *database.Hash
	if !ok {
		hash = database.NewHash()
		db.Set(key, hash, 0)
	} else {
		existingHash, isHash := val.(*database.Hash)
		if !isHash {
			return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		hash = existingHash
	}

	addedOrUpdatedCount := 0
	for i := 1; i < len(args); i += 2 {
		field := string(args[i].Bulk)
		value := string(args[i+1].Bulk)
		if hash.HSet(field, value) {
			addedOrUpdatedCount++
		}
	}
	return resp.NewInteger(int64(addedOrUpdatedCount))
}

func HGetCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.NewError("ERR wrong number of arguments for 'hget' command")
	}
	key := string(args[0].Bulk)
	field := string(args[1].Bulk)

	val, ok := db.Get(key)
	if !ok {
		return resp.NewNullBulkString()
	}
	hash, isHash := val.(*database.Hash)
	if !isHash {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	value, found := hash.HGet(field)
	if !found {
		return resp.NewNullBulkString()
	}
	return resp.NewBulkString([]byte(value))
}

func HDelCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.NewError("ERR wrong number of arguments for 'hdel' command")
	}
	key := string(args[0].Bulk)
	fields := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		fields[i] = string(arg.Bulk)
	}

	val, ok := db.Get(key)
	if !ok {
		return resp.NewInteger(0)
	}
	hash, isHash := val.(*database.Hash)
	if !isHash {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	deletedCount := hash.HDel(fields...)
	return resp.NewInteger(int64(deletedCount))
}

func HLenCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'hlen' command")
	}
	key := string(args[0].Bulk)

	val, ok := db.Get(key)
	if !ok {
		return resp.NewInteger(0)
	}
	hash, isHash := val.(*database.Hash)
	if !isHash {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return resp.NewInteger(int64(hash.HLen()))
}

func SAddCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.NewError("ERR wrong number of arguments for 'sadd' command")
	}
	key := string(args[0].Bulk)
	members := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		members[i] = string(arg.Bulk)
	}

	val, ok := db.Get(key)
	var set *database.Set
	if !ok {
		set = database.NewSet()
		db.Set(key, set, 0)
	} else {
		existingSet, isSet := val.(*database.Set)
		if !isSet {
			return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		set = existingSet
	}
	addedCount := set.SAdd(members...)
	return resp.NewInteger(int64(addedCount))
}

func SRemCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.NewError("ERR wrong number of arguments for 'srem' command")
	}
	key := string(args[0].Bulk)
	members := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		members[i] = string(arg.Bulk)
	}

	val, ok := db.Get(key)
	if !ok {
		return resp.NewInteger(0)
	}
	set, isSet := val.(*database.Set)
	if !isSet {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	removedCount := set.SRem(members...)
	return resp.NewInteger(int64(removedCount))
}

func SIsMemberCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.NewError("ERR wrong number of arguments for 'sismember' command")
	}
	key := string(args[0].Bulk)
	member := string(args[1].Bulk)

	val, ok := db.Get(key)
	if !ok {
		return resp.NewInteger(0)
	}
	set, isSet := val.(*database.Set)
	if !isSet {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	if set.SIsMember(member) {
		return resp.NewInteger(1)
	}
	return resp.NewInteger(0)
}

func SCardCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'scard' command")
	}
	key := string(args[0].Bulk)

	val, ok := db.Get(key)
	if !ok {
		return resp.NewInteger(0)
	}
	set, isSet := val.(*database.Set)
	if !isSet {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return resp.NewInteger(int64(set.SCard()))
}

func ZAddCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) < 3 || len(args)%2 != 1 {
		return resp.NewError("ERR wrong number of arguments for 'zadd' command")
	}

	key := string(args[0].Bulk)

	val, ok := db.Get(key)
	var zset *database.ZSet
	if !ok {
		zset = database.NewZSet()
		db.Set(key, zset, 0)
	} else {
		existingZSet, isZSet := val.(*database.ZSet)
		if !isZSet {
			return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		zset = existingZSet
	}

	addedCount := 0
	for i := 1; i < len(args); i += 2 {
		scoreStr := string(args[i].Bulk)
		member := string(args[i+1].Bulk)

		score, err := strconv.ParseFloat(scoreStr, 64)
		if err != nil {
			return resp.NewError("ERR value is not a valid float")
		}
		addedCount += zset.ZAdd(score, member)
	}
	return resp.NewInteger(int64(addedCount))
}

func ZScoreCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.NewError("ERR wrong number of arguments for 'zscore' command")
	}
	key := string(args[0].Bulk)
	member := string(args[1].Bulk)

	val, ok := db.Get(key)
	if !ok {
		return resp.NewNullBulkString()
	}
	zset, isZSet := val.(*database.ZSet)
	if !isZSet {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	score, found := zset.ZScore(member)
	if !found {
		return resp.NewNullBulkString()
	}
	return resp.NewBulkString([]byte(strconv.FormatFloat(score, 'f', -1, 64)))
}

func ZRemCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return resp.NewError("ERR wrong number of arguments for 'zrem' command")
	}
	key := string(args[0].Bulk)
	members := make([]string, len(args)-1)
	for i, arg := range args[1:] {
		members[i] = string(arg.Bulk)
	}

	val, ok := db.Get(key)
	if !ok {
		return resp.NewInteger(0)
	}
	zset, isZSet := val.(*database.ZSet)
	if !isZSet {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	removedCount := zset.ZRem(members...)
	return resp.NewInteger(int64(removedCount))
}

func ZCardCommand(db *database.Database, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.NewError("ERR wrong number of arguments for 'zcard' command")
	}
	key := string(args[0].Bulk)

	val, ok := db.Get(key)
	if !ok {
		return resp.NewInteger(0)
	}
	zset, isZSet := val.(*database.ZSet)
	if !isZSet {
		return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return resp.NewInteger(int64(zset.ZCard()))
}
