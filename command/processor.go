package command

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/HORUSCRIME/goredis/database"
	"github.com/HORUSCRIME/goredis/resp"
)

type HandlerFunc func(db *database.Database, args []resp.Value) resp.Value

type Processor struct {
	handlers map[string]HandlerFunc
	db       *database.Database 
	mu       sync.RWMutex
}

func NewProcessor(db *database.Database) *Processor {
	p := &Processor{
		db:       db,
		handlers: make(map[string]HandlerFunc),
	}
	p.registerDefaultHandlers()
	return p
}

func (p *Processor) registerDefaultHandlers() {
	p.Register("PING", PingCommand)
	p.Register("ECHO", EchoCommand)
	p.Register("SET", SetCommand)
	p.Register("GET", GetCommand)
	p.Register("DEL", DelCommand)
	p.Register("EXISTS", ExistsCommand)
	p.Register("TYPE", TypeCommand)

	p.Register("LPUSH", LPushCommand)
	p.Register("RPUSH", RPushCommand)
	p.Register("LPOP", LPopCommand)
	p.Register("RPOP", RPopCommand)
	p.Register("LLEN", LLenCommand)

	p.Register("HSET", HSetCommand)
	p.Register("HGET", HGetCommand)
	p.Register("HDEL", HDelCommand)
	p.Register("HLEN", HLenCommand)

	p.Register("SADD", SAddCommand)
	p.Register("SREM", SRemCommand)
	p.Register("SISMEMBER", SIsMemberCommand)
	p.Register("SCARD", SCardCommand)

	p.Register("ZADD", ZAddCommand)
	p.Register("ZSCORE", ZScoreCommand)
	p.Register("ZREM", ZRemCommand)
	p.Register("ZCARD", ZCardCommand)
}

func (p *Processor) Register(cmd string, handler HandlerFunc) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.handlers[strings.ToUpper(cmd)] = handler
}


func (p *Processor) Process(cmdValue resp.Value) resp.Value {
	if cmdValue.Type != resp.ArrayType || len(cmdValue.Array) == 0 {
		return resp.NewError("ERR invalid command format")
	}

	commandName := strings.ToUpper(string(cmdValue.Array[0].Bulk))
	args := cmdValue.Array[1:]

	p.mu.RLock()
	handler, ok := p.handlers[commandName]
	p.mu.RUnlock()

	if !ok {
		return resp.NewError(fmt.Sprintf("ERR unknown command '%s'", commandName))
	}

	for _, arg := range args {
		if arg.Type != resp.BulkStringType {
			return resp.NewError("ERR arguments must be bulk strings")
		}
	}

	log.Printf("Executing command: %s, args: %v", commandName, args)
	return handler(p.db, args)
}
