pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                git branch: 'master',
                    url: 'https://github.com/abhinashrd/e_commerce.git'
            }
        }

        stage('Build Services') {
            parallel {
                stage('Build Order Service') {
                    steps {
                        bat '''
                        cd services\\ordersvc
                        go build -o ordersvc.exe
                        '''
                    }
                }
                stage('Build Product Service') {
                    steps {
                        bat '''
                        cd services\\productsvc
                        go build -o productsvc.exe
                        '''
                    }
                }
                stage('Build User Service') {
                    steps {
                        bat '''
                        cd services\\usersvc
                        go build -o usersvc.exe
                        '''
                    }
                }
            }
        }

        stage('Run Tests') {
            parallel {
                stage('Test Order Service') {
                    steps {
                        bat '''
                        cd services\\ordersvc
                        go test ./...
                        '''
                    }
                }
                stage('Test Product Service') {
                    steps {
                        bat '''
                        cd services\\productsvc
                        go test ./...
                        '''
                    }
                }
                stage('Test User Service') {
                    steps {
                        bat '''
                        cd services\\usersvc
                        go test ./...
                        '''
                    }
                }
            }
        }

        stage('Build & Push Docker Images') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'dockerhub-cred', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    bat """
                    echo %DOCKER_PASS% | docker login -u %DOCKER_USER% --password-stdin
                    
                    docker build -t %DOCKER_USER%/ordersvc:latest ./services/ordersvc
                    docker build -t %DOCKER_USER%/productsvc:latest ./services/productsvc
                    docker build -t %DOCKER_USER%/usersvc:latest ./services/usersvc

                    docker push %DOCKER_USER%/ordersvc:latest
                    docker push %DOCKER_USER%/productsvc:latest
                    docker push %DOCKER_USER%/usersvc:latest
                    """
                }
            }
        }
    }
}
