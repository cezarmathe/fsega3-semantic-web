# semweb

`make local-up` will start the environment with `docker-compose.yml` (replace
`podman` with `docker` inside the `Makefile` if you need).
- `json-server` is available at `http://localhost:4000`.
- `rdf4j-workbench` is available at `http://localhost:8080/rdf4j-workbench`.
- `rdf4j-server` is available at `http://localhost:8080/rdf4j-server`.

`make run-frontend` will start the frontend at `http://localhost:3000`.

`make run-backend` will start the backend at `http://localhost:8000`
