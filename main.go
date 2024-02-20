//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
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

package main

import (
	"time"

	"github.com/eliona-smart-building-assistant/go-eliona/app"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-eliona/dashboard"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// The main function starts the app by starting all services necessary for this app and waits
// until all services are finished.
func main() {
	log.Info("main", "Starting the app.")

	// Set default database to use boil.*G functions.
	database := db.Database(app.AppName())
	defer database.Close()
	boil.SetDB(database)

	// Set the database logging level.
	if log.Lev() >= log.TraceLevel {
		boil.DebugMode = true
		boil.DebugWriter = log.GetWriter(log.TraceLevel, "database")
	}

	// Necessary to close used init resources, because db.Pool() is used in this app.
	defer db.ClosePool()

	// Init the app before the first run.
	app.Init(db.Pool(), app.AppName(),
		app.ExecSqlFile("conf/init.sql"),
		asset.InitAssetTypeFiles("resources/asset-types/*.json"),
		dashboard.InitWidgetTypeFiles("resources/widget-types/*.json"),
	)

	// Starting the service to collect the data for this app.
	common.WaitForWithOs(
		common.Loop(doAnything, time.Second),
		listenApi,
	)

	log.Info("main", "Terminate the app.")
}
