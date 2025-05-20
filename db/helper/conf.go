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

	"github.com/eliona-smart-building-assistant/go-eliona/frontend"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"github.com/lib/pq"

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
		log.Fatal("conf", "Database not initialized. Call InitDB() first.")
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

func UpsertConfig(ctx context.Context, config appmodel.Configuration) (appmodel.Configuration, error) {
	af, err := json.Marshal(config.AssetFilter)
	if err != nil {
		return appmodel.Configuration{}, fmt.Errorf("marshalling asset filter: %w", err)
	}

	stmt := Configuration.INSERT()
	// TODO: There might be a better way to solve this?
	if config.Id != 0 {
		// If ID is provided, include it in the INSERT
		stmt = Configuration.INSERT(
			Configuration.ID,
			Configuration.APIAccessChangeMe,
			Configuration.RefreshInterval,
			Configuration.RequestTimeout,
			Configuration.AssetFilter,
			Configuration.Active,
			Configuration.Enable,
			Configuration.ProjectIds,
			Configuration.UserID,
		).VALUES(
			config.Id,
			config.ApiAccessChangeMe,
			config.RefreshInterval,
			config.RequestTimeout,
			Json(af),
			config.Active,
			config.Enable,
			pq.StringArray(config.ProjectIDs),
			frontend.GetEnvironment(ctx).UserId,
		).ON_CONFLICT(
			Configuration.ID,
		).DO_UPDATE(
			SET(
				Configuration.APIAccessChangeMe.SET(Configuration.EXCLUDED.APIAccessChangeMe),
				Configuration.RefreshInterval.SET(Configuration.EXCLUDED.RefreshInterval),
				Configuration.RequestTimeout.SET(Configuration.EXCLUDED.RequestTimeout),
				Configuration.AssetFilter.SET(Configuration.EXCLUDED.AssetFilter),
				Configuration.Active.SET(Configuration.EXCLUDED.Active),
				Configuration.Enable.SET(Configuration.EXCLUDED.Enable),
				Configuration.ProjectIds.SET(Configuration.EXCLUDED.ProjectIds),
			),
		)
	} else {
		// If ID is 0, omit it to allow auto-increment
		stmt = Configuration.INSERT(
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
			Json(af),
			config.Active,
			config.Enable,
			pq.StringArray(config.ProjectIDs),
			frontend.GetEnvironment(ctx).UserId,
		)
	}

	stmt = stmt.RETURNING(
		Configuration.AllColumns,
	)

	var updatedConfig model.Configuration
	err = stmt.QueryContext(ctx, GetDB().db, &updatedConfig)
	if err != nil {
		return appmodel.Configuration{}, err
	}

	return toAppConfig(updatedConfig)
}

func GetConfig(ctx context.Context, id int64) (appmodel.Configuration, error) {
	var dbConfig model.Configuration
	err := Configuration.
		SELECT(Configuration.AllColumns).
		WHERE(Configuration.ID.EQ(Int(id))).
		QueryContext(ctx, GetDB().db, &dbConfig)
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

func InsertAssetWithDetails(ctx context.Context, config appmodel.Configuration, projectID string, globalAssetID string, assetID int32, providerID string, isRoot bool) error {
	asset := appmodel.Asset{
		Config:        config,
		ProjectID:     projectID,
		GlobalAssetID: globalAssetID,
		ProviderID:    providerID,
		IsRoot:        isRoot,
		AssetID:       assetID,
	}

	return InsertAsset(ctx, asset)
}

func InsertAsset(ctx context.Context, asset appmodel.Asset) error {
	stmt := Asset.INSERT(
		Asset.ConfigurationID,
		Asset.ProjectID,
		Asset.GlobalAssetID,
		Asset.ProviderID,
		Asset.AssetID,
		Asset.IsRoot,
	).VALUES(
		asset.Config.Id,
		asset.ProjectID,
		asset.GlobalAssetID,
		asset.ProviderID,
		asset.AssetID,
		asset.IsRoot,
	)

	_, err := stmt.ExecContext(ctx, GetDB().db)
	return err
}

func GetAssetId(ctx context.Context, config appmodel.Configuration, projectID, gai string) (*int32, error) {
	var dest struct {
		ID int32
	}
	stmt := Asset.SELECT(
		Asset.ID,
	).WHERE(
		Asset.ConfigurationID.EQ(Int(config.Id)).AND(
			Asset.ProjectID.EQ(String(projectID))).AND(
			Asset.GlobalAssetID.EQ(String(gai))),
	)
	err := stmt.QueryContext(ctx, GetDB().db, &dest)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("getting asset ID: %v", err)
	}

	return &dest.ID, nil
}

type assetWithConfig struct {
	model.Asset
	model.Configuration
}

func GetAssetById(assetId int32) (appmodel.Asset, error) {
	var asset assetWithConfig
	err := SELECT(
		Asset.AllColumns, Configuration.AllColumns,
	).FROM(
		Asset.INNER_JOIN(Configuration, Configuration.ID.EQ(Asset.ConfigurationID)),
	).WHERE(
		Asset.ID.EQ(Int32(assetId)),
	).Query(GetDB().db, &asset)
	if errors.Is(err, sql.ErrNoRows) {
		return appmodel.Asset{}, ErrNotFound
	} else if err != nil {
		return appmodel.Asset{}, fmt.Errorf("fetching asset %v: %v", assetId, err)
	}

	config, err := toAppConfig(asset.Configuration)
	if err != nil {
		return appmodel.Asset{}, fmt.Errorf("translating configuration: %v", err)
	}

	return toAppAsset(asset.Asset, config), nil
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
		UserId:            dbCfg.UserID,
	}, nil
}

func toAppAsset(dbAsset model.Asset, config appmodel.Configuration) appmodel.Asset {
	a := appmodel.Asset{
		ID:            dbAsset.ID,
		Config:        config,
		ProjectID:     dbAsset.ProjectID,
		GlobalAssetID: dbAsset.GlobalAssetID,
		ProviderID:    dbAsset.ProviderID,
		IsRoot:        dbAsset.IsRoot,
	}
	if dbAsset.AssetID != nil {
		a.AssetID = *dbAsset.AssetID
	}
	return a
}

func GetRootAssets() ([]appmodel.Asset, error) {
	var assets []assetWithConfig
	err := SELECT(
		Asset.AllColumns, Configuration.AllColumns,
	).FROM(
		Asset.INNER_JOIN(Configuration, Configuration.ID.EQ(Asset.ConfigurationID)),
	).WHERE(
		Asset.IsRoot.EQ(Bool(true)),
	).Query(GetDB().db, &assets)
	if err != nil {
		return nil, fmt.Errorf("fetching root assets: %v", err)
	}

	appAssets := make([]appmodel.Asset, 0, len(assets))
	for _, asset := range assets {
		config, err := toAppConfig(asset.Configuration)
		if err != nil {
			return nil, fmt.Errorf("translating configuration: %v", err)
		}
		appAssets = append(appAssets, toAppAsset(asset.Asset, config))
	}
	return appAssets, nil
}
