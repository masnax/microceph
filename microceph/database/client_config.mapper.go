package database

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/canonical/lxd/lxd/db/query"
	"github.com/canonical/lxd/shared/api"
	"github.com/canonical/microcluster/v2/cluster"
)

var _ = api.ServerEnvironment{}

var clientConfigItemObjects = cluster.RegisterStmt(`
SELECT client_config.id, core_cluster_members.name AS host, client_config.key, client_config.value
  FROM client_config
  JOIN core_cluster_members ON client_config.member_id = core_cluster_members.id
  ORDER BY core_cluster_members.id, client_config.key
`)

var clientConfigItemObjectsByKey = cluster.RegisterStmt(`
SELECT client_config.id, core_cluster_members.name AS host, client_config.key, client_config.value
  FROM client_config
  JOIN core_cluster_members ON client_config.member_id = core_cluster_members.id
  WHERE ( client_config.key = ? )
  ORDER BY core_cluster_members.id, client_config.key
`)

var clientConfigItemObjectsByHost = cluster.RegisterStmt(`
SELECT client_config.id, core_cluster_members.name AS host, client_config.key, client_config.value
  FROM client_config
  JOIN core_cluster_members ON client_config.member_id = core_cluster_members.id
  WHERE ( host = ? )
  ORDER BY core_cluster_members.id, client_config.key
`)

var clientConfigItemObjectsByKeyAndHost = cluster.RegisterStmt(`
SELECT client_config.id, core_cluster_members.name AS host, client_config.key, client_config.value
  FROM client_config
  JOIN core_cluster_members ON client_config.member_id = core_cluster_members.id
  WHERE ( client_config.key = ? AND host = ? )
  ORDER BY core_cluster_members.id, client_config.key
`)

var clientConfigItemID = cluster.RegisterStmt(`
SELECT client_config.id FROM client_config
  JOIN core_cluster_members ON client_config.member_id = core_cluster_members.id
  WHERE core_cluster_members.name = ? AND client_config.key = ?
`)

var clientConfigItemCreate = cluster.RegisterStmt(`
INSERT INTO client_config (member_id, key, value)
  VALUES ((SELECT core_cluster_members.id FROM core_cluster_members WHERE core_cluster_members.name = ?), ?, ?)
`)

var clientConfigItemDeleteByKey = cluster.RegisterStmt(`
DELETE FROM client_config WHERE key = ?
`)

var clientConfigItemDeleteByHost = cluster.RegisterStmt(`
DELETE FROM client_config WHERE member_id = (SELECT core_cluster_members.id FROM core_cluster_members WHERE core_cluster_members.name = ?)
`)

var clientConfigItemDeleteByKeyAndHost = cluster.RegisterStmt(`
DELETE FROM client_config WHERE key = ? AND member_id = (SELECT core_cluster_members.id FROM core_cluster_members WHERE core_cluster_members.name = ?)
`)

var clientConfigItemUpdate = cluster.RegisterStmt(`
UPDATE client_config
  SET member_id = (SELECT core_cluster_members.id FROM core_cluster_members WHERE core_cluster_members.name = ?), key = ?, value = ?
 WHERE id = ?
`)

// clientConfigItemColumns returns a string of column names to be used with a SELECT statement for the entity.
// Use this function when building statements to retrieve database entries matching the ClientConfigItem entity.
func clientConfigItemColumns() string {
	return "client_config.id, core_cluster_members.name AS host, client_config.key, client_config.value"
}

// getClientConfigItems can be used to run handwritten sql.Stmts to return a slice of objects.
func getClientConfigItems(ctx context.Context, stmt *sql.Stmt, args ...any) ([]ClientConfigItem, error) {
	objects := make([]ClientConfigItem, 0)

	dest := func(scan func(dest ...any) error) error {
		c := ClientConfigItem{}
		err := scan(&c.ID, &c.Host, &c.Key, &c.Value)
		if err != nil {
			return err
		}

		objects = append(objects, c)

		return nil
	}

	err := query.SelectObjects(ctx, stmt, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"client_config\" table: %w", err)
	}

	return objects, nil
}

