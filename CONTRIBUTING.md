# Contributing

## Introduction

Please note we have a code of conduct, please follow it in all your interactions with the project.

## Asking For Help

If you need help with _ShellExpand_, feel free to open an issue with your question. We probably won't have time to answer it, but maybe someone else in the community will.

If you want some of our time to help you use _ShellExpand_, feel free to contact us as `hello` @ `ganbarodigital.com`. We'll need to charge you for that time, and we'll agree any fee up front with you before any work starts.

## Reporting A Bug

### Is It A Bug In Your Code?

Most of the features in _ShellExpand_ have to call your [expansion callbacks](README.md#expansion-callbacks). Is the bug caused by a problem in one of your callbacks?

It helps us enormously if you spend a few minutes ruling that out before opening new bug reports on GitHub.

### Existing Issue?

Has someone already reported the same bug? Please take a look through our open issues, just in case they have. If there's an issue already open, please add your bug report to that issue. It helps everyone a lot.

If can't find an open issue that describes your bug, please do open a new issue.

### Creating A New Issue

When you create a new bug report, please include the following:

1. UNIX shell script to demonstate the correct behaviour - the smaller, the better
2. sample Golang code that demonstrates the bug - again, the smaller, the better please!

You only need to include a UNIX shell script in your bug report if _ShellExpand_ is returning the wrong output. Definitely not needed if _ShellExpand_ is crashing or getting stuck in an infinite loop somewhere in the parser!

If you can include a unit test (preferably one that would go in [expand_test.go](expand_test.go)) to demonstrate the problem, that's guaranteed to catch our attention and bump your bug report up to the front of the line.

## Bug Fixes and Feature Requests

### Before You Submit A Pull Request

Please take a look through our open issues to see if someone else has already reported the problem or already raised a request for your feature.

* If they have, please comment on that issue. It helps everyone if we can avoid scattering related information across several open issues.
* If they haven't, please do create a new issue.

It saves everyone a lot of time if we discuss your proposal together before we start looking at code.

### Creating The Feature Request

When you create the feature request, please include the following:

1. UNIX shell script to demonstrate the new feature - the smaller, the better
2. a description of the new feature you're proposing
3. a pointer to where we'll find that feature described in `bash`'s man page
4. a list of any backwards-compatibility breaks

We're aiming for as close to 100% compatibility with the behaviour in the latest version of `bash`. The only exception we'll make is around error handling and reporting.

* If you open a feature request asking us to break compatibility with how `bash` does string expansion, we'll almost certainly reject it.
* I'm sorry, the same goes for feature requests asking us to add string expansion features that `bash` doesn't have.
* And the same goes for feature requests that fall outside the scope of this package.

We'll discuss your feature request with you, and agree the test cases that your pull request will need to satisfy to get accepted.

### Raising The Pull Request

All pull requests need to be against the latest `develop` branch, please, __not__ `master`.

Please make sure your pull request includes:

1. Unit tests for your feature.
2. Your feature itself.
3. Updates to [README.md](README.md) to describe your feature.

We have 100% code coverage, and we can't accept any pull requests that contain untested code. Think of 100% code coverage as the barest minimum; 100% feature coverage is our prefered goal to aim for!

We don't care all that much about your commit history in the pull request, as long as the commit messages aren't offensive. Whether you're working on this in your own time or not, please treat this as a professional working environment. If it wouldn't be tolerated in the work place, it's not appropriate for this project either.

Please keep each pull request down to a single feature at a time. It's quicker for us to test, review and accept pull requests if they're small and easy to understand.

If you've got a whole pile of features that you need us to review and accept to help you in your day job, the best way to help us help you is to book some time for us to work with you on this. We'll need to charge you for that time, and any fee will be agreed up front before the work starts.

## Code of Conduct

### Our Pledge

In the interest of fostering an open and welcoming environment, we as
contributors and maintainers pledge to making participation in our project and
our community a harassment-free experience for everyone, regardless of age, body
size, disability, ethnicity, gender identity and expression, level of experience,
nationality, personal appearance, race, religion, or sexual identity and
orientation.

### Our Standards

Examples of behavior that contributes to creating a positive environment
include:

* Using welcoming and inclusive language
* Being respectful of differing viewpoints and experiences
* Gracefully accepting constructive criticism
* Focusing on what is best for the community
* Showing empathy towards other community members

Examples of unacceptable behavior by participants include:

* The use of sexualized language or imagery and unwelcome sexual attention or
advances
* Trolling, insulting/derogatory comments, and personal or political attacks
* Public or private harassment
* Publishing others' private information, such as a physical or electronic
  address, without explicit permission
* Other conduct which could reasonably be considered inappropriate in a
  professional setting

### Our Responsibilities

Project maintainers are responsible for clarifying the standards of acceptable
behavior and are expected to take appropriate and fair corrective action in
response to any instances of unacceptable behavior.

Project maintainers have the right and responsibility to remove, edit, or
reject comments, commits, code, wiki edits, issues, and other contributions
that are not aligned to this Code of Conduct, or to ban temporarily or
permanently any contributor for other behaviors that they deem inappropriate,
threatening, offensive, or harmful.

### Scope

This Code of Conduct applies both within project spaces and in public spaces
when an individual is representing the project or its community. Examples of
representing a project or community include using an official project e-mail
address, posting via an official social media account, or acting as an appointed
representative at an online or offline event. Representation of a project may be
further defined and clarified by project maintainers.

### Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be
reported by contacting the project team at `go_shellexpand` @ `ganbarodigital.com`. All
complaints will be reviewed and investigated and will result in a response that
is deemed necessary and appropriate to the circumstances. The project team is
obligated to maintain confidentiality with regard to the reporter of an incident.
Further details of specific enforcement policies may be posted separately.

Project maintainers who do not follow or enforce the Code of Conduct in good
faith may face temporary or permanent repercussions as determined by other
members of the project's leadership.

### Attribution

This Code of Conduct is adapted from the [Contributor Covenant][homepage], version 1.4,
available at [http://contributor-covenant.org/version/1/4][version]. Based on [a template](https://gist.githubusercontent.com/PurpleBooth/b24679402957c63ec426/raw/5c4f62c1e50c1e6654e76e873aba3df2b0cdeea2/Good-CONTRIBUTING.md-template.md) published by [PurpleBooth](https://gist.github.com/PurpleBooth)

[homepage]: http://contributor-covenant.org
[version]: http://contributor-covenant.org/version/1/4/