## Docker Example Usage
### Prepare
This tool requires a payload sent by Github(WebHook).
Add the webhook to Jenkins in Github Repository.
The payload are pass to this tool through environment variable.

For this action, the following plugins must be installed on Jenkins.
```plugins
"Git plugin", "Docker plugin", "Generic Webhook Trigger Plugin", "GitHub plugin" and "Go Plugin"
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
See [example provenance](../example/provenance.slsa) in this repository.
