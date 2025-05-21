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
	"database/sql"
	"log"
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

//
// Todo: Define anything for configuration like structures and methods to read and process configuration
//
