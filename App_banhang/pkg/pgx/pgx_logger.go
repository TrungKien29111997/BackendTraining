package pgx

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
	"user-management-api/pkg/logger"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type PgxZerologTracer struct {
	Logger         zerolog.Logger
	SlowQueryLimit time.Duration
}

type QueryInfo struct {
	QueryName     string
	OperationType string
	CleanSQL      string
	OriginalSQL   string
}

func parseSQL(sql string) QueryInfo {
	info := QueryInfo{
		OriginalSQL: sql,
	}
	// 0: originalSQL, 1: queryName, 2: operationType
	if match := sqlcNameRegex.FindStringSubmatch(sql); len(match) == 3 {
		info.QueryName = match[1]
		info.OperationType = strings.ToUpper(match[2])
	}

	cleanSQL := commentRegex.ReplaceAllString(sql, "")
	cleanSQL = strings.TrimSpace(cleanSQL)
	cleanSQL = spaceRegex.ReplaceAllString(cleanSQL, " ")
	info.CleanSQL = cleanSQL
	return info
}

func formatArg(arg any) string {

	val := reflect.ValueOf(arg)

	if arg == nil || val.Kind() == reflect.Ptr && val.IsNil() {
		return "NULL"
	}

	if val.Kind() == reflect.Ptr {
		arg = val.Elem().Interface()
	}

	switch v := arg.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case bool:
		return fmt.Sprintf("%t", v)
	case int, int8, int16, int32, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format(time.RFC3339))
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(fmt.Sprintf("%v", v), "'", "''"))
	}
}

func replacePlaceHolders(sql string, args []any) string {
	for i, arg := range args {
		placeHolder := fmt.Sprintf("$%d", i+1)
		sql = strings.ReplaceAll(sql, placeHolder, formatArg(arg))
	}
	return sql
}

var (
	sqlcNameRegex = regexp.MustCompile(`-- name:\s(\w+)\s*:(\w+)`)
	spaceRegex    = regexp.MustCompile(`\s+`)
	commentRegex  = regexp.MustCompile(`-- [^\n\r]*`)
)

func (t *PgxZerologTracer) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	sql, _ := data["sql"].(string)
	args, _ := data["args"].([]any)
	duration, _ := data["duration"].(time.Duration)

	queryInfo := parseSQL(sql)

	var finalSQL string
	if len(args) > 0 {
		finalSQL = replacePlaceHolders(queryInfo.CleanSQL, args)
	} else {
		finalSQL = queryInfo.CleanSQL
	}

	baseLogger := t.Logger.With().
		Str("trace_id", logger.GetTraceID(ctx)).
		Dur("duration", duration).
		Str("sql_original", queryInfo.OriginalSQL).
		Str("sql_clean", finalSQL).
		Str("query_name", queryInfo.QueryName).
		Str("operation_type", queryInfo.OperationType).
		Interface("args", args)

	logger := baseLogger.Logger()

	if msg == "Query" {
		if duration > t.SlowQueryLimit {
			logger.Warn().Str("event", "Slow Query").Msg("slow SQL query")
			return
		}
		logger.Info().Str("event", "Query").Msg("Executed SQL")
		return
	}
}
