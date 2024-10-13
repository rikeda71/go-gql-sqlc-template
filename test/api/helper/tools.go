package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func PostGraphQLRequest(query Query, server *echo.Echo) ([]byte, error) {
	req := httptest.NewRequest(echo.POST, "/graphql", query.RequestBody())
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if got := cmp.Diff(http.StatusOK, rec.Code); got != "" {
		return nil, errors.New("unexpected response code: " + got)
	}
	fmt.Printf("query: %s\nresponse: %s", query, rec.Body.String())
	return rec.Body.Bytes(), nil
}

// ExecuteSQLsFromDir ディレクトリ内のSQLファイルを昇順にソートして実行する
func ExecuteSQLsFromDir(dir string, conn *pgxpool.Pool, purpose string) error {
	fmt.Println("===============================")
	fmt.Println(purpose)
	fmt.Println("===============================")
	migrationFiles, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	sqlFiles := make([]string, 0, len(migrationFiles))
	for _, file := range migrationFiles {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, path.Join(dir, file.Name()))
		}
	}
	// alphabetically sort
	sort.Strings(sqlFiles)
	// migrate database
	var errs []error
	for _, file := range sqlFiles {
		slog.Debug("read file", "filename", file)
		content, err := os.ReadFile(file)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed to read file "+file))
			continue
		}
		arr := strings.Split(strings.Split(string(content), "-- migrate:down")[0], ";")
		for _, str := range arr[:len(arr)-1] {
			if _, err = conn.Exec(context.Background(), str+";"); err != nil {
				errs = append(errs, errors.Wrap(err, "failed to execute sql in "+file))
				continue
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
