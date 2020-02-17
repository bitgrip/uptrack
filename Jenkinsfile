pipeline {
  agent { node { label 'docker' } }
  options {
    buildDiscarder(logRotator(numToKeepStr: '3'))
    disableConcurrentBuilds()
  }
  environment {
    DOCKER_USER = "${userID()}"
    DOCKER_GROUP = "${groupID('docker')}"
    DOCKER_VERSION_SUFFIX = '-0.1.5'
    IMAGE_PUSH_REGISTRY = 'https://docker-push.bitgrip.berlin'
    IMAGE_REPOSITORY = 'bitgrip'
    GIT_HOST = "bitbucket.org"
    GIT_EMAIL = "bolt@bitgrip.berlin"
    GIT_LOGIN = "boltbitgrip"
  }
  stages {
    stage('Build and Push Uptrack Image') {
      agent {
        docker {
          label 'docker'
          image 'docker:18.05.0-ce'
          args '-u $DOCKER_USER:$DOCKER_GROUP -v /var/run/docker.sock:/var/run/docker.sock'
          reuseNode true
          //here we expose docker socket to container. Now we can build docker images in the same way as on host machine where docker daemon is installed'
        }
      }
      steps {
        script {
          env.EXTENDED_IMAGE_TAG = "git describe --tags"
          def image = docker.build("${IMAGE_REPOSITORY}/uptrack:${EXTENDED_IMAGE_TAG}", "--squash --pull . --build-arg BUILD_DATE=\$(date -u +'%Y-%m-%dT%H:%M:%SZ')  VCS_REF='git rev-parse --short HEAD'")
          docker.withRegistry("${IMAGE_PUSH_REGISTRY}", 'bg-system') {
            image.push()
          }
        }
      }


    }

  }
}
