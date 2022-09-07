pipeline {
  agent {
    label 'ubuntu_docker_label'
  }
  tools {
    go "Go 1.13"
  }
  options {
    checkoutToSubdirectory('src/github.com/Infoblox-CTO/heka-ui')
  }
  environment {
    GOPATH = "$WORKSPACE"
    DIRECTORY = "src/github.com/Infoblox-CTO/heka-ui"
  }
  stages {
    stage("Test") {
      steps {
        sh "cd $DIRECTORY && make test"
      }
    }
    stage("Build") {
      steps {
        withDockerRegistry([credentialsId: "dockerhub-bloxcicd", url: ""]) {
          sh "cd $DIRECTORY && make docker && make push"
        }
      }
    }
    stage("Push merge") {
      when {
        not { changeRequest() }
        not { buildingTag() }
      }
      steps {
        withDockerRegistry([credentialsId: "dockerhub-bloxcicd", url: ""]) {
          sh "cd $DIRECTORY && make push-latest"
          sh "cd $DIRECTORY && make push"
        }
      }
    }
    stage("Push Release/Tag") {
      when {
        buildingTag()
      }
      steps {
        withDockerRegistry([credentialsId: "dockerhub-bloxcicd", url: ""]) {
          sh 'cd $DIRECTORY && make push IMAGE_VERSION=${TAG_NAME}'
        }
      }
    }

  }
  post {
    success {
      finalizeBuild()
    }
    always {
      sh "cd $DIRECTORY && make clean"
      cleanWs()
    }
  }
}
