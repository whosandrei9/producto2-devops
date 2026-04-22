pipeline {
    agent any

    environment {
        APP_NAME = 'producto2-uoc'
        MINIKUBE = '/usr/local/bin/minikube'
        KUBECTL  = '/usr/local/bin/kubectl'
    }

    stages {
        stage('Preparar variables') {
            steps {
                script {
                    if (env.BRANCH_NAME == 'main') {
                        env.K8S_NAMESPACE = 'produccion'
                        env.REPLICAS = '2'
                        env.APP_ENV = 'produccion'
                    } else if (env.BRANCH_NAME == 'test') {
                        env.K8S_NAMESPACE = 'pruebas'
                        env.REPLICAS = '2'
                        env.APP_ENV = 'pruebas'
                    } else {
                        error("Solo se despliegan automáticamente las ramas main y test.")
                    }

                    env.IMAGE_TAG = "${APP_NAME}:${BRANCH_NAME}-${BUILD_NUMBER}"
                }
            }
        }

        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Asegurar cluster minikube') {
            steps {
                sh """
                    if ! ${MINIKUBE} status >/dev/null 2>&1; then
                        ${MINIKUBE} start --driver=docker --nodes 2
                    fi
                """
            }
        }

        stage('Construir imagen en minikube') {
            steps {
                sh """
                    ${MINIKUBE} image build --all -t ${IMAGE_TAG} .
                """
            }
        }

        stage('Desplegar en Kubernetes') {
            steps {
                sh """
                    ${KUBECTL} create namespace ${K8S_NAMESPACE} --dry-run=client -o yaml | ${KUBECTL} apply -f -

                    sed -e "s|__NAMESPACE__|${K8S_NAMESPACE}|g" \
                        -e "s|__IMAGE__|${IMAGE_TAG}|g" \
                        -e "s|__REPLICAS__|${REPLICAS}|g" \
                        -e "s|__APP_ENV__|${APP_ENV}|g" \
                        k8s/deployment.yaml | ${KUBECTL} apply -f -

                    sed -e "s|__NAMESPACE__|${K8S_NAMESPACE}|g" \
                        k8s/service.yaml | ${KUBECTL} apply -f -

                    ${KUBECTL} rollout status deployment/${APP_NAME} -n ${K8S_NAMESPACE} --timeout=180s
                """
            }
        }

        stage('Verificar despliegue') {
            steps {
                sh """
                    ${KUBECTL} get pods -n ${K8S_NAMESPACE} -o wide
                    ${KUBECTL} get svc -n ${K8S_NAMESPACE}
                """
            }
        }
    }
}
