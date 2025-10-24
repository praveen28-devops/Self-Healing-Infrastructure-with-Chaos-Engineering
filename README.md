🛡️ Self-Healing Infrastructure with Chaos Engineering
Project Type: DevOps & SRE
Skill Focus: Kubernetes - Prometheus - Grafana - LitmusChaos - GoLang - Postgres - Automation
Automation Level: 100% (Zero-Touch Infrastructure)
📘 Project Overview

This project demonstrates an industry-ready, automated Kubernetes platform that ensures application resiliency and availability through:

    💠 Self-Healing Infrastructure: Automatically recovers from pod or database failures

    📊 Active Monitoring: Prometheus + Grafana track and visualize microservice health

    ⚡ Chaos Engineering: Uses LitmusChaos to simulate real-world failures (pod kill, DB outage, latency)

    🧠 Automation + Zero Touch: Resilience achieved through CronJobs + Prometheus alerts + hands-off remediation

🧩 Architecture

Core Components:

    Golang Microservice — Sample “Voting App” emitting custom Prometheus metrics

    PostgreSQL — Application relational datastore (deployed via Bitnami Helm Chart)

    Prometheus & Grafana — Observability, metrics visualization, and alerting

    LitmusChaos Agent — Simulates pod kill and DB chaos to test self-healing

    Kubernetes (Minikube/Docker) — Container orchestration platform

    CronJob Controller — Resets failed or unhealthy pods via custom Prometheus triggers

🔁 What the System Can Do
Capability	Description
💥 Automatic Recovery	Restarts failed or stuck pods using alert-based CronJobs
🧩 Chaos Experimentation	Periodically injects controlled chaos to test resiliency
⏲️ Real-Time Health Dashboards	Live charts for app availability, endpoint latency & DB health
🔐 Error Pattern Detection	Through Prometheus alert rules and application probes
💡 SRE Best Practices	Combines monitoring, fault tolerance, and automated recovery
⚙️ Step-by-Step Technical Implementation
1. Environment Setup

    Kubernetes cluster launched via Minikube (Docker driver)

    Installed Helm 3 and kubectl for deployment management

2. Application & Database

    Containerized Go App, image hosted on GitHub Container Registry (GHCR)

    PostgreSQL deployed through the Bitnami Helm Chart

    Database credentials stored and managed via Kubernetes Secrets

3. Monitoring Configuration

    Prometheus installed to scrape custom metrics such as go_app_http_requests_total

    Grafana dashboards visualized live metrics and triggered alerts on anomalies

    Exposed application through Kubernetes Service for browser access

4. Chaos Engineering Integration

    Deployed LitmusChaos portal and registered the self-agent

    Performed chaos tests (Pod Kill, DB Network Delay, CPU spikes)

    Verified automatic stabilization after injected failures

5. Automated Remediation

    Created a CronJob for:

        Detecting “CrashLoopBackOff” pods

        Restarting or redeploying automatically

    Alerts from Prometheus triggered self-healing workflows

    Validated zero manual recovery required

🧠 Key Learnings & Troubleshooting Highlights
Challenge	Resolution
❌ Database auth failures (CrashLoopBackOff)	Fixed via Kubernetes secret synchronization
⚙️ MongoDB user mismatch during agent setup	Resolved using Litmus configuration patch
⛓️ Agent stuck in PENDING	Fixed namespace & context mismatch via YAML re-apply
📈 Prometheus /metrics endpoint not scraped	Corrected ServiceMonitor configuration & selector labels
🚀 CronJob remediation not triggering	Attached RBAC permissions and label selectors properly
🧮 Observability Insights

    Custom metric: go_app_http_requests_total

    Alert Trigger: CPU > 85% or service unresponsive for >15s

    Grafana Dashboards: Visualized HTTP requests, DB connections, and app uptime

    LitmusChaos Integration: Validates fault injection impact and recovery time

📸 Screenshots

    App Service (VS Code + Minikube) – service exposed for testing

    Prometheus Dashboard – live HTTP request metrics

    Pod Failure Logs (Debug Example) – sample recovery scenario

    LitmusChaos Portal – registered self-healing infrastructure agent

🚀 Outcomes

    ✅ Demonstrated DevOps & SRE automation proficiency

    ✅ Achieved resilient microservice architecture

    ✅ Implemented hands-free recovery system

    ✅ Adopted observability-driven management (Prometheus + Grafana)

💼 Recruiter / HR Key Takeaway

    This project bridges DevOps theory and production-ready execution.
    It proves real fault tolerance, zero downtime recovery, and continuous monitoring — exactly what modern cloud-native engineering teams demand.

🧩 Tools & Frameworks
Category	Tools Used
Containerization	Docker, Kubernetes, Minikube
Monitoring	Prometheus, Grafana
Chaos Engineering	LitmusChaos
Programming	Golang
Database	PostgreSQL
Automation	CronJob + Prometheus API
Visualization	Grafana Dashboards
Version Control	GitHub (GHCR for image hosting)
📚 How to Recreate This Project

bash
# 1. Start Minikube
minikube start --driver=docker

# 2. Deploy PostgreSQL via Helm
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-postgres bitnami/postgresql

# 3. Deploy Go App
kubectl apply -f go-app-deployment.yaml
kubectl apply -f go-app-service.yaml

# 4. Setup Monitoring Tools
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install monitoring prometheus-community/kube-prometheus-stack --namespace monitoring --create-namespace

# 5. Setup LitmusChaos
helm repo add litmuschaos https://litmuschaos.github.io/litmus-helm/
helm install litmus litmuschaos/litmus --namespace litmus --create-namespace

# 6. Test Self-Healing
kubectl delete pod -l app=go-app  # Observe auto-restart

🧾 Author

Created by: Praveen – Cloud & DevOps Engineer
Specialization: Kubernetes | AWS | OCI | CI/CD | SRE | Observability Automation
⭐ If you find this useful, give the repository a star on GitHub and connect on LinkedIn for collaboration discussions!