// getClientConfigItemsRaw can be used to run handwritten query strings to return a slice of objects.
func getClientConfigItemsRaw(ctx context.Context, tx *sql.Tx, sql string, args ...any) ([]ClientConfigItem, error) {
	objects := make([]ClientConfigItem, 0)

	dest := func(scan func(dest ...any) error) error {
		c := ClientConfigItem{}
		err := scan(&c.ID, &c.Host, &c.Key, &c.Value)
		if err != nil {
			return err
		}

		objects = append(objects, c)

		return nil
	}

	err := query.Scan(ctx, tx, sql, dest, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"client_config\" table: %w", err)
	}

	return objects, nil
}

// GetClientConfigItems returns all available ClientConfigItems.
// generator: ClientConfigItem GetMany
func GetClientConfigItems(ctx context.Context, tx *sql.Tx, filters ...ClientConfigItemFilter) ([]ClientConfigItem, error) {
	var err error

	// Result slice.
	objects := make([]ClientConfigItem, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var sqlStmt *sql.Stmt
	args := []any{}
	queryParts := [2]string{}

	if len(filters) == 0 {
		sqlStmt, err = cluster.Stmt(tx, clientConfigItemObjects)
		if err != nil {
			return nil, fmt.Errorf("Failed to get \"clientConfigItemObjects\" prepared statement: %w", err)
		}
	}

	for i, filter := range filters {
		if filter.Key != nil && filter.Host != nil {
			args = append(args, []any{filter.Key, filter.Host}...)
			if len(filters) == 1 {
				sqlStmt, err = cluster.Stmt(tx, clientConfigItemObjectsByKeyAndHost)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"clientConfigItemObjectsByKeyAndHost\" prepared statement: %w", err)
				}

				break
			}

			query, err := cluster.StmtString(clientConfigItemObjectsByKeyAndHost)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"clientConfigItemObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Key != nil && filter.Host == nil {
			args = append(args, []any{filter.Key}...)
			if len(filters) == 1 {
				sqlStmt, err = cluster.Stmt(tx, clientConfigItemObjectsByKey)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"clientConfigItemObjectsByKey\" prepared statement: %w", err)
				}

				break
			}

			query, err := cluster.StmtString(clientConfigItemObjectsByKey)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"clientConfigItemObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Host != nil && filter.Key == nil {
			args = append(args, []any{filter.Host}...)
			if len(filters) == 1 {
				sqlStmt, err = cluster.Stmt(tx, clientConfigItemObjectsByHost)
				if err != nil {
					return nil, fmt.Errorf("Failed to get \"clientConfigItemObjectsByHost\" prepared statement: %w", err)
				}

				break
			}

			query, err := cluster.StmtString(clientConfigItemObjectsByHost)
			if err != nil {
				return nil, fmt.Errorf("Failed to get \"clientConfigItemObjects\" prepared statement: %w", err)
			}

			parts := strings.SplitN(query, "ORDER BY", 2)
			if i == 0 {
				copy(queryParts[:], parts)
				continue
			}

			_, where, _ := strings.Cut(parts[0], "WHERE")
			queryParts[0] += "OR" + where
		} else if filter.Host == nil && filter.Key == nil {
			return nil, fmt.Errorf("Cannot filter on empty ClientConfigItemFilter")
		} else {
			return nil, fmt.Errorf("No statement exists for the given Filter")
		}
	}

	// Select.
	if sqlStmt != nil {
		objects, err = getClientConfigItems(ctx, sqlStmt, args...)
	} else {
		queryStr := strings.Join(queryParts[:], "ORDER BY")
		objects, err = getClientConfigItemsRaw(ctx, tx, queryStr, args...)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"client_config\" table: %w", err)
	}

	return objects, nil
}

// GetClientConfigItem returns the ClientConfigItem with the given key.
// generator: ClientConfigItem GetOne
func GetClientConfigItem(ctx context.Context, tx *sql.Tx, host string, key string) (*ClientConfigItem, error) {
	filter := ClientConfigItemFilter{}
	filter.Host = &host
	filter.Key = &key

	objects, err := GetClientConfigItems(ctx, tx, filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"client_config\" table: %w", err)
	}

	switch len(objects) {
	case 0:
		return nil, api.StatusErrorf(http.StatusNotFound, "ClientConfigItem not found")
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one \"client_config\" entry matches")
	}
}

