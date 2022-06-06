pipeline {
    agent any
    stages {
        stage('Setup') {
            steps {
                sh '''
                    make install-vault
                    . mk-include/bin/vault-setup
                    . vault-sem-get-secret semaphore-secrets-global
                    . vault-sem-get-secret artifactory-docker-helm
                    . vault-sem-get-secret testbreak-reporting
                    . vault-sem-get-secret ssh_id_rsa
                    . vault-sem-get-secret netrc
                    . vault-sem-get-secret ssh_config
                    . vault-sem-get-secret gitconfig
                    chmod 400 ~/.ssh/id_rsa
                    make init-ci
                    export "GOPATH=$(go env GOPATH)"
                    export "PATH=${GOPATH}/bin:${PATH}"
                    git config --global url."git@github.com:".insteadOf "https://github.com/"
                '''
                sh '''
                    make deps ARGS=--vendor-only
                    make test
                    make release-ci
                '''
            }
        }
    }
    post {
        always {
            make testbreak-after
        }
    }
}
