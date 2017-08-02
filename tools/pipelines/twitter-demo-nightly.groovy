def run(helm_home) {
  def namespace = "twitter-nightly-${env.BUILD_NUMBER}"
  try {
    stage("Run tests") {
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
          "TS_NAME=${namespace}"
        ]) {
          ansiColor("xterm") {
            sh("./scenarios/twitter-stats/twitter-stats.sh up test down")
          }
        }
      }
    }
  } finally {
    stage("Cleanup") {
      timeout(15) {
        withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
          def ok = false
          while (!ok) {
            sleep(5)
            ok = sh(script: "./kubectl delete ns --ignore-not-found ${namespace}", returnStatus: true) == 0
          }
        }
      }
    }
  }
}

return this;
