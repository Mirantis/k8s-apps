def run(helm_home) {
  withEnv(['HELM_HOME=' + helm_home]) {
    def chart_names = []
    dir("charts") {
      for (chart_dir in findFiles()) {
        chart_names << chart_dir.name
      }
    }

    stage("Dependencies") {
      sh('go run tools/pre-test-local-repos.go')
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
  def image_tag = env.BUILD_NUMBER + '-' + env.GERRIT_CHANGE_NUMBER + '-' + env.GERRIT_PATCHSET_NUMBER
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
          sh('docker build --tag nexus-scc.ng.mirantis.net:5000/' + image_name + ':' + image_tag + ' ' + image_name)
          sh('docker push nexus-scc.ng.mirantis.net:5000/' + image_name + ':' + image_tag)
        }
      }
    }
    parallel stages
  }
  stage("Run tests") {
    def params = []

    for (prefix in ['', 'logCollector.', 'testImage.', 'filebeatImage.', 'spark.', 'zeppelin.']) {
      params << prefix + 'image.repository=nexus-scc.ng.mirantis.net:5000/'
      params << prefix + 'image.tag=' + image_tag
    }

    try {
      withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
        withEnv([
          'HELM_HOME=' + helm_home,
          'HELM_CMD=' + pwd() + '/helm',
          'KUBECTL_CMD=' + pwd() + '/kubectl',
        ]) {
          if (env.GERRIT_CHANGE_NUMBER == '6417') {
              sh("set -o pipefail; exec 3>&1; go test -v -timeout 90m -args --charts --image-repo nexus-scc.ng.mirantis.net:5000/${image_tag}/ --verify-version --exclude tweepub,tweetics,kibana,logstash --prefix j${image_tag}- 2>&1 3>&- | tee /dev/fd/3 | ./go-junit-report > report.xml 3>&-")
          } else {
              sh("set -o pipefail; exec 3>&1; go test -v -timeout 90m -args --charts --image-repo nexus-scc.ng.mirantis.net:5000/${image_tag}/ --exclude tweepub,tweetics,kibana,logstash --prefix j${image_tag}- 2>&1 3>&- | tee /dev/fd/3 | ./go-junit-report > report.xml 3>&-")
          }
        }
      }
    } finally {
      junit('report.xml')
    }
  }
}
return this;
