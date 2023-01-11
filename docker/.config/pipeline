def artifact_path
def artifact_name

pipeline {
    agent any
    tools {
        git "git"
        go 'go 1.18'

    }

    environment {
        CREDENTIAL_ID = 'GIT_ACCOUNT'
        Repository_Generator = 'https://github.com/Samsung/slsa-jenkins-generator'
    }

    stages {
        stage('build target ') {
            steps {
                script {
                    branch = "${env.payload_ref}"
                    branch = "${branch}".split("refs/heads/")[1]
                    git branch: "${branch}", credentialsId: "$CREDENTIAL_ID", url: "$env.payload_repository_clone_url"

                     //TODO: replace with real build command
                    sh 'echo "hello world" > output.txt'
                    artifact_path = sh(script: 'pwd', returnStdout: true).trim()
                    artifact_name = "output.txt"
                }
            }
        }

        stage('generate provenance') {
            steps {
                dir("../workspace") {
                    sh 'rm -rf slsa-jenkins-generator && mkdir slsa-jenkins-generator'

                    dir("slsa-jenkins-generator") {
                        git branch: "main", credentialsId: "$CREDENTIAL_ID", url: "$Repository_Generator"

                        // 'run generator via docker'
                        echo 'run generator via docker'
                        sh "docker build . -t scia:slsa-generator"
                        sh "printenv > ./envlist && docker run --env-file ./envlist -v \"${artifact_path}\":\"/artifacts\" scia:slsa-generator -a artifacts/${artifact_name} -o artifacts"
                    }
                }
            }
        }
    }
}
