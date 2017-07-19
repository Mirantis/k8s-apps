def run(helm_home) {
  stage("Add repo") {
    sh('./helm repo add mirantisworkloads https://mirantisworkloads.storage.googleapis.com/')
  }

  stage("Run tests") {
    try {
      withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
        withEnv([
          'HELM_HOME=' + helm_home,
          'HELM_CMD=' + pwd() + '/helm',
          'KUBECTL_CMD=' + pwd() + '/kubectl',
        ]) {
          sh("set -o pipefail; exec 3>&1; go test -v -timeout 90m -args --charts --image-repo mirantisworkloads --exclude tweepub,tweetics,kibana,logstash --prefix j${env.BUILD_NUMBER}-nightly- 2>&1 3>&- | tee /dev/fd/3 | ./go-junit-report > report.xml 3>&-")
        }
      }
    } finally {
      junit('report.xml')
    }
  }
}
return this;
