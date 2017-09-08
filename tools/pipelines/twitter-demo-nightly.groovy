def run(helm_home, namespace, kubernetes_domain) {
  stage("Run tests") {
    sh("find charts/ -name values.yaml | xargs sed -i -e 's/kubernetes_domain: cluster.local/kubernetes_domain: ${kubernetes_domain}/g'")
    withCredentials([
      file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG'),
      string(credentialsId: "twitter-demo-ts-app-key", variable: "TS_APP_KEY"),
      string(credentialsId: "twitter-demo-ts-app-secret", variable: "TS_APP_SECRET"),
      string(credentialsId: "twitter-demo-ts-token-key", variable: "TS_TOKEN_KEY"),
      string(credentialsId: "twitter-demo-ts-token-secret", variable: "TS_TOKEN_SECRET")
    ]) {
      withEnv([
        "HELM_HOME=" + helm_home,
        "TS_HELM_CMD=" + pwd() + "/helm",
        "TS_KUBECTL_CMD=" + pwd() + "/kubectl",
        "TS_NAME=${namespace}",
        "TS_USE_INTERNAL_IP=yes"
      ]) {
        ansiColor("xterm") {
          sh("./scenarios/twitter-stats/twitter-stats.sh up test down")
        }
      }
    }
  }
}

return this;
