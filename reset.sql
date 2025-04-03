--  This file is part of the Eliona project.
--  Copyright © 2025 IoTEC AG. All Rights Reserved.
--  ______ _ _
-- |  ____| (_)
-- | |__  | |_  ___  _ __   __ _
-- |  __| | | |/ _ \| '_ \ / _` |
-- | |____| | | (_) | | | | (_| |
-- |______|_|_|\___/|_| |_|\__,_|
--
--  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
--  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
--  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
--  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
--  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

-- This idempotent script resets the database to a defined state ready for testing.
-- The only thing that remains after testing then are the incremented auto-increment values and app
-- registration (which you can optionally remove as well by uncommenting the last command).

SET SCHEMA 'public';

DELETE FROM versioning.patches
WHERE app_name = 'app-name';

INSERT INTO public.eliona_store (app_name, category, version)
VALUES ('app-name', 'app', 'v0.0.0')
ON CONFLICT (app_name) DO UPDATE SET version = 'v0.0.0';

INSERT INTO public.eliona_app (app_name, enable)
VALUES ('app-name', 't')
ON CONFLICT (app_name) DO UPDATE SET initialized_at = null;

DROP SCHEMA IF EXISTS app_schema_name CASCADE;

DELETE FROM heap
WHERE asset_id IN (
	SELECT asset_id
	FROM asset
	WHERE asset_type LIKE E'app\\_schema\\_name\\_%'
);

DELETE FROM attribute_schema
WHERE asset_type LIKE E'app\\_schema\\_name\\_%';

DELETE FROM asset
WHERE asset_type LIKE E'app\\_schema\\_name\\_%';

DELETE FROM asset_type
WHERE asset_type LIKE E'app\\_schema\\_name\\_%';

DELETE FROM public.widget_data
WHERE widget_id IN (
	SELECT public.widget.id
	FROM public.widget
		JOIN public.dashboard USING (dashboard_id)
	WHERE public.dashboard.name LIKE 'App Name%'
);

DELETE FROM public.widget
WHERE dashboard_id IN (
	SELECT dashboard_id
	FROM public.dashboard
	WHERE name LIKE 'App Name%'
);

DELETE FROM public.dashboard
WHERE name LIKE 'App Name%';

-- DELETE FROM eliona_app WHERE app_name = 'app-name';
-- DELETE FROM eliona_store WHERE app_name = 'app-name';
