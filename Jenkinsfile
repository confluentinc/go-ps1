def config = jobConfig {
    realJobPrefixes = ['cli']
    timeoutHours = 1
}

def job = {
    if (config.isPrJob) {
        configureGitSSH("github/confluent_jenkins", "private_key")

        stage('Setup Go and Dependencies') {
            withVaultEnv([["github/confluent_jenkins", "user", "GIT_USER"],
                ["github/confluent_jenkins", "access_token", "GIT_TOKEN"],
                ["sonatype/confluent", "user", "SONATYPE_OSSRH_USER"],
                ["sonatype/confluent", "password", "SONATYPE_OSSRH_PASSWORD"]]){
                withEnv(["GIT_CREDENTIAL=${env.GIT_USER}:${env.GIT_TOKEN}", "GIT_USER=${env.GIT_USER}", "GIT_TOKEN=${env.GIT_TOKEN}"]) {
                    sh '''#!/bin/bash -i
                        export GOVER=1.17.6
                        wget "https://golang.org/dl/go${GOVER}.linux-amd64.tar.gz" --quiet --output-document go${GOVER}.tar.gz
                        tar -C $(pwd)/.. -xzf go${GOVER}.tar.gz
                        echo "export GOROOT=$(pwd)/../go" >> ~/.bashrc
                        echo "export GOPATH=$(pwd)/../go/path" >> ~/.bashrc
                        echo "export GOBIN=$(pwd)/../go/bin" >> ~/.bashrc
                        echo "export modulePath=$(pwd)/../go/src/github.com/confluentinc/go-ps1" >> ~/.bashrc
                        source ~/.bashrc
                        echo "GOROOT IS ${GOROOT}\n"
                        mkdir -p $GOPATH/bin
                        mkdir -p $GOROOT/bin
                        echo "export PATH=${GOPATH}/bin:${GOROOT}/bin:${GOBIN}:$PATH" >> ~/.bashrc
                        source ~/.bashrc
                        cat ~/.bashrc
                        echo "machine github.com\n\tlogin $GIT_USER\n\tpassword $GIT_TOKEN" > ~/.netrc
                        make jenkins-deps || exit 1
                    '''
                }
            }
        }

        stage('Build, Test, and Release') {
            withVaultEnv([["github/confluent_jenkins", "user", "GIT_USER"],
                ["github/confluent_jenkins", "access_token", "GIT_TOKEN"],
                ["sonatype/confluent", "user", "SONATYPE_OSSRH_USER"],
                ["sonatype/confluent", "password", "SONATYPE_OSSRH_PASSWORD"]]){
                withEnv(["GIT_CREDENTIAL=${env.GIT_USER}:${env.GIT_TOKEN}", "GIT_USER=${env.GIT_USER}", "GIT_TOKEN=${env.GIT_TOKEN}"]) {
                    sh '''#!/bin/bash -i
                        source ~/.bashrc
                        make deps ARGS=--vendor-only || exit 1
                        make test || exit 1
                        make release-ci || exit 1
                    '''
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
