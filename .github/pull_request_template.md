<!--  All Submissions -->


## Community Note
<!-- Please leave the community note as is. -->
* Please vote on this PR by adding a :thumbsup: [reaction](https://blog.github.com/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/) to the original PR to help the community and maintainers prioritize for review
* Please do not leave comments along the lines of "+1", "me too" or "any updates", they generate extra noise for PR followers and do not help prioritize for review


## Description

<!-- Please include a description below with the reason for the PR, what it is doing, what it is trying to accomplish, and anything relevant for a reviewer to know. 

If this is a breaking change for users please detail how it cannot be avoided and why it should be made in a minor version of the provider -->


## Changes to existing Resource / Data Source

- [ ] I have added an explanation of what my changes do and why I'd like you to include them (This may be covered by linking to an issue above, but may benefit from additional explanation).
- [ ] I have written new tests for my resource or datasource changes & updated any relevant documentation.
- [ ] I have successfully run tests with my changes locally. If not, please provide details on testing challenges that prevented you running the tests.
- [ ] (For changes that include a **state migration only**). I have manually tested the migration path between relevant versions of the provider.


## Testing 

- [ ] My submission includes Test coverage as described in the [Contribution Guide](../blob/main/contributing/topics/guide-new-resource.md) and the tests pass. (if this is not possible for any reason, please include details of why you did or could not add test coverage)

<!-- Please include testing logs or evidence here or an explanation on why no testing evidence can be provided. 

For state migrations please test the changes locally and provide details here, such as the versions involved in testing the migration path. For further details on testing state migration changes please see our guide on [state migrations](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/contributing/topics/guide-state-migrations.md#testing) in the contributor documentation. -->


## Change Log

Below please provide what should go into the changelog (if anything) conforming to the [Changelog Format documented here](../blob/main/contributing/topics/maintainer-changelog.md).

<!-- Replace the changelog example below with your entry. One resource per line. -->

* `ad_resource` - support for the `thing1` property [GH-00000]


<!-- What type of PR is this? -->
This is a (please select all that apply):

- [ ] Bug Fix
- [ ] New Feature (ie adding a service, resource, or data source)
- [ ] Enhancement
- [ ] Breaking Change


## Related Issue(s)
Fixes #0000

<!-- heimdall_github_prtemplate:grc-pci_dss-2024-01-05 -->

## Rollback Plan

If a change needs to be reverted, we will publish an updated version of the provider.

## Changes to Security Controls

Are there any changes to security controls (access controls, encryption, logging) in this pull request? If so, explain.

> [!NOTE] 
> If this PR changes meaningfully during the course of review please update the title and description as required.
