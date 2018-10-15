# Contributing to the MultiApps CF CLI Plugin

## Did you find a bug?
1. Check if the bug has already been reported and has an open [Issue](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/issues).

2. If there is none, create one by using the provided [Issue Template](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/issues/new/choose) for bugs.

3. Try to be as detailed as possible when describing the bug. Every bit of information helps!

## Do you have a question or need support?
If you need any support or have any questions regarding the project, you can drop us a message on [Slack](https://cloudfoundry.slack.com/?redir=%2Fmessages%2Fmultiapps-dev) or open an [Issue](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/issues) and we shall get back to you.

## Do you want to contribute to the code base?

### Fork the project
1. To develop your contribution to the project, first [fork](https://help.github.com/articles/fork-a-repo/) this repository in your own github account. 

2. When developing make sure to keep your fork up to date with the origin's master branch or the release branch you want to contribute a fix to.

### How to build, develop and install?
1. To build a new version of the plugin follow the [Development](https://github.com/cloudfoundry-incubator//multiapps-cli-plugin#development) instructions.

2. If you have added new dependencies into the CF plugin make sure to update them as described in the [Adding Dependencies](https://github.com/cloudfoundry-incubator//multiapps-cli-plugin#adding-dependency-into-the-multiapps-cli-plugin).

3. To install the plugin follow the [Installation](https://github.com/cloudfoundry-incubator//multiapps-cli-plugin#installation) instructions.

### Testing
1. Running the tests is done with the [ginkgo](https://github.com/onsi/ginkgo) framework. Once you have it installed just execute `ginkgo -r` from the project's root directory and it will run all the tests.

2. If you are developing new functionality make sure to add tests covering the new scenarios where applicable!

3. The [spring-music](https://github.com/nvvalchev/spring-music) contains a handy sample MTA archive to test your MultiApps CF CLI Plugin against the MultiApps Controller.

### Formatting
Having the same style of formatting across the project helps a lot with readability. To format the project's source code run the following command from the root directory of the project:
```
gofmt -w cli clients commands testutil ui util
```

## Creating a pull request
When creating a pull request please use the provided template. Don't forget to link the [Issue](https://github.com/cloudfoundry-incubator/multiapps-cli-plugin/issues) if there is one related to your pull request!
