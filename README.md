### Seismograph
Performance tests data storage

#### Developing
Install [docker/docker-compose](https://docs.docker.com/compose/install/)

Install deps, you'll need [yarn](https://classic.yarnpkg.com/en/docs/install/#mac-stable)
```
cd back && make install-deps
cd front && yarn install
```
```
cd back

// backend tests
make test

// html cover
make cover

// run with dev frontend & compose, go to http://localhost:3000 in browser
make stop && make start_rebuild

in other tab

docker-compose logs -f

// e2e tests
cd back && make stop && make start_rebuild
cd front && yarn run cypress run

or for cypress UI for debug

cd front && yarn run cypress open
```
### Cloud providers setup:

#### AWS
Get credentials and default region files, put it on deployment host
```
- ${HOME}/.aws/credentials
- ${HOME}/.aws/config
```
### Other
Local Prometheus setup:
```
cd back && ./scripts/local_prometheus.sh
```
Urls:
```
http://localhost:3000 - UI
http://localhost:10500 - API
http://localhost:9000 - MinIO UI
http://localhost:5432 - PostgreSQL
http://localhost:10500/metrics - Prometheus metrics url
http://localhost:10500/debug/pprof/ - Pprof
http://localhost:9090 - local Prometheus UI
```

#### TODO list
- [x] Add report with test metadata for load test pg/minio
- [x] Create tests list page
- [x] Create test review page (echarts-react)
- [x] Create clusters page
- [x] Create initial aws cluster bootstrap
- [ ] Add swaggo for public API
- [ ] Create tests compare functionality
- [ ] Create edPELT changepoint detection
- [ ] Add basic AUT info metadata for tests
- [ ] Add prometheus import
- [ ] Add go benchmarks support