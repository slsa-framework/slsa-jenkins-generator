## Plugin Example Usage

### Prepare
This tool requires a payload sent by Github(WebHook).
Add the webhook to Jenkins in Github Repository.
The payload are pass to this tool through environment variable.

For this action, the following plugins must be installed on Jenkins.
```plugins
"Git plugin", "Docker plugin", "Generic Webhook Trigger Plugin" and "GitHub plugin"
```

### slsa-jenkins-generator Plugin Build
```plugins
cd slsa-jenkins-generator/plugin
mvn clean package -DskipTests
```

### slsa-jenkins-generator Plugin Install 
If an administrator manually copies a plugin archive into the plugins directory, it should be named with a .hpi suffix to match the file names used by plugins installed from the update center.

Or
#####
Navigate to the Manage Jenkins > Manage Plugins page in the web UI.
Click on the Advanced tab.
Choose the .hpi file from your system or enter a URL to the archive file under the Deploy Plugin section.
Deploy the plugin file.
Once a plugin file has been uploaded, the Jenkins controller must be manually restarted in order for the changes to take effect.



### Usage
Below is sample to insert into your build or release pipeline.

```plugins
        stage('generate provenance') {
            steps {
                dir("../workspace") {
                    dir("slsa-jenkins-generator") {
                        script {
                            env_path = sh(script: 'pwd', returnStdout: true).trim()
                            sh "printenv > ${env_path}/.env"
                            env_path = "${env_path}"+"/.env"
                        }
                        step([$class: 'ProvenanceGenerator', artifact:artifact_path, output:output_path, envpath:env_path])
                    }
                }
            }
        }
```
More details, see the [example pipeline](.config/pipeline) or [example config_for_Job](.config/config.xml) in this repository.

### Output
See [example provenance](../example/provenance.slsa) in this repository.
