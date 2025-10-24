# 🛡️ Self-Healing Infrastructure with Chaos Engineering

**Project Type:** DevOps & SRE
**Skill Focus:** Kubernetes | Prometheus | Grafana | LitmusChaos | GoLang | PostgreSQL | Automation
**Automation Level:** Zero-Touch Infrastructure

---

## 📘 Project Overview

This project demonstrates an industry-ready, automated Kubernetes platform ensuring application resiliency and availability through:

* **💠 Self-Healing Infrastructure:** Automatically recovers from pod or database failures.
* **📊 Active Monitoring:** Prometheus + Grafana track and visualize microservice health.
* **⚡ Chaos Engineering:** LitmusChaos simulates real-world failures (pod kill, DB outage, latency).
* **🧠 Automation & Zero Touch:** Resilience achieved through CronJobs + Prometheus alerts for hands-off remediation.

---

## 🧩 Architecture

**Core Components:**

| Component                    | Description                                                    |
| ---------------------------- | -------------------------------------------------------------- |
| Golang Microservice          | Sample “Voting App” emitting custom Prometheus metrics         |
| PostgreSQL                   | Application relational datastore (via Bitnami Helm Chart)      |
| Prometheus & Grafana         | Observability, metrics visualization, and alerting             |
| LitmusChaos Agent            | Simulates pod kill and DB chaos for self-healing validation    |
| Kubernetes (Minikube/Docker) | Container orchestration platform                               |
| CronJob Controller           | Detects and resets failed/unhealthy pods via Prometheus alerts |

---

## 🔁 System Capabilities

| Capability                 | Description                                                      |
| -------------------------- | ---------------------------------------------------------------- |
| 💥 Automatic Recovery      | Restarts failed or stuck pods using alert-based CronJobs         |
| 🧩 Chaos Experimentation   | Periodically injects controlled chaos to test resiliency         |
| ⏲️ Real-Time Dashboards    | Live charts for app availability, endpoint latency & DB health   |
| 🔐 Error Pattern Detection | Prometheus alert rules and application probes identify anomalies |
| 💡 SRE Best Practices      | Combines monitoring, fault tolerance, and automated recovery     |

---

## ⚙️ Minimal Technical Implementation

1. **Setup Environment:**

   * Kubernetes cluster (Minikube/Docker)
   * Helm 3 and kubectl installed

2. **Deploy Application & Database:**

   * PostgreSQL via Bitnami Helm Chart
   * Go App container deployed on Kubernetes
   * Credentials managed with Kubernetes Secrets

3. **Setup Monitoring & Chaos Engineering:**

   * Prometheus scrapes metrics; Grafana visualizes dashboards
   * LitmusChaos deployed for chaos testing

4. **Automated Remediation:**

   * CronJob detects failed pods (`CrashLoopBackOff`)
   * Prometheus alerts trigger auto-recovery

---

## 🧠 Key Learnings & Troubleshooting

| Challenge                          | Resolution                                          |
| ---------------------------------- | --------------------------------------------------- |
| ❌ Database auth failures           | Kubernetes Secret synchronization                   |
| ⚙️ MongoDB user mismatch           | LitmusChaos configuration patch                     |
| ⛓️ Agent stuck in PENDING          | Corrected namespace & context via YAML re-apply     |
| 📈 Prometheus endpoint not scraped | Fixed ServiceMonitor config & labels                |
| 🚀 CronJob not triggering          | Proper RBAC permissions and label selectors applied |

---

## 📈 Observability Insights

* **Metric:** `go_app_http_requests_total`
* **Alert Triggers:** CPU > 85% or service unresponsive > 15s
* **Dashboards:** HTTP requests, DB connections, app uptime
* **Chaos Testing:** Validates fault injection impact and recovery

---

## 🚀 Outcomes

* ✅ Demonstrated DevOps & SRE automation proficiency
* ✅ Achieved resilient microservice architecture
* ✅ Implemented hands-free recovery system
* ✅ Adopted observability-driven management

**Recruiter Key Takeaway:**
This project bridges DevOps theory with production-ready execution, proving real fault tolerance, zero downtime recovery, and continuous monitoring.

---

## 🧩 Tools & Frameworks

| Category          | Tools                           |
| ----------------- | ------------------------------- |
| Containerization  | Docker, Kubernetes, Minikube    |
| Monitoring        | Prometheus, Grafana             |
| Chaos Engineering | LitmusChaos                     |
| Programming       | Golang                          |
| Database          | PostgreSQL                      |
| Automation        | CronJob + Prometheus API        |
| Visualization     | Grafana Dashboards              |
| Version Control   | GitHub (GHCR for image hosting) |

---

## 📚 Quick Start (Minimal Commands)

```bash
# Start Minikube
minikube start --driver=docker

# Deploy PostgreSQL
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-postgres bitnami/postgresql

# Deploy Go App
kubectl apply -f go-app-deployment.yaml
kubectl apply -f go-app-service.yaml

# Setup Monitoring
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install monitoring prometheus-community/kube-prometheus-stack --namespace monitoring --create-namespace

# Setup LitmusChaos
helm repo add litmuschaos https://litmuschaos.github.io/litmus-helm/
helm install litmus litmuschaos/litmus --namespace litmus --create-namespace

# Test Self-Healing
kubectl delete pod -l app=go-app  # Observe auto-restart
```

---

## 🧾 Author

**Created by:** Praveen – Cloud & DevOps Engineer
**Specialization:** Kubernetes | AWS | OCI | CI/CD | SRE | Observability Automation
