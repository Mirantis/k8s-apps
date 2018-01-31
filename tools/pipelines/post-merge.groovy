def run(helm_home, namespace_prefix, kubernetes_domain, buildId, docker_repository) {
  stage("Dependencies") {
    sh('go run tools/pre-test-local-repos.go')
    sh("go get github.com/kubernetes/apimachinery/pkg/util/yaml")
    sh("cp -r charts tmp-charts")
  }
  stage("Run tests") {
    try {
      withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
        withEnv([
          'HELM_HOME=' + helm_home,
          'HELM_CMD=' + pwd() + '/helm',
          'KUBECTL_CMD=' + pwd() + '/kubectl',
        ]) {
          sh("set -o pipefail; exec 3>&1; go test -v -timeout 90m -args --images --charts --kubernetes-domain ${kubernetes_domain} --image-repo ${docker_repository}/${buildId} --build-images-opts='--no-cache' --exclude tweepub,tweetics,tweeviz,tweeviz_api,tweeviz_ui,kibana,logstash,bus-floating-data,monocular --prefix ${namespace_prefix}- 2>&1 3>&- | tee /dev/fd/3 | ./go-junit-report > report.xml 3>&-")
        }
      }
    } finally {
      junit('report.xml')
    }
  }

  stage("Push images") {
    sh("go test -v -timeout 60m -args --charts=false --push --image-repo ${docker_repository}/${buildId}")
  }

  stage("Add chart repo") {
    withEnv(["HELM_HOME=${helm_home}"]) {
      sh('./helm repo add mirantisworkloads https://mirantisworkloads.storage.googleapis.com/')
    }
  }

  stage("Build packages") {
    withEnv([
      'HELM_HOME=' + helm_home,
      'HELM_CMD=' + pwd() + '/helm',
    ]) {
      sh("rm -rf charts && mv tmp-charts charts")
      sh("./tools/build-packages.sh")
    }
  }

  stage("Push packages") {
    withCredentials([file(credentialsId: "gcloud-mirantisworkloads", variable: "GCLOUD_KEYPATH"),
                     string(credentialsId: "gcloud-mirantisworkloads-project", variable: "GCLOUD_PROJECT")]) {
      sh("./tools/push-packages.sh")
    }
  }
}

return this;
