# Welcome

You have stumbled upon a top secret part of the internet where rocket scientists hang out.

It would be a multi-national security concern if you joined a [Discord channel](https://discord.gg/sWdSk6m) where the rocket scientists answer any questions thrown their way.

Of course we are kidding, Space Cloud welcomes you!

## How to Contribute to Space Cloud? (General Guidelines)

There are multiple ways to contribute to Space Cloud:
- Code contributions.
  - Looking for a starting point? Search for issues with the labels [good first issue](https://github.com/spaceuptech/space-cloud/labels/good%20first%20issue) or [help wanted](https://github.com/spaceuptech/space-cloud/labels/help%20wanted).
- Creating well-described issues with reproducible steps.
- Community contributions (making examples, tutorials, spreading the word, or sending out [this tweet](https://twitter.com/intent/tweet?url=&text=Hey%2C%20I%20just%20found%20a%20cool%20open%20source%20alternative%20for%20Firebase%20and%20Heroku.%20It's%20%40SpaceUpTech's%20%23SpaceCloud.%20Check%20it%20out%3A%20https%3A%2F%2Fgithub.com%2Fspaceuptech%2Fspace-cloud)).
- Contributing to discussions and offering your views and feedback on community calls.

- Every PR should be attached to an issue. If your contribution doesn't align with an existing issue, create a new issue and link it with your PR.
- Space Cloud's [motivation](https://docs.spaceuptech.com/introduction/motivation/) and [design goals](https://docs.spaceuptech.com/introduction/design-goals/) are the driving factors for feature ideas.

> If you are confused about where to start, which issue to pick or how to solve one, we welcome your questions in the [#contributions channel](https://discord.gg/sWdSk6m) on Discord.

## Sign the CLA
### Before your PRs can be merged, you need to sign the [Contributor License Agreement](https://cla-assistant.io/spaceuptech/space-cloud).

## Project Structure

- The core components of Space Cloud are `gateway`, `runner`, `metric-proxy` and `space-cli`.
- The codebase is a monorepo consisting of the core components and their sub modules.

## Branching Strategy and Guidelines

- Space Cloud follows a version and release based branching practice. Active development happens on the version branches starting with the prefix `v`.
- For creating a PR, you should create a fork from the latest version branch. Avoid creating forks from master because a version release might force you to rebase your commits.