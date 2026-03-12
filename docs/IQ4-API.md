# Rain Bird IQ4 cloud API reference

Research notes on the IQ4 cloud platform API used by the Rain Bird 2.0 app.

## Overview

Rain Bird 2.0 app migrated controller management to the IQ4 cloud platform. Schedule/program data is no longer accessible via the local controller API (CDT protocol returns empty responses on firmware 4.98+). Basic operations (zone control, water budget, rain delay) still work locally, but reading/writing schedules requires the cloud API.

## Authentication

### Flow

OpenID Connect implicit flow against IdentityServer4:

1. **GET** login page with OIDC params
2. **POST** credentials + antiforgery token
3. Follow redirects to get `access_token` from fragment

### Details

- **Identity server**: `https://iq4server.rainbird.com/coreidentityserver`
- **Client ID**: `C5A6F324-3CD3-4B22-9F78-B4835BA55D25`
- **Redirect URI**: `https://iq4.rainbird.com/auth.html`
- **Scopes**: `coreAPI.read coreAPI.write openid profile`
- **Response type**: `id_token token`
- **Token type**: JWT Bearer
- **Token lifetime**: ~2 hours (exp claim)

### AWS WAF challenge

The login endpoint is protected by AWS WAF JavaScript challenge (`x-amzn-waf-action: challenge`). Simple HTTP clients (curl, fetch) may get blocked with HTTP 202. The Go client in this project handles the flow when WAF is not active.

### Username format

Username is the Rain Bird account username (not the email address).

## API base URL

```
https://iq4server.rainbird.com/coreapi/api/
```

All requests require `Authorization: Bearer <jwt_token>` header.

## Data model

```
Company → Sites → Controllers (Satellites)
                    ├── Stations (physical valve zones, 4 or 8 per controller)
                    └── Programs (A, B, C – irrigation schedules)
                        ├── Start times (when to run, multiple per program)
                        ├── Program steps (station → runtime mapping)
                        └── Seasonal adjust (% scaling of runtimes)
```

## Read endpoints

### Sites

**GET** `/Site/GetSites`

Returns array of sites for the current company.

### Controllers

**GET** `/Satellite/GetSatelliteList`

Returns all controllers across all sites. Note: the `isConnected` field here is unreliable for MQTT controllers.

**GET** `/Satellite/isConnected?satelliteIds=X&satelliteIds=Y`

Returns real-time MQTT connection status. Use this instead of the field from `GetSatelliteList`.

**GET** `/Satellite/GetSatellite?satelliteId={id}`

Returns detailed controller info including capabilities.

### Stations

**GET** `/Station/GetStationListForSatellite?satelliteId={id}`

Returns station list with adjustment factors, flow rates, landscape/sprinkler type, etc.

Notable fields:
- `areaLevel2Name` – landscape type (e.g. "Grass")
- `areaLevel3Name` – sprinkler type (e.g. "Rotors/Impacts", "Spray Heads")
- `arc` – spray arc in degrees
- `soilType` – soil classification
- `precRateFinal` – precipitation rate
- `cropCoefficient` – crop coefficient for ET calculations

**GET** `/Station/GetStation?stationId={id}`

Returns full station detail. Note: some fields (like area levels) are only present in the list endpoint.

### Programs

**GET** `/Program/GetProgramList?satelliteId={id}`

Returns all programs for a controller.

**GET** `/Program/GetProgramListForMultiSites`

Returns all programs across all sites.

**GET** `/Program/GetProgram?programId={id}`

Full program detail. Notable fields:
- `programAdjust` – seasonal adjustment percentage (e.g. 130 = 130%)
- `weekDays` – 7-char binary string, **Sunday-first**: Su Mo Tu We Th Fr Sa
- `type` – 0 = Week Days, 2 = Odd, 3 = Odd31, 4 = Even, 5 = Cyclic, 6 = Month dates, 7 = Calendar date
- `programStep` – always empty in response; use separate endpoints

### Start times

**GET** `/Program/GetScheduledStartTimes?satelliteId={id}`

