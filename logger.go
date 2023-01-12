package logo

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewLogger(coll *mongo.Collection) *Logger {
	return &Logger{
		coll: coll,
	}
}

type Logger struct {
	coll *mongo.Collection
}

func (l *Logger) sendLog(log gin.H) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	print("-------------------")
	_, err := l.coll.InsertOne(ctx, log)
	print(err)
	return err
}

func (l *Logger) LogMiddleware(ctx *gin.Context) {
	req := *ctx.Request
	start := time.Now()
	ctx.Next()
	log := gin.H{
		"method":            req.Method,
		"status":            ctx.Writer.Status(),
		"host":              req.Host,
		"uri":               req.URL.String(),
		"header":            req.Header,
		"content-length":    req.ContentLength,
		"transfer_encoding": req.TransferEncoding,
		"trailer":           req.Trailer,
		"remote_addr":       req.RemoteAddr,
		"remote_uri":        req.RequestURI,
		"date":              time.Now().UTC().Local().Add(time.Hour * 3),
		"latency":           time.Since(start).Seconds(),
	}
	err := l.sendLog(log)
	if err != nil {
		panic(err)
	}
}