// GetClientConfigItemID return the ID of the ClientConfigItem with the given key.
// generator: ClientConfigItem ID
func GetClientConfigItemID(ctx context.Context, tx *sql.Tx, host string, key string) (int64, error) {
	stmt, err := cluster.Stmt(tx, clientConfigItemID)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"clientConfigItemID\" prepared statement: %w", err)
	}

	row := stmt.QueryRowContext(ctx, host, key)
	var id int64
	err = row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, api.StatusErrorf(http.StatusNotFound, "ClientConfigItem not found")
	}

	if err != nil {
		return -1, fmt.Errorf("Failed to get \"client_config\" ID: %w", err)
	}

	return id, nil
}

// ClientConfigItemExists checks if a ClientConfigItem with the given key exists.
// generator: ClientConfigItem Exists
func ClientConfigItemExists(ctx context.Context, tx *sql.Tx, host string, key string) (bool, error) {
	_, err := GetClientConfigItemID(ctx, tx, host, key)
	if err != nil {
		if api.StatusErrorCheck(err, http.StatusNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// CreateClientConfigItem adds a new ClientConfigItem to the database.
// generator: ClientConfigItem Create
func CreateClientConfigItem(ctx context.Context, tx *sql.Tx, object ClientConfigItem) (int64, error) {
	// Check if a ClientConfigItem with the same key exists.
	exists, err := ClientConfigItemExists(ctx, tx, object.Host, object.Key)
	if err != nil {
		return -1, fmt.Errorf("Failed to check for duplicates: %w", err)
	}

	if exists {
		return -1, api.StatusErrorf(http.StatusConflict, "This \"client_config\" entry already exists")
	}

	args := make([]any, 3)

	// Populate the statement arguments.
	args[0] = object.Host
	args[1] = object.Key
	args[2] = object.Value

	// Prepared statement to use.
	stmt, err := cluster.Stmt(tx, clientConfigItemCreate)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"clientConfigItemCreate\" prepared statement: %w", err)
	}

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf("Failed to create \"client_config\" entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to fetch \"client_config\" entry ID: %w", err)
	}

	return id, nil
}

// DeleteClientConfigItem deletes the ClientConfigItem matching the given key parameters.
// generator: ClientConfigItem DeleteOne-by-Key-and-Host
func DeleteClientConfigItem(ctx context.Context, tx *sql.Tx, key string, host string) error {
	stmt, err := cluster.Stmt(tx, clientConfigItemDeleteByKeyAndHost)
	if err != nil {
		return fmt.Errorf("Failed to get \"clientConfigItemDeleteByKeyAndHost\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(key, host)
	if err != nil {
		return fmt.Errorf("Delete \"client_config\": %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n == 0 {
		return api.StatusErrorf(http.StatusNotFound, "ClientConfigItem not found")
	} else if n > 1 {
		return fmt.Errorf("Query deleted %d ClientConfigItem rows instead of 1", n)
	}

	return nil
}

// DeleteClientConfigItems deletes the ClientConfigItem matching the given key parameters.
// generator: ClientConfigItem DeleteMany-by-Key
func DeleteClientConfigItems(ctx context.Context, tx *sql.Tx, key string) error {
	stmt, err := cluster.Stmt(tx, clientConfigItemDeleteByKey)
	if err != nil {
		return fmt.Errorf("Failed to get \"clientConfigItemDeleteByKey\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(key)
	if err != nil {
		return fmt.Errorf("Delete \"client_config\": %w", err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	return nil
}

// UpdateClientConfigItem updates the ClientConfigItem matching the given key parameters.
// generator: ClientConfigItem Update
func UpdateClientConfigItem(ctx context.Context, tx *sql.Tx, host string, key string, object ClientConfigItem) error {
	id, err := GetClientConfigItemID(ctx, tx, host, key)
	if err != nil {
		return err
	}

	stmt, err := cluster.Stmt(tx, clientConfigItemUpdate)
	if err != nil {
		return fmt.Errorf("Failed to get \"clientConfigItemUpdate\" prepared statement: %w", err)
	}

	result, err := stmt.Exec(object.Host, object.Key, object.Value, id)
	if err != nil {
		return fmt.Errorf("Update \"client_config\" entry failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query updated %d rows instead of 1", n)
	}

	return nil
}
