# Contributing to the MultiApps CF CLI Plugin

## Did you find a bug?
* Check if the bug has already been reported and has an open [Issue](https://github.com/cloudfoundry/multiapps-cli-plugin/issues).

* If there is none, create one by using the provided [Issue Template](https://github.com/cloudfoundry/multiapps-cli-plugin/issues/new/choose) for bugs.

* Try to be as detailed as possible when describing the bug. Every bit of information helps!

## Do you have a question or need support?
If you need any support or have any questions regarding the project, you can drop us a message on [Slack](https://cloudfoundry.slack.com/?redir=%2Fmessages%2Fmultiapps-dev) or open an [Issue](https://github.com/cloudfoundry/multiapps-cli-plugin/issues) and we shall get back to you.

## Do you want to contribute to the code base?

### Starter GitHub Issues
If you are looking for what you can contribute to the project, check the GitHub Issues labeled as [Good First Issue](https://github.com/cloudfoundry/multiapps-cli-plugin/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) to find items that are marked as more beginner friendly.

### Fork the project
* To develop your contribution to the project, first [fork](https://help.github.com/articles/fork-a-repo/) this repository in your own github account.

* To clone the project into your Go workspace check the [Cloning the repository](https://github.com/cloudfoundry/multiapps-cli-plugin#cloning-the-repository) section.

* When developing make sure to keep your fork up to date with the origin's master branch or the release branch you want to contribute a fix to.

### How to build, develop and install?
* To build a new version of the plugin follow the [Development](https://github.com/cloudfoundry/multiapps-cli-plugin#development) instructions.

* If you have added new dependencies into the CF plugin make sure to update them as described in the [Adding Dependencies](https://github.com/cloudfoundry/multiapps-cli-plugin#adding-dependency-into-the-multiapps-cli-plugin).

* To install the plugin follow the [Installation](https://github.com/cloudfoundry/multiapps-cli-plugin#installation) instructions.

### Testing
* Running the tests is done with the [ginkgo](https://github.com/onsi/ginkgo) framework. Once you have it installed just execute `ginkgo -r` from the project's root directory and it will run all the tests.

* If you are developing new functionality make sure to add tests covering the new scenarios where applicable!

* The [spring-music](https://github.com/nvvalchev/spring-music) contains a handy sample MTA archive to test your MultiApps CF CLI Plugin against the MultiApps Controller.

### Formatting
Having the same style of formatting across the project helps a lot with readability. To format the project's source code run the following command from the root directory of the project:
```
gofmt -w cli clients commands testutil ui util
```

## Creating a pull request
When creating a pull request please use the provided template. Don't forget to link the [Issue](https://github.com/cloudfoundry/multiapps-cli-plugin/issues) if there is one related to your pull request!
