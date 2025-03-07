//  This file is part of the Eliona project.
//  Copyright © 2025 IoTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package dbhelper

import (
	appmodel "app-name/app/model"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/eliona-smart-building-assistant/go-eliona/frontend"

	"app-name/db/generated/postgres/app_schema_name/model"
	. "app-name/db/generated/postgres/app_schema_name/table"

	. "github.com/go-jet/jet/v2/postgres"
)

// DBHelper is a singleton struct managing the database connection and queries.
type DBHelper struct {
	db *sql.DB
}

var (
	instance *DBHelper
)

// InitDB initializes the database connection ONCE.
func InitDB(db *sql.DB) {
	instance = &DBHelper{
		db: db,
	}
}

// GetDB returns the singleton database instance.
func GetDB() *DBHelper {
	if instance == nil {
		log.Fatal("Database not initialized. Call InitDB() first.")
	}
	return instance
}

// CloseDB gracefully shuts down the database connection.
func CloseDB() error {
	if instance != nil && instance.db != nil {
		return instance.db.Close()
	}
	return nil
}

var ErrBadRequest = errors.New("bad request")
var ErrNotFound = errors.New("not found")

func InsertConfig(ctx context.Context, config appmodel.Configuration) error {
	af, err := json.Marshal(config.AssetFilter)
	if err != nil {
		return fmt.Errorf("marshalling asset filter: %w", err)
	}

	stmt := Configuration.INSERT(
		Configuration.APIAccessChangeMe,
		Configuration.RefreshInterval,
		Configuration.RequestTimeout,
		Configuration.AssetFilter,
		Configuration.Active,
		Configuration.Enable,
		Configuration.ProjectIds,
		Configuration.UserID,
	).VALUES(
		config.ApiAccessChangeMe,
		config.RefreshInterval,
		config.RequestTimeout,
		af,
		config.Active,
		config.Enable,
		config.ProjectIDs,
		frontend.GetEnvironment(ctx).UserId,
	)

	_, err = stmt.ExecContext(ctx, GetDB().db)
	return err
}

func GetConfig(ctx context.Context, id int64) (appmodel.Configuration, error) {
	var dbConfig model.Configuration
	stmt := Configuration.SELECT(Configuration.AllColumns).
		WHERE(Configuration.ID.EQ(Int(id)))

	err := stmt.QueryContext(ctx, GetDB().db, &dbConfig)
	if errors.Is(err, sql.ErrNoRows) {
		return appmodel.Configuration{}, ErrNotFound
	} else if err != nil {
		return appmodel.Configuration{}, err
	}

	return toAppConfig(dbConfig)
}

func DeleteConfig(ctx context.Context, id int64) error {
	stmt := Configuration.DELETE().
		WHERE(Configuration.ID.EQ(Int(id)))

	r, err := stmt.ExecContext(ctx, GetDB().db)
	if err != nil {
		return err
	}
	rows, _ := r.RowsAffected()
	if rows == 0 {
		return ErrNotFound
	}
	if rows > 1 {
		return fmt.Errorf("unexpected deletion: deleted %d rows", rows)
	}
	return nil
}

func GetConfigs(ctx context.Context) ([]appmodel.Configuration, error) {
	var dbConfigs []model.Configuration
	stmt := Configuration.SELECT(Configuration.AllColumns)

	err := stmt.QueryContext(ctx, GetDB().db, &dbConfigs)
	if err != nil {
		return nil, err
	}

	configs := make([]appmodel.Configuration, len(dbConfigs))
	for i, dbCfg := range dbConfigs {
		configs[i], err = toAppConfig(dbCfg)
		if err != nil {
			return nil, err
		}
	}
	return configs, nil
}

func SetConfigActiveState(ctx context.Context, id int64, state bool) error {
	stmt := Configuration.UPDATE(Configuration.Active).
		SET(state).
		WHERE(Configuration.ID.EQ(Int(id)))

	_, err := stmt.ExecContext(ctx, GetDB().db)
	return err
}

func InsertAsset(ctx context.Context, asset appmodel.Asset) error {
	stmt := Asset.INSERT(
		Asset.ConfigurationID,
		Asset.ProjectID,
		Asset.GlobalAssetID,
		Asset.ProviderID,
		Asset.AssetID,
	).VALUES(
		asset.Config.Id,
		asset.ProjectID,
		asset.GlobalAssetID,
		asset.ProviderID,
		asset.AssetID,
	)

	_, err := stmt.ExecContext(ctx, GetDB().db)
	return err
}

func GetAssetById(ctx context.Context, assetId int64) (appmodel.Asset, error) {
	var dbAsset model.Asset
	stmt := Asset.SELECT(Asset.AllColumns).WHERE(Asset.ID.EQ(Int(assetId)))
	err := stmt.QueryContext(ctx, GetDB().db, &dbAsset)
	if errors.Is(err, sql.ErrNoRows) {
		return appmodel.Asset{}, ErrNotFound
	} else if err != nil {
		return appmodel.Asset{}, err
	}

	config, err := GetConfig(ctx, dbAsset.ConfigurationID)
	if err != nil {
		return appmodel.Asset{}, err
	}

	return toAppAsset(dbAsset, config), nil
}

func toAppConfig(dbCfg model.Configuration) (appmodel.Configuration, error) {
	var assetFilter [][]appmodel.FilterRule
	err := json.Unmarshal([]byte(dbCfg.AssetFilter), &assetFilter)
	if err != nil {
		return appmodel.Configuration{}, err
	}

	return appmodel.Configuration{
		Id:                dbCfg.ID,
		ApiAccessChangeMe: dbCfg.APIAccessChangeMe,
		RefreshInterval:   dbCfg.RefreshInterval,
		RequestTimeout:    dbCfg.RequestTimeout,
		AssetFilter:       assetFilter,
		Active:            dbCfg.Active,
		Enable:            dbCfg.Enable,
		ProjectIDs:        dbCfg.ProjectIds,
		// https://github.com/go-jet/jet/issues/351
		// https://github.com/go-jet/jet/issues/53
		UserId: dbCfg.UserID,
	}, nil
}

func toAppAsset(dbAsset model.Asset, config appmodel.Configuration) appmodel.Asset {
	a := appmodel.Asset{
		ID:            dbAsset.ID,
		Config:        config,
		ProjectID:     dbAsset.ProjectID,
		GlobalAssetID: dbAsset.GlobalAssetID,
		ProviderID:    dbAsset.ProviderID,
	}
	if dbAsset.AssetID != nil {
		a.AssetID = *dbAsset.AssetID
	}
	return a
}

func GetRootAssets() ([]appmodel.Asset, error) {
	assets, err := dbgen.Assets(
		dbgen.AssetWhere.IsRoot.EQ(true),
	).AllG(context.Background())
	if err != nil {
		return nil, fmt.Errorf("fetching root assets: %v", err)
	}

	appAssets := make([]appmodel.Asset, 0, len(assets))
	for _, asset := range assets {
		c, err := asset.Configuration().OneG(context.Background())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		if err != nil {
			return nil, fmt.Errorf("fetching configuration: %v", err)
		}
		config, err := toAppConfig(c)
		if err != nil {
			return nil, fmt.Errorf("translating configuration: %v", err)
		}
		appAssets = append(appAssets, toAppAsset(*asset, config))
	}
	return appAssets, nil
}
