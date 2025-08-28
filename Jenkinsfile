pipeline {
    agent any

    options {
        timestamps()
        ansiColor('xterm')
    }

    environment {
        DOCKER_COMPOSE_FILE = 'docker-compose.yml'
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build images') {
            steps {
                bat '''
                    echo Building Docker images...
                    docker compose -f %DOCKER_COMPOSE_FILE% build --no-cache
                '''
            }
        }

        stage('Unit tests') {
            steps {
                bat '''
                    echo Running Go unit tests...

                    docker run --rm -v "%cd%\\usersvc:/app" -w /app golang:1.22 go test ./...
                    docker run --rm -v "%cd%\\productsvc:/app" -w /app golang:1.22 go test ./...
                    docker run --rm -v "%cd%\\ordersvc:/app" -w /app golang:1.22 go test ./...
                '''
            }
        }

        stage('Start stack') {
            steps {
                bat '''
                    echo Starting full stack...
                    docker compose -f %DOCKER_COMPOSE_FILE% up -d
                '''
            }
        }

        stage('Wait for health') {
            steps {
                bat '''
                    echo Waiting for services to become healthy...
                    timeout /t 20
                '''
            }
        }

        stage('Integration smoke test') {
            steps {
                bat '''
                    echo Running smoke test...
                    curl -s http://localhost:8081/healthz
                    curl -s http://localhost:8082/healthz
                    curl -s http://localhost:8083/healthz
                '''
            }
        }

        stage('Push images (optional)') {
            when {
                expression { return env.BRANCH_NAME == 'master' }
            }
            steps {
                bat '''
                    echo Pushing Docker images to registry (if configured)...
                    REM docker login -u %DOCKER_USER% -p %DOCKER_PASS%
                    REM docker push <your-repo>/usersvc:latest
                    REM docker push <your-repo>/productsvc:latest
                    REM docker push <your-repo>/ordersvc:latest
                '''
            }
        }
    }

    post {
        always {
            bat 'docker compose -f %DOCKER_COMPOSE_FILE% down || echo Nothing to clean'
            echo "✔ Pipeline finished (cleanup done)"
        }
        failure {
            echo "❌ Pipeline failed"
        }
    }
}