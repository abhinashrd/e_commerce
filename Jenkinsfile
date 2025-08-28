pipeline {
    agent any

    environment {
        DOCKER_BUILDKIT = '1'
        COMPOSE_DOCKER_CLI_BUILD = '1'
    }

    options {
        ansiColor('xterm')
        timestamps()
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build images') {
            steps {
                ansiColor('xterm') {
                    sh '''
                      echo "🐳 Building Docker images..."
                      docker compose -f docker-compose.yml build --no-cache
                    '''
                }
            }
        }

        stage('Unit tests') {
            steps {
                ansiColor('xterm') {
                    sh '''
                      echo "🔍 Running unit tests in each service..."

                      for svc in usersvc productsvc ordersvc; do
                        echo "➡️ Testing $svc ..."
                        docker run --rm -v $PWD/$svc:/app -w /app golang:1.22 go test ./...
                      done
                    '''
                }
            }
        }

        stage('Start stack') {
            steps {
                ansiColor('xterm') {
                    sh '''
                      echo "🚀 Starting full stack with docker-compose..."
                      docker compose -f docker-compose.yml up -d
                    '''
                }
            }
        }

        stage('Wait for health') {
            steps {
                ansiColor('xterm') {
                    sh '''
                      echo "⏳ Waiting for services to become healthy..."
                      sleep 15
                    '''
                }
            }
        }

        stage('Integration smoke test') {
            steps {
                ansiColor('xterm') {
                    sh '''
                      echo "🧪 Running smoke tests..."

                      curl -f http://localhost:8080/healthz || exit 1
                      curl -f http://localhost:8082/healthz || exit 1
                      curl -f http://localhost:8083/healthz || exit 1

                      echo "✅ All services passed smoke test"
                    '''
                }
            }
        }

        stage('Push images (optional)') {
            when {
                expression { return env.BRANCH_NAME == 'main' }
            }
            steps {
                ansiColor('xterm') {
                    sh '''
                      echo "📦 Pushing images to Docker Hub (if configured)..."
                      # Example:
                      # docker tag go-ecommerce-usersvc your-dockerhub-user/go-ecommerce-usersvc:latest
                      # docker push your-dockerhub-user/go-ecommerce-usersvc:latest
                    '''
                }
            }
        }
    }

    post {
        always {
            ansiColor('xterm') {
                sh 'docker compose -f docker-compose.yml down -v || true'
            }
            echo '✅ Cleanup finished'
        }
        failure {
            echo '❌ Pipeline failed'
        }
        success {
            echo '🎉 Pipeline succeeded'
        }
    }
}