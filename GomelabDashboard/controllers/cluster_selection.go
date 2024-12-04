package controllers

import (
	"fmt"
	"gomelabdashboard/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetServices(c *gin.Context) {
	clusterIndex := c.Query("cluster") // Cluster index from query params

	// Load Kubernetes clients
	clients, err := utils.LoadKubeClientsFromSingleFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load kubeconfigs"})
		return
	}

	// Validate cluster index
	clusterID, err := strconv.Atoi(clusterIndex)
	if err != nil || clusterID < 0 || clusterID >= len(clients) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cluster ID"})
		return
	}

	// Get Kubernetes client for the selected cluster
	client := clients[clusterID]

	// Fetch services
	services, err := client.CoreV1().Services("").List(c, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list services", "details": err.Error()})
		return
	}

	// Fetch endpoints
	endpoints, err := client.CoreV1().Endpoints("").List(c, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list endpoints", "details": err.Error()})
		return
	}

	// Map service to its details
	var serviceDetails []map[string]interface{}
	for _, service := range services.Items {
		serviceInfo := map[string]interface{}{
			"name":         service.Name,
			"namespace":    service.Namespace,
			"clusterIP":    service.Spec.ClusterIP,
			"type":         service.Spec.Type,
			"ports":        service.Spec.Ports,
			"selector":     service.Spec.Selector,
			"annotations":  service.Annotations,
			"creationTime": service.CreationTimestamp.String(),
		}

		// Add endpoint information if available
		for _, ep := range endpoints.Items {
			if ep.Name == service.Name && ep.Namespace == service.Namespace {
				var endpointAddresses []string
				for _, subset := range ep.Subsets {
					for _, address := range subset.Addresses {
						endpointAddresses = append(endpointAddresses, fmt.Sprintf("%s:%d", address.IP, subset.Ports[0].Port))
					}
				}
				serviceInfo["endpoints"] = endpointAddresses
			}
		}

		serviceDetails = append(serviceDetails, serviceInfo)
	}

	c.JSON(http.StatusOK, gin.H{"services": serviceDetails})
}
