# slsa-jenkins-generator

A proof-of-concept SLSA provenance generator for Jenkins.

## Background

[SLSA](https://github.com/slsa-framework/slsa) is a framework intended to codify
and promote secure software supply-chain practices. SLSA helps trace software
artifacts (e.g. binaries) back to the build and source control systems that
produced them using in-toto's
[Attestation](https://github.com/in-toto/attestation/blob/main/spec/README.md)
metadata format.

## Description

This proof-of-concept Jenkins demonstrates an initial SLSA integration
conformant with SLSA Level 1. This provenance can be uploaded to the native
artifact store or to any other artifact repository.

While there are no integrity guarantees on the produced provenance at L1,
publishing artifact provenance in a common format opens up opportunities for
automated analysis and auditing. Additionally, moving build definitions into
source control and onto well-supported, secure build systems represents a marked
improvement from the ecosystem's current state.

## Example Usage

### Prepare
This tool requires a payload sent by Github(WebHook).
Add the webhook to Jenkins in Github Repository.
The payload are pass to this tool through environment variable.

For this action, the following plugins must be installed on Jenkins.
```plugins
"Git", "Docker plugin", "Generic Webhook Trigger Plugin" and "GitHub"
```

### Usage
Below is sample to insert into your build or release pipeline.

```plugins
        stage('generate provenance') {
            steps {
                dir("../workspace") {
                    sh 'rm -rf slsa-jenkins-generator && mkdir slsa-jenkins-generator'

                    dir("slsa-jenkins-generator") {
                        git branch: "main", credentialsId: "$CREDENTIAL_ID", url: "$Repository_Generator"
                        sh "docker build . -t scia:slsa-generator"
                        sh "printenv > ./envlist && docker run --env-file ./envlist -v \"${artifact_path}\":\"/artifacts\" scia:slsa-generator -a artifacts/${artifact_name} -o artifacts"
                    }
                }
            }
        }
```
More details, see the [example pipeline](.config/pipeline) or [example config_for_Job](.config/config.xml) in this repository.

### Output
See [example provenance](example/provenance.slsa) in this repository.


