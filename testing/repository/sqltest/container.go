package sqltest

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/transientvariable/config"
	"github.com/transientvariable/support"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	pg "github.com/testcontainers/testcontainers-go/modules/postgres"

	_ "github.com/transientvariable/support/testing"
)

const (
	database = "defaultdb"
	username = "defaultuser"
	password = "defaultpassword"
	image    = "postgres:17"
)

type PostgresContainer struct {
	closed        bool
	container     *pg.PostgresContainer
	ctx           context.Context
	ctxCancel     context.CancelFunc
	ctxParent     context.Context
	database      string
	image         string
	migrationsDir string
	mutex         sync.Mutex
	uri           *url.URL
	userInfo      *url.Userinfo
}

func NewPostgresContainer(options ...func(*PostgresContainer)) (*PostgresContainer, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	pgc := &PostgresContainer{
		migrationsDir: filepath.Join(dir, config.ValueMustResolve(RepositorySQLMigrationPath)),
		userInfo:      userInfo(),
	}

	for _, opt := range options {
		opt(pgc)
	}

	if pgc.ctxParent == nil {
		pgc.ctxParent = context.Background()
	}
	pgc.ctx, pgc.ctxCancel = context.WithCancel(pgc.ctxParent)

	if pgc.database == "" {
		pgc.database = database
	}

	if pgc.image == "" {
		pgc.image = image
	}

	u := pgc.userInfo.Username()
	p, _ := pgc.userInfo.Password()

	opts := []testcontainers.ContainerCustomizer{
		pg.WithDatabase(pgc.database),
		pg.WithUsername(u),
		pg.WithPassword(p),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(5 * time.Second)),
	}

	c, err := pg.Run(pgc.ctx, pgc.image, opts...)
	if err != nil {
		return nil, fmt.Errorf("postgres_test_container: %w", err)
	}
	pgc.container = c

	s, err := c.ConnectionString(pgc.ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("postgres_test_container: %w", err)
	}

	pgc.uri, err = url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("postgres_test_container: %w", err)
	}
	return pgc, nil
}

func (c *PostgresContainer) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.closed {
		c.closed = true
		defer c.ctxCancel()
		if err := c.container.Terminate(c.ctx); err != nil {
			return fmt.Errorf("postgres_test_container: %w", err)
		}
		return nil
	}
	return errors.New("postgres_test_container: already closed")
}

func (c *PostgresContainer) Container() *pg.PostgresContainer {
	return c.container
}

func (c *PostgresContainer) Database() string {
	return c.database
}

func (c *PostgresContainer) Image() string {
	return c.image
}

func (c *PostgresContainer) MigrationsDir() string {
	return c.migrationsDir
}

func (c *PostgresContainer) Password() (string, bool) {
	if c.userInfo == nil {
		return c.userInfo.Password()
	}
	return "", false
}

func (c *PostgresContainer) URI() (*url.URL, error) {
	if c.uri != nil {
		return url.Parse(c.uri.String())
	}
	return nil, errors.New("postgres_test_container: uri is nil")
}

func (c *PostgresContainer) Username() string {
	if c.userInfo == nil {
		return c.userInfo.Username()
	}
	return ""
}

func (c *PostgresContainer) String() string {
	m := map[string]any{
		"database":      c.database,
		"image":         c.image,
		"migrationsDir": c.migrationsDir,
	}

	if c.userInfo != nil && c.userInfo.Username() != "" {
		m["userinfo"] = map[string]any{
			"username": c.userInfo.Username(),
		}
	}

	if c.uri != nil {
		m["uri"] = c.uri.String()
	}
	return string(support.ToJSONFormatted(m))
}

func userInfo() *url.Userinfo {
	if u, _ := config.Value(RepositorySQLUsername); u != "" {
		if p, _ := config.Value(RepositorySQLPassword); p != "" {
			return url.UserPassword(u, p)
		}
	}
	return url.UserPassword(username, password)
}

func WithContext(ctx context.Context) func(*PostgresContainer) {
	return func(c *PostgresContainer) {
		c.ctxParent = ctx
	}
}

func WithDatabase(database string) func(*PostgresContainer) {
	return func(c *PostgresContainer) {
		c.database = strings.TrimSpace(database)
	}
}

func WithImage(image string) func(*PostgresContainer) {
	return func(c *PostgresContainer) {
		c.image = strings.TrimSpace(image)
	}
}
