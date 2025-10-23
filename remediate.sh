#!/bin/bash
# Script: remediate.sh
# Purpose: Checks Alertmanager for critical alerts and restarts the failing application.

# --- Configuration ---
# Internal Kubernetes Service URL for Alertmanager (running in monitoring namespace)
ALERTMANAGER_URL="http://alertmanager-monitoring-kube-prometheus-alertmanager.monitoring:9093/api/v2/alerts"
APP_DEPLOYMENT="go-app"
APP_NAMESPACE="default"
ALERT_SEVERITY="critical"

echo "--- $(date) ---"
echo "Checking Alertmanager for critical alerts..."

# 1. Query Alertmanager for currently FIRING alerts
ALERT_STATUS=$(curl -s -G $ALERTMANAGER_URL --data-urlencode "filter=status=firing,severity=${ALERT_SEVERITY}")

# Check if the response contains an array of firing alerts (meaning alerts are active)
FIRING_COUNT=$(echo $ALERT_STATUS | jq '. | length')

if [[ $FIRING_COUNT -gt 0 ]]; then
    echo "🚨 Critical Alert detected! Firing count: $FIRING_COUNT"
    echo "Attempting remediation action: Restarting $APP_DEPLOYMENT in $APP_NAMESPACE..."

    # 2. Remediation Action: Force a rolling restart of the application deployment
    # The kubectl binary is automatically available in the Minikube environment
    /usr/local/bin/kubectl rollout restart deployment $APP_DEPLOYMENT -n $APP_NAMESPACE

    if [ $? -eq 0 ]; then
        echo "✅ Remediation successful! Deployment restart triggered."
    else
        echo "❌ Remediation FAILED! Could not restart deployment."
    fi
else
    echo "✅ System is healthy. No critical alerts firing."
fi

echo "-------------------"