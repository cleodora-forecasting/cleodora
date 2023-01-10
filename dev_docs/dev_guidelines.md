# Dev Guidelines

The following should be prioritized, in particular over feature completeness:

* Long-term maintainability, including a high degree of automation and tests
* Easy setup and excellent 'Getting Started' documentation for new contributors
* Simplicity of use
* Personal use and privacy

When in doubt, new features will not be added. It's better to do a few things,
but do them very well.


## Git

* A commit should do one thing only e.g. fix a bug, fix whitespace in one or
  multiple files, implement (part of) a feature
  * Reason: It's easier to review, revert and cherry-pick
* A commit should not break the build, tests etc. to always leave the system in
  a runnable state
  * Reason: Don't hinder other developers and make tracking down issues with
    the help of `git bisect` easier
* Small commits are better
  * Reason: It's easier to review
* Good commit messages are important. They should include a summary of the WHAT
  and an explanation of the WHY. [More
  details](https://gist.github.com/robertpainsi/b632364184e70900af4ab688decf6f53)
  * Reason: The Git log is the history of the repository and should be
    understandable e.g. to write the changelog. The WHY is required to be able
    to understand and question past decisions in the future (Chesterton's
    fence).
* Long running secondary branches or huge pull requests are undesirable.
  Instead merge frequently.
  * Reason: Allow other developers to build on top of your work, avoiding
    gigantic merge conflicts.
* The Git repository is a [monorepo](https://en.wikipedia.org/wiki/Monorepo)
  and should stay that way. Creating more repositories is generally
  undesirable.
  * Reason: It is easier to handle dependencies between parts of the project
    (e.g. documentation and code, API and clients) as well as easier to do
    cross-cutting modifications (e.g. rename forecast 'summary' to 'title').


## Tests

* Prefer testing functionality over testing implementation. This means leaning
  towards _integration tests_.
  * Reason: It makes it possible to refactor code without touching the tests.
    The tests serve as documentation. The functionality is the important part
    you actually want to be tested, not the implementation.
* Mock as little as possible, but as much as necessary to ensure the tests are
  fast enough, stable and isolated from each other.
  * Reason: Mocks tie you to the implementation and prevent you from treating
    the module as a blackbox (e.g. you need to know which mock is called in
    what order, how many times etc. when in fact you should only care about the
    final result). Extensive mocking can also force you to re-implement parts
    of your application inside the mocks.

Some useful resources:

* [TDD, Where Did It All Go Wrong (Ian
  Cooper)](https://www.youtube.com/watch?v=EZ05e7EMOLM)
* [TDD Revisited (Ian Cooper)](https://www.youtube.com/watch?v=Jh50o8-wWCA)
* https://kentcdodds.com/blog/static-vs-unit-vs-integration-vs-e2e-tests
* https://kentcdodds.com/blog/common-mistakes-with-react-testing-library
