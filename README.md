# MyDocker T17 scratch canary

This public repository is dedicated to the MyDocker GitHub scratch-repository
qualification. It is not a production repository and contains no product
source, provider credentials, repository secrets, or reusable customer state.

The only workflow is manually dispatched by the protected MyDocker provider
runner. It performs no checkout, executes no pull-request-controlled code, and
uses GitHub's short-lived OIDC endpoint only to validate the fixed
GitHub-hosted workflow identity. The protected MyDocker driver owns all
reconciliation, replay fencing, and terminal cleanup.
