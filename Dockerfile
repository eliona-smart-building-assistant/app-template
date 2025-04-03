#  This file is part of the Eliona project.
#  Copyright © 2025 IoTEC AG. All Rights Reserved.
#  ______ _ _
# |  ____| (_)
# | |__  | |_  ___  _ __   __ _
# |  __| | | |/ _ \| '_ \ / _` |
# | |____| | | (_) | | | | (_| |
# |______|_|_|\___/|_| |_|\__,_|
#
#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
#  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
#  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
#  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
#  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

FROM eliona/base-golang:latest-1-alpine AS build

RUN apk add git

WORKDIR /
COPY . ./

RUN go mod download

RUN DATE=$(date) && \
    GIT_COMMIT=$(git rev-list -1 HEAD) && \
    go build -ldflags "-X 'app-name/api/services.BuildTimestamp=$DATE' -X 'app-name/api/services.GitCommit=$GIT_COMMIT'" -o ../app-build

FROM eliona/base-alpine:latest AS target

COPY --from=build /app-build ./
COPY db/*.sql ./db/
COPY resources/ ./resources/
COPY openapi.yaml ./
COPY metadata.json ./

ENV TZ=Europe/Zurich
CMD [ "/app-build" ]
