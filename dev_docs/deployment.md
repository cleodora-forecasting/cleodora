# Deployment

## Docker

```bash
make clean
make build
cp dist/cleosrv_linux_amd64_v1/cleosrv .
DOCKER_TAG=0.1.0.dev.`git rev-parse --short HEAD`
echo "DOCKER_TAG: ${DOCKER_TAG}"
docker build --tag cleodora/cleodora:${DOCKER_TAG} .
docker run -p 8080:8080 -v cleodora_data:/data cleodora/cleodora:${DOCKER_TAG}
docker push cleodora/cleodora:${DOCKER_TAG}
#docker tag cleodora/cleodora:${DOCKER_TAG} cleodora/cleodora:latest
#docker push cleodora/cleodora:latest
```

Always start the container with a named volume (and keep using the same name,
`-v cleodora_data:/data` in the example below, even when updating the image):

```bash
docker run -p 8080:8080 -v cleodora_data:/data cleodora:VERSION
```

The reason is that this image will use an anonymous volume `/data` by default
to store the data. This means if you just stop a container and start a new one,
you will lose your data (e.g. list of forecasts). There are some other things
you can do to avoid this (but the best and easiest is using a named volume as
mentioned above):

* Use a bind mount.
* Before deleting the old container, start the new one with `--volumes-from`
  option to use the same (anonymous) volume. Then you can delete the old
  container.
* Disaster recovery: Find the anonymous volume and copy the data into a new
  volume. This will only work if the volume hasn't been deleted (luckily
  `docker rm` does not delete such volumes by default). See for example [this
  link](https://github.com/moby/moby/issues/31154#issuecomment-360531460).


## fly.io (demo.cleodora.org)

```bash
make clean
make build
flyctl deploy --local-only # use local Docker to build
```


## Release

* Ensure the changelog and download links on the website are up to date
* Git repo needs to be completely clean. No untracked files!
* Create a GitHub token for goreleaser
  * https://github.com/settings/personal-access-tokens/new
  * Access on the cleodora-forecasting organization and repository cleodora-forecasting/cleodora
  * Give the token no organization permissions and the following repository permissions:
    * Read access to metadata
    * Read and Write access to code

```bash
cp website/content/docs/changelog.md .
# Remove everything except the current release from changelog.md
# Also remove the version title because it becomes redundant
vim changelog.md
GITHUB_TOKEN=secret ./bin/goreleaser release --release-notes changelog.md
```
