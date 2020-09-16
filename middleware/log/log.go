package log

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gin-items/library/setting"
)

var (
	AccessLogger *zap.Logger
	ErrorLogger  *zap.Logger
)

func InitLogger() {
	go initAccessLogger()
	go initErrorLog()
}

func initAccessLogger() {
	cfg := setting.Config().Log.AccessLog
	//日志文件
	fileName := path.Join(cfg.FilePath, cfg.FileName)
	writer := getLogWriter(fileName)
	encoder := getEncoder()
	var level = new(zapcore.Level)
	err := level.UnmarshalText([]byte("info"))
	if err != nil {
		fmt.Printf("init access_logger failed, err:%v\n", err)
		return
	}
	core := zapcore.NewCore(encoder, zapcore.AddSync(writer), level)
	AccessLogger = zap.New(core, zap.AddCaller())
	return
}

func initErrorLog() {
	cfg := setting.Config().Log.ErrorLog
	//日志文件
	fileName := path.Join(cfg.FilePath, cfg.FileName)
	writer := getLogWriter(fileName)
	encoder := getEncoder()
	var level = new(zapcore.Level)
	err := level.UnmarshalText([]byte("info"))
	if err != nil {
		fmt.Printf("init error_logger failed, err:%v\n", err)
		return
	}
	core := zapcore.NewCore(encoder, zapcore.AddSync(writer), level)
	ErrorLogger = zap.New(core, zap.AddCaller())
	return
}

func getLogWriter(filename string) io.Writer {
	// 日志切割，保存15天日志，按天分割日志
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		filename+".%Y%m%d",
		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(filename),
		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return logWriter
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(i time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(i.Format("2006-01-02 15:04:05"))
	}
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder // 日志级别转小写
	encoderConfig.EncodeDuration = func(duration time.Duration, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendInt64(int64(duration) / 1000000)
	}
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // 以包/文件:行号 格式化调用堆栈
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 日志记录
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 执行时间
		costTime := time.Since(startTime)

		logger.Info(
			c.Request.URL.Path,
			zap.Int("status", c.Writer.Status()),
			zap.Duration("cost", costTime),
			zap.String("method", c.Request.Method),
			zap.String("uri", c.Request.RequestURI),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}

// recover掉项目可能出现的panic，并使用zap记录相关日志，替换gin默认的recover
func Recovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

// 日志记录到 MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

// 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

// 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
