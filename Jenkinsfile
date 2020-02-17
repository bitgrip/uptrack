pipeline {
  agent { node { label 'docker' } }
  options {
    buildDiscarder(logRotator(numToKeepStr: '3'))
    disableConcurrentBuilds()
  }
  environment {
    DOCKER_USER = "${userID()}"
    DOCKER_GROUP = "${groupID('docker')}"
    IMAGE_PUSH_REGISTRY = 'https://docker-push.bitgrip.berlin'
    IMAGE_REPOSITORY = 'bitgrip'
    GIT_HOST = "bitbucket.org"
    GIT_EMAIL = "bolt@bitgrip.berlin"
    GIT_LOGIN = "boltbitgrip"


    EXTENDED_IMAGE_TAG = """${sh(
            returnStdout: true,
            script: 'git describe --tags'
    ).trim()}"""

    BUILD_DATE = """${sh(
            returnStdout: true,
            script: "date -u '+%Y-%m-%dT%H:%M:%SZ'"
    ).trim()}"""

    VCS_REF = """${sh(
            returnStdout: true,
            script: 'git rev-parse --short HEAD'
    ).trim()}"""


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


          def image = docker.build("${IMAGE_REPOSITORY}/uptrack:${EXTENDED_IMAGE_TAG}", "--squash --pull . --build-arg BUILD_DATE=${BUILD_DATE}  VCS_REF=${VCS_REF} --short HEAD'")
          docker.withRegistry("${IMAGE_PUSH_REGISTRY}", 'bg-system') {
            image.push()
          }
        }
      }


    }

  }
}
