b:
	./scripts/build.sh
	# pack build cnb-test -B gcr.io/buildpacks/builder:v1 -b . -p sample-app
	docker rmi -f cnb-test
	# ./scripts/package.sh -v 2.1.0
	pack build cnb-test -B gcr.io/buildpacks/builder:v1 -b . -p sample-app

c:
	docker rmi -f cnb-test
	pack build cnb-test -B gcr.io/buildpacks/builder:v1 -b ghcr.io/cage1016/aaa:latest -p sample-app