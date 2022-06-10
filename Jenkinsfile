def config = jobConfig {
    realJobPrefixes = ['cli']
    timeoutHours = 1
}

def job = {
    if (config.isPrJob) {
        configureGitSSH("github/confluent_jenkins", "private_key")
        def mavenSettingsFile = "/home/jenkins/.m2/settings.xml"
        withMavenSettings("maven/jenkins_maven_global_settings", "settings", "MAVEN_GLOBAL_SETTINGS_FILE", mavenSettingsFile) {

            stage('Setup Go and Dependencies, and Run Tests') {
                withVaultEnv([["github/confluent_jenkins", "user", "GIT_USER"],
                    ["github/confluent_jenkins", "access_token", "GIT_TOKEN"],
                    ["sonatype/confluent", "user", "SONATYPE_OSSRH_USER"],
                    ["sonatype/confluent", "password", "SONATYPE_OSSRH_PASSWORD"]]){
                    withEnv(["GIT_CREDENTIAL=${env.GIT_USER}:${env.GIT_TOKEN}", "GIT_USER=${env.GIT_USER}", "GIT_TOKEN=${env.GIT_TOKEN}"]) {
                        withVaultFile([["gradle/gradle_properties_maven", "gradle_properties_file",
                            "gradle.properties", "GRADLE_PROPERTIES_FILE"]]) {
                            sh '''#!/usr/bin/env bash
                                export HASH=$(git rev-parse --short=7 HEAD)
                                wget "https://golang.org/dl/go1.17.6.linux-amd64.tar.gz" --quiet --output-document go1.17.6.tar.gz
                                tar -C $(pwd)/.. -xzf go1.17.6.tar.gz
                                export GOROOT=$(pwd)/../go
                                export GOPATH=$(pwd)/../go/path
                                export GOBIN=$(pwd)/../go/bin
                                export modulePath=$(pwd)/../go/src/github.com/confluentinc/go-ps1
                                mkdir -p $GOPATH/bin
                                mkdir -p $GOROOT/bin
                                export PATH=$GOPATH/bin:$GOROOT/bin:$GOBIN:$PATH
                                echo "machine github.com\n\tlogin $GIT_USER\n\tpassword $GIT_TOKEN" > ~/.netrc
                                make jenkins-deps || exit 1
                                make deps ARGS=--vendor-only || exit 1
                                make test || exit 1
                                make release-ci || exit 1
                            '''
                        }
                    }
                }
            }
        }
    }
}


def post = {
        stage("Cleanup") {
            sh '''
                make testbreak-after
            '''
        }
}

runJob config, job, post
