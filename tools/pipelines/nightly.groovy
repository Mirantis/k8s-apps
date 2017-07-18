def run(helm_home) {
  withEnv(['HELM_HOME=' + helm_home]) {
    def chart_names = []
    dir("charts") {
      for (chart_dir in findFiles()) {
        chart_names << chart_dir.name
      }
    }

    stage("Dependencies") {
      sh('./helm repo add mirantisworkloads https://mirantisworkloads.storage.googleapis.com/')
      for (chart_name in chart_names) {
        sh('./helm dependency update charts/' + chart_name);
      }
    }

    def stages = [:]
    for (_chart_name in chart_names) {
      def chart_name = _chart_name;
      stages[chart_name] = {
        stage(chart_name + " lint") {
          sh('./helm lint charts/' + chart_name);
        }
      }
    }
    parallel stages
  }
  dir("images") {
    def stages = [:]
    def image_names = []
    for (image_dir in findFiles()) {
      image_names << image_dir.name
    }
    for (_image_name in image_names) {
      def image_name = _image_name
      stages[image_name] = {
        stage(image_name + " image build") {
          sh('docker build --no-cache --pull ' + image_name)
        }
      }
    }
    parallel stages
  }
  stage("Run tests") {
    try {
      withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
        withEnv([
          'HELM_HOME=' + helm_home,
          'HELM_CMD=' + pwd() + '/helm',
          'KUBECTL_CMD=' + pwd() + '/kubectl',
        ]) {
          sh('set -o pipefail; exec 3>&1; go test -v -timeout 90m -args --exclude tweepub,tweetics,kibana,logstash --prefix j' + env.BUILD_NUMBER + '-nightly- 2>&1 3>&- | tee /dev/fd/3 | ./go-junit-report > report.xml 3>&-')
        }
      }
    } finally {
      junit('report.xml')
    }
  }
}
return this;
