# Code Structure

* **[cleoc/](../cleoc/)**: Contains the code for the `cleoc` CLI client. Get
  started with the `cmd` folder.
* **[cleosrv/](../cleosrv/)**: Contains the code for the `cleosrv` application,
  including the GraphQL API. It does not include the React frontend (it is
  bundled during the build). Get started with `main.go` .
* **[cleoutils/](../cleoutils/)**: Contains code that is re-used among cleosrv
  and cleoc and helpers. The purpose of `tools.go` is to install specific
  versions of tools such as goreleaser.
* **[design/](../design/)**: Contains logos and other graphical resources.
* **[dev_docs/](../dev_docs/)**: Documentation for development. User
  documentation can be found on the website.
* **[e2e_tests/](../e2e_tests/)**: Contains end to end tests that test the
  entire application stack.
* **[frontend/](../frontend/)**: Contains the React frontend that is bundled
  with cleosrv during build. Get started with `src/index.tsx` .
* **[scripts/](../scripts/)**: Contains some helper scripts. The are expected
  to be executed from the top level directory.
* **[website/](../website/)**: Contains the source for the cleodora.org
  website. In the `content` directory you'll find the interesting stuff.
* **[Makefile](../Makefile)**: Run the build, lint the software etc. Execute
  `make` to see the possible targets.
* **[schema.graphql](../schema.graphql)**: The GraphQL schema used by `cleoc`,
  `cleosrv` and the frontend.
