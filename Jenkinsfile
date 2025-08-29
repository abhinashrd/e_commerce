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

        stage('Build Docker Images') {
            parallel {
                stage('Order Service Image') {
                    steps {
                        bat 'docker build -t ordersvc:latest ./services/ordersvc'
                    }
                }
                stage('Product Service Image') {
                    steps {
                        bat 'docker build -t productsvc:latest ./services/productsvc'
                    }
                }
                stage('User Service Image') {
                    steps {
                        bat 'docker build -t usersvc:latest ./services/usersvc'
                    }
                }
            }
        }

        stage('Login to DockerHub') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'dockerhub-cred', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    bat """
                    echo %DOCKER_PASS% | docker login -u %DOCKER_USER% --password-stdin
                    """
                }
            }
        }

        stage('Push Images') {
            steps {
                bat '''
                docker tag ordersvc:latest abhinashd/ordersvc:latest
                docker tag productsvc:latest abhinashd/productsvc:latest
                docker tag usersvc:latest abhinashd/usersvc:latest

                docker push abhinashd/ordersvc:latest
                docker push abhinashd/productsvc:latest
                docker push abhinashd/usersvc:latest
                '''
            }
        }
    }
}
