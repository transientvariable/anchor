package sqltest

const (
	// RepositorySQL configuration path.
	//
	// String: <root>.repository.sql
	RepositorySQL = "repository.sql"

	// RepositorySQLUsername configuration path.
	//
	// String: <root>.repository.sql.username
	RepositorySQLUsername = RepositorySQL + ".username"

	// RepositorySQLPassword configuration path.
	//
	// String: <root>.repository.sql.password
	RepositorySQLPassword = RepositorySQL + ".password"

	// RepositorySQLConn configuration path.
	//
	// String: <root>.repository.sql.connection
	RepositorySQLConn = RepositorySQL + ".connection"

	// RepositorySQLConnMaxIdle configuration path.
	//
	// String: <root>.repository.sql.connection.maxIdle
	RepositorySQLConnMaxIdle = RepositorySQLConn + ".maxIdle"

	// RepositorySQLConnMaxIdleTime configuration path.
	//
	// String: <root>.repository.sql.connection.maxIdleTime
	RepositorySQLConnMaxIdleTime = RepositorySQLConn + ".maxIdleTime"

	// RepositorySQLConnMaxOpen configuration path.
	//
	// String: <root>.repository.sql.connection.maxOpen
	RepositorySQLConnMaxOpen = RepositorySQLConn + ".maxOpen"

	// RepositorySQLMigrationPath configuration path.
	//
	// String: <root>.repository.sql.migration.path
	RepositorySQLMigrationPath = RepositorySQL + ".migration.path"

	// RepositorySQLMigrationVerboseLogging configuration path.
	//
	// String: <root>.repository.sql.migration.verboseLogging
	RepositorySQLMigrationVerboseLogging = RepositorySQL + ".migration.verboseLogging"

	// RepositorySQLURI configuration path.
	//
	// String: <root>.repository.sql.uri
	RepositorySQLURI = RepositorySQL + ".uri"
)