Returns start times grouped by program. Date portion (1999-09-09) is meaningless – only the time part matters.

**GET** `/StartTime/GetAllStartTimes?includeProgram=false&includeProgramGroup=false`

Returns all start times across all programs.

### Run times / program steps

**GET** `/ProgramStep/GetProgramsAssignedAndRunTimeBySatelliteId?satelliteId={id}`

Returns runtime assignments per station, including both base and adjusted runtimes.

**GET** `/ProgramStep/GetProgramStepById?programStepId={id}`

Returns the full program step detail (needed for updates).

### Other endpoints

- **GET** `/ProgramGroup/GetProgramTypes` – frequency type enum
- **GET** `/User/GetUser` – current user info
- **GET** `/User/GetUserCurrentDateTime` – server time
- **GET** `/Company/GetCompanyPreferences` – company settings
- **GET** `/CultureSettings/GetCultureSettingDetail` – locale settings
- **GET** `/WeatherSource/GetWeatherSource` – weather data source
- **GET** `/License/IsCloud` – cloud license check

## Write endpoints

| Operation | Method | Endpoint | Body |
|-----------|--------|----------|------|
| Update program | `PUT` | `/Program/UpdateProgram` | Full program object from `GetProgram` |
| Update program step | `PUT` | `/ProgramStep/UpdateProgramStep` | Full step from `GetProgramStepById` |
| Create start time | `POST` | `/StartTime/CreateStartTime` | Start time object |
| Delete start time | `PATCH` | `/StartTime/v2/UpdateBatches` | `{"add":[],"update":[],"delete":{"id":<programId>,"ids":[<startTimeId>]}}` |
| Create program steps | `POST` | `/ProgramStep/CreateProgramSteps` | `[{"actionId":"RunStation","programId":"<id>","runTimeLong":null,"stationId":<id>}]` |
| Delete program steps | `DELETE` | `/ProgramStep/DeleteProgramSteps` | Array of step IDs |

### Important notes on write operations

- **CreateProgramSteps**: `actionId` must be the string `"RunStation"` (not an int), `programId` must be a string, `runTimeLong` should be `null`. Runtime is set separately via `UpdateProgramStep`.
- **DeleteStartTime**: The `DELETE /StartTime/DeleteStartTime` endpoint returns 403 for some start times. Use `PATCH /StartTime/v2/UpdateBatches` instead – it works reliably.
- **UpdateProgram**: Send the full program object back (GET it first, modify fields, PUT it back).
- **runTimeLong**: Uses .NET ticks (100-nanosecond units). 10 minutes = 6000000000.

## Local API – what still works

These SIP commands work on RC2 with firmware 4.98:
- `get_model_and_version()` – model and version
- `get_available_stations()` – available stations
- `get_zone_states()` – active zones
- `irrigate_zone(zone, minutes)` – run a zone
- `stop_irrigation()` – stop all
- `water_budget(program)` – seasonal adjust
- `get_rain_delay()` / `set_rain_delay(days)`
- `get_rain_sensor_state()`
- `get_current_irrigation()`
- `get_wifi_params()` – MAC, IP, RSSI, SSID
- `get_serial_number()`
- `get_settings()` – program count, location
- `get_controller_firmware_version()`

## What does NOT work locally

- `get_schedule()` – RC2 NAKs the legacy SIP `RetrieveScheduleRequest` (0x20)
- CDT batch get commands (Universal protocol) – return empty data blocks
- Schedule data is managed exclusively via the IQ4 cloud on firmware 4.98+

## References

- [HA issue #142123](https://github.com/home-assistant/core/issues/142123) – RC2 2.0 migration, community API discovery
- [pyrainbird issue #481](https://github.com/allenporter/pyrainbird/issues/481) – ESP-ME3 schedule support
- [IQ4 BMS/API training PDF](https://www.rainbird.com/sites/default/files/media/documents/2023-10/iq4_bms-api_training_0.pdf) – official API overview
- IQ4 web app: `https://iq4.rainbird.com/`
