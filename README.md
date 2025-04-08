# RDPAlarm

## Build

Run `./prebuild.sh` at project root folder then build with:

Go Tool Argument: `-ldflags="-s -w -H=windowsgui" -trimpath` in production.

In debug environment: Set environment variable `IS_IN_DEBUG=1` and remove `-H=windowsgui` while compiling.

## Usage

Copy `assets/rdpalert_pushconf.json.example` to `rdpalert_pushconf.json` and modify, then save it within the same directory of executable.

Note: Device key must be within the same format as the application shown.

In config, `notificationLevel` is optional, possible values are one of: `active, passive, timeSensitive`.

Import and enable `assets/RDPAlert.xml` to task scheduler and change the executable path accordingly, then enable it.

## License

 RDPAlarm
 Copyright (C) 2023  Patmeow Limited
 
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU Affero General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 
 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU Affero General Public License for more details.
 
 You should have received a copy of the GNU Affero General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